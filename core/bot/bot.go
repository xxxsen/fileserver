package bot

import (
	"bytes"
	"context"
	"encoding/hex"
	"fileserver/core"
	"fileserver/proto/fileserver/fileinfo"
	"fileserver/utils"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	lru "github.com/hnlq715/golang-lru"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type TGBot struct {
	c         *config
	bot       *tgbotapi.BotAPI
	metaCache *lru.Cache
}

const (
	defaultMaxTGBotFileSize   = 4 * 1024 * 1024 * 1024
	defaultTGBotFileBlockSize = 20 * 1024 * 1024
)

func New(opts ...Option) (*TGBot, error) {
	c := &config{
		fsize:   defaultMaxTGBotFileSize,
		blksize: defaultTGBotFileBlockSize,
		tmpdir:  os.TempDir(),
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.chatid == 0 || len(c.token) == 0 {
		return nil, errs.New(errs.ErrParam, "invalid chatid/token")
	}
	botClient, err := newBotClient(c.chatid, c.token)
	if err != nil {
		return nil, errs.Wrap(errs.ErrServiceInternal, "new bot client fail", err)
	}
	metaCache, _ := lru.New(10000)

	return &TGBot{c: c, metaCache: metaCache, bot: botClient}, nil
}

func asyncUpdate(chatid int64, bot *tgbotapi.BotAPI) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		logutil.GetLogger(context.Background()).With(zap.Int64("userid", chatid), zap.Int64("senderid", update.Message.Chat.ID),
			zap.String("message", update.Message.Text)).
			Info("recv message from remote")
	}
	return nil
}

func newBotClient(chatid int64, token string) (*tgbotapi.BotAPI, error) {
	//parse config
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errs.Wrap(errs.ErrServiceInternal, "init bot fail", err)
	}
	go func() {
		err := asyncUpdate(chatid, bot)
		if err != nil {
			logutil.GetLogger(context.Background()).With(zap.Error(err)).Error("async update bot fail")
		}
	}()
	return bot, nil
}

func (c *TGBot) BlockSize() int64 {
	return c.c.blksize
}

func (c *TGBot) MaxFileSize() int64 {
	return c.c.fsize
}

func (c *TGBot) uploadOne(ctx context.Context, r io.Reader, sz int64) (string, string, error) {
	sname := uuid.NewString()
	mReader := MD5Reader(r)
	cReader := CountReader(mReader)
	freader := tgbotapi.FileReader{
		Name:   sname,
		Reader: cReader,
	}
	doc := tgbotapi.NewDocument(c.c.chatid, freader)
	doc.DisableNotification = true
	msg, err := c.bot.Send(doc)
	if err != nil {
		return "", "", errs.Wrap(errs.ErrIO, "send document fail", err)
	}
	if int64(cReader.GetCount()) != sz {
		return "", "", errs.New(errs.ErrIO, "send document size not match, write:%d, need:%d", cReader.GetCount(), sz)
	}
	return msg.Document.FileID, hex.EncodeToString(mReader.GetSum()), nil
}

func (c *TGBot) singleFileUpload(ctx context.Context, uctx *core.FileUploadRequest) (string, error) {
	fileid, ck, err := c.uploadOne(ctx, uctx.ReadSeeker, uctx.Size)
	if err != nil {
		return "", errs.Wrap(errs.ErrIO, "upload one part fail", err)
	}
	if len(uctx.MD5) != 0 && ck != uctx.MD5 {
		return "", errs.New(errs.ErrParam, "checksum not match, calc:%s, get:%s", ck, uctx.MD5)
	}
	return fileid, nil
}

func (c *TGBot) multipartFileUpload(ctx context.Context, uctx *core.FileUploadRequest) (string, error) {
	blkcount := utils.CalcFileBlockCount(uint64(uctx.Size), uint64(c.BlockSize()))
	blklist := make([]string, 0, blkcount)
	for i := 0; i < blkcount; i++ {
		partreader := io.LimitReader(uctx.ReadSeeker, c.BlockSize())
		blkidsz := utils.CalcBlockSize(uint64(uctx.Size), uint64(c.BlockSize()), i)
		if blkidsz == 0 {
			return "", errs.New(errs.ErrParam, "invalid blkidsize, id:%d, get:%d", i, blkidsz)
		}
		fid, _, err := c.uploadOne(ctx, partreader, int64(blkidsz))
		if err != nil {
			return "", errs.Wrap(errs.ErrIO, fmt.Sprintf("upload block fail, id:%d", i), err)
		}
		blklist = append(blklist, fid)
	}
	filectx := &fileinfo.BotUploadContext{
		FileSize:  proto.Int64(uctx.Size),
		BlockSize: proto.Int64(c.BlockSize()),
		Blocks:    blklist,
	}
	raw, err := utils.EncodeBotUploadContext(filectx)
	if err != nil {
		return "", errs.Wrap(errs.ErrMarshal, "encode upload ctx fail", err)
	}
	fid, _, err := c.uploadOne(ctx, bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		fid, _, err = c.uploadOne(ctx, bytes.NewReader(raw), int64(len(raw)))
	}
	if err != nil {
		return "", errs.Wrap(errs.ErrIO, "save multipart context fail", err)
	}
	return fid, nil
}

func (c *TGBot) FileUpload(ctx context.Context, uctx *core.FileUploadRequest) (*core.FileUploadResponse, error) {
	var (
		uploader       = c.singleFileUpload
		filetype int32 = int32(fileinfo.BotConstants_BOT_FILE_TYPE_SINGLE)
	)
	if uctx.Size > c.BlockSize() {
		uploader = c.multipartFileUpload
		filetype = int32(fileinfo.BotConstants_BOT_FILE_TYPE_MULTIPART)
	}
	fileid, err := uploader(ctx, uctx)
	if err != nil {
		return nil, errs.Wrap(errs.ErrIO, "upload file fail", err)
	}
	extra, err := utils.EncodeBotFileExtra(&fileinfo.BotFileExtra{
		ChatId:    proto.Int64(c.c.chatid),
		FileType:  proto.Int32(filetype),
		BlockSize: proto.Int64(c.BlockSize()),
		FileSize:  proto.Int64(uctx.Size),
	})
	if err != nil {
		return nil, errs.Wrap(errs.ErrMarshal, "encode bot extra fail", err)
	}
	if err != nil {
		return nil, err
	}
	return &core.FileUploadResponse{
		Key:   fileid,
		Extra: extra,
	}, nil
}

func (c *TGBot) singleFileDownload(ctx context.Context, key string, downat int64) (io.ReadCloser, error) {
	return NewPartReader(ctx, c.bot, key, downat), nil
}

func (c *TGBot) getMultiblockMeta(ctx context.Context, fid string) (*fileinfo.BotUploadContext, error) {
	if v, ok := c.metaCache.Get(fid); ok {
		return v.(*fileinfo.BotUploadContext), nil
	}
	r, err := c.singleFileDownload(ctx, fid, 0)
	if err != nil {
		return nil, errs.Wrap(errs.ErrIO, "get download meta fail", err)
	}
	defer r.Close()
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errs.Wrap(errs.ErrIO, "read meta data fail", err)
	}
	finfo, err := utils.DecodeBotUploadContext(raw)
	if err != nil {
		return nil, errs.Wrap(errs.ErrUnknown,
			fmt.Sprintf("decode upload context fail, ctxdata:%s, ctxdata len:%d", hex.EncodeToString(raw), len(raw)),
			err,
		)
	}
	c.metaCache.Add(fid, finfo)
	return finfo, nil
}

func (c *TGBot) multipartFileDownload(ctx context.Context, fctx *core.FileDownloadRequest) (io.ReadCloser, error) {
	finfo, err := c.getMultiblockMeta(ctx, fctx.Key)
	if err != nil {
		return nil, errs.Wrap(errs.ErrIO, "get meta fail", err)
	}
	return NewMultipartReader(ctx, c.bot, &MultipartMeta{
		StartAt:  fctx.StartAt,
		FileSize: finfo.GetFileSize(),
		BlkSize:  finfo.GetBlockSize(),
		BlkList:  finfo.GetBlocks(),
	}), nil
}

func (c *TGBot) FileDownload(ctx context.Context, fctx *core.FileDownloadRequest) (*core.FileDownloadResponse, error) {
	bctx, err := utils.DecodeBotFileExtra(fctx.Extra)
	if err != nil {
		return nil, errs.Wrap(errs.ErrUnmarshal, "decode bot file context fail", err)
	}
	if bctx.GetChatId() != c.c.chatid {
		return nil, errs.New(errs.ErrParam, "chatid not match, file chatid:%d, current chatid:%d", bctx.GetChatId(), c.c.chatid)
	}
	var r io.ReadCloser
	switch bctx.GetFileType() {
	case int32(fileinfo.BotConstants_BOT_FILE_TYPE_SINGLE):
		r, err = c.singleFileDownload(ctx, fctx.Key, fctx.StartAt)
	case int32(fileinfo.BotConstants_BOT_FILE_TYPE_MULTIPART):
		r, err = c.multipartFileDownload(ctx, fctx)
	default:
		return nil, errs.New(errs.ErrParam, "not support file type:%d", bctx.GetFileType())
	}
	if err != nil {
		return nil, errs.Wrap(errs.ErrIO, "get file reader fail", err)
	}
	return &core.FileDownloadResponse{Reader: r}, nil
}

func (c *TGBot) BeginFileUpload(ctx context.Context, fctx *core.BeginFileUploadRequest) (*core.BeginFileUploadResponse, error) {
	//TODO:
	panic(1)
}

func (c *TGBot) PartFileUpload(ctx context.Context, pctx *core.PartFileUploadRequest) (*core.PartFileUploadResponse, error) {
	//TODO:
	panic(1)
}

func (c *TGBot) FinishFileUpload(ctx context.Context, fctx *core.FinishFileUploadRequest) (*core.FinishFileUploadResponse, error) {
	//TODO:
	panic(1)
}
