package bot

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fileserver/core"
	"fileserver/proto/fileserver/fileinfo"
	"fileserver/utils"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	lru "github.com/hnlq715/golang-lru"
	"github.com/xxxsen/common/errs"
	"google.golang.org/protobuf/proto"
)

type TGBot struct {
	c         *config
	bot       *tgbotapi.BotAPI
	metaCache *lru.Cache
}

const (
	defaultMaxTGBotFileSize    = 4 * 1024 * 1024 * 1024
	defaultTGBotFileBlockSize  = 20 * 1024 * 1024
	defaultPartDataStoreLength = 160
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

func newBotClient(chatid int64, token string) (*tgbotapi.BotAPI, error) {
	//parse config
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errs.Wrap(errs.ErrServiceInternal, "init bot fail", err)
	}
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

func (c *TGBot) singleFileUpload(ctx context.Context, uctx *core.FileUploadRequest) (string, string, error) {
	fileid, ck, err := c.uploadOne(ctx, uctx.ReadSeeker, uctx.Size)
	if err != nil {
		return "", "", errs.Wrap(errs.ErrIO, "upload one part fail", err)
	}
	if len(uctx.MD5) != 0 && ck != uctx.MD5 {
		return "", "", errs.New(errs.ErrParam, "checksum not match, calc:%s, get:%s", ck, uctx.MD5)
	}
	return fileid, ck, nil
}

func (c *TGBot) multipartFileUpload(ctx context.Context, uctx *core.FileUploadRequest) (string, string, error) {
	blkcount := utils.CalcFileBlockCount(uint64(uctx.Size), uint64(c.BlockSize()))
	blklist := make([]string, 0, blkcount)
	md5reader := MD5Reader(uctx.ReadSeeker)
	for i := 0; i < blkcount; i++ {
		partreader := io.LimitReader(md5reader, c.BlockSize())
		blkidsz := utils.CalcBlockSize(uint64(uctx.Size), uint64(c.BlockSize()), i)
		if blkidsz == 0 {
			return "", "", errs.New(errs.ErrParam, "invalid blkidsize, id:%d, get:%d", i, blkidsz)
		}
		fid, _, err := c.uploadOne(ctx, partreader, int64(blkidsz))
		if err != nil {
			return "", "", errs.Wrap(errs.ErrIO, fmt.Sprintf("upload block fail, id:%d", i), err)
		}
		blklist = append(blklist, fid)
	}
	ck := hex.EncodeToString(md5reader.GetSum())
	if len(uctx.MD5) > 0 && ck != uctx.MD5 {
		return "", "", errs.New(errs.ErrParam, "checksum not match, calc:%s, carry:%s", ck, uctx.MD5)
	}
	fid, err := c.writeMultiPartToBot(ctx, uint64(uctx.Size), uint32(c.BlockSize()), blklist)
	if err != nil {
		return "", "", errs.Wrap(errs.ErrIO, "save multipart context fail", err)
	}
	return fid, ck, nil
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
	fileid, cksum, err := uploader(ctx, uctx)
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
		Key:      fileid,
		Extra:    extra,
		CheckSum: cksum,
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

func (c *TGBot) buildTmpDir(xfid string) string {
	return c.c.tmpdir + string(filepath.Separator) + "tgbot_upload" + string(filepath.Separator) + xfid
}

func (c *TGBot) BeginFileUpload(ctx context.Context, fctx *core.BeginFileUploadRequest) (*core.BeginFileUploadResponse, error) {
	xfid := uuid.NewString()
	if err := os.MkdirAll(c.buildTmpDir(xfid), os.ModeDir|os.ModePerm); err != nil {
		return nil, errs.Wrap(errs.ErrServiceInternal, "make dir fail", err)
	}
	upid, err := utils.EncodeUploadID(&fileinfo.UploadIdCtx{
		FileSize:  proto.Uint64(uint64(fctx.FileSize)),
		FileKey:   proto.String(xfid),
		BlockSize: proto.Uint32(uint32(c.BlockSize())),
	})
	if err != nil {
		return nil, err
	}
	return &core.BeginFileUploadResponse{UploadID: upid}, nil
}

func (c *TGBot) isExistFile(file string) (bool, error) {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *TGBot) storePartInfo(dir string, partid int, partkey string, partsize int64, ck string) error {
	raw, err := utils.EncodePartPair(&fileinfo.PartPair{
		PartId:   proto.Int32(int32(partid)),
		PartKey:  proto.String(partkey),
		PartSize: proto.Int64(partsize),
		Md5Value: proto.String(ck),
	})
	if err != nil {
		return errs.Wrap(errs.ErrMarshal, "encode pb fail", err)
	}
	if len(raw) > defaultPartDataStoreLength {
		return fmt.Errorf("part data too long, size:%d", len(raw))
	}
	blockdata := make([]byte, defaultPartDataStoreLength)
	binary.BigEndian.PutUint16(blockdata, uint16(len(raw)))
	copy(blockdata[2:], raw)
	filename := fmt.Sprintf("%s%sdata", dir, string(filepath.Separator))
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	pos := int64(partid * defaultPartDataStoreLength)
	if _, err := f.WriteAt(blockdata, pos); err != nil {
		return errs.Wrap(errs.ErrIO, "write block data fail", err)
	}
	return nil
}

func (c *TGBot) PartFileUpload(ctx context.Context, pctx *core.PartFileUploadRequest) (*core.PartFileUploadResponse, error) {
	uctx, err := utils.DecodeUploadID(pctx.UploadId)
	if err != nil {
		return nil, err
	}
	bkcnt := utils.CalcFileBlockCount(uctx.GetFileSize(), uint64(uctx.GetBlockSize()))
	if pctx.PartId == 0 || pctx.PartId > uint64(bkcnt) {
		return nil, errs.New(errs.ErrParam, "invalid partid:%d", pctx.PartId)
	}
	if pctx.PartId != uint64(bkcnt) && pctx.Size != int64(uctx.GetBlockSize()) {
		return nil, errs.New(errs.ErrParam, "invalid part size, partid:%d, blksize:%d", pctx.PartId, uctx.GetBlockSize())
	}
	if pctx.Size == 0 {
		return nil, errs.New(errs.ErrParam, "empty size")
	}
	dir := c.buildTmpDir(uctx.GetFileKey())
	exist, err := c.isExistFile(dir)
	if err != nil {
		return nil, errs.Wrap(errs.ErrIO, "check upload folder fail", err)
	}
	if !exist {
		return nil, errs.New(errs.ErrNotFound, "upload id not found")
	}
	fileid, ck, err := c.uploadOne(ctx, pctx.ReadSeeker, pctx.Size)
	if err != nil {
		return nil, errs.Wrap(errs.ErrIO, "upload part file fail", err)
	}
	if len(pctx.MD5) > 0 && pctx.MD5 != ck {
		return nil, errs.New(errs.ErrParam, "checksum not match, get:%d, real:%d", pctx.MD5, ck)
	}
	if err := c.storePartInfo(dir, int(pctx.PartId), fileid, pctx.Size, ck); err != nil {
		return nil, errs.Wrap(errs.ErrIO, "store partinfo to disk fail", err)
	}
	return &core.PartFileUploadResponse{}, nil
}

func (c *TGBot) readStorePartInfo(dir string) ([]*fileinfo.PartPair, error) {
	filename := fmt.Sprintf("%s%sdata", dir, string(filepath.Separator))
	exist, err := c.isExistFile(filename)
	if err != nil {
		return nil, errs.Wrap(errs.ErrIO, "check part info fail", err)
	}
	if !exist {
		return nil, errs.New(errs.ErrParam, "part info not found")
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errs.Wrap(errs.ErrIO, "read part info fail", err)
	}
	if len(data)%defaultPartDataStoreLength != 0 {
		return nil, errs.New(errs.ErrParam, "invalid part info, size:%d", len(data))
	}
	blkcount := utils.CalcFileBlockCount(uint64(len(data)), defaultPartDataStoreLength) - 1 //first block is an empty block
	rs := make([]*fileinfo.PartPair, 0, blkcount)
	for i := 0; i < blkcount; i++ {
		partid := i + 1
		start := partid * defaultPartDataStoreLength
		part := data[start : start+defaultPartDataStoreLength]
		length := binary.BigEndian.Uint16(part)
		if length == 0 {
			return nil, errs.New(errs.ErrParam, "invalid part, partid:%d", partid)
		}

		realdata := part[2 : 2+length]
		pair, err := utils.DecodePartPair(realdata)
		if err != nil {
			return nil, errs.Wrap(errs.ErrUnknown, "decode part info fail", err)
		}
		if int(pair.GetPartId()) != i+1 {
			return nil, errs.New(errs.ErrParam, "partid not at its real loc, data partid:%d, partid:%d", partid, pair.GetPartId())
		}
		rs = append(rs, pair)
	}
	return rs, nil
}

func (c *TGBot) FinishFileUpload(ctx context.Context, fctx *core.FinishFileUploadRequest) (*core.FinishFileUploadResponse, error) {
	uctx, err := utils.DecodeUploadID(fctx.UploadId)
	if err != nil {
		return nil, err
	}
	dir := c.buildTmpDir(uctx.GetFileKey())
	exist, err := c.isExistFile(dir)
	if err != nil {
		return nil, errs.Wrap(errs.ErrIO, "check upload folder fail", err)
	}
	if !exist {
		return nil, errs.New(errs.ErrNotFound, "upload id not found")
	}
	parts, err := c.readStorePartInfo(dir)
	if err != nil {
		return nil, errs.Wrap(errs.ErrIO, "read store part info fail", err)
	}
	if len(parts) == 0 {
		return nil, errs.New(errs.ErrParam, "no file part found")
	}
	var calcSize int64
	blks := make([]string, 0, len(parts))
	md5s := make([]string, 0, len(parts))
	for _, item := range parts {
		calcSize += item.GetPartSize()
		blks = append(blks, item.GetPartKey())
		md5s = append(md5s, item.GetMd5Value())
	}
	if calcSize != int64(uctx.GetFileSize()) {
		return nil, errs.New(errs.ErrParam, "file size not match, calc:%d, uctx:%d", calcSize, uctx.GetFileSize())
	}
	filekey := parts[0].GetPartKey()
	filetype := fileinfo.BotConstants_BOT_FILE_TYPE_SINGLE
	if len(parts) > 0 {
		filetype = fileinfo.BotConstants_BOT_FILE_TYPE_MULTIPART
		filekey, err = c.writeMultiPartToBot(ctx, uctx.GetFileSize(), uctx.GetBlockSize(), blks)
		if err != nil {
			return nil, errs.Wrap(errs.ErrIO, "save parts to bot fail", err)
		}
	}
	cks := c.buildETag(md5s)
	extra, err := utils.EncodeBotFileExtra(&fileinfo.BotFileExtra{
		ChatId:    proto.Int64(c.c.chatid),
		FileType:  proto.Int32(int32(filetype)),
		BlockSize: proto.Int64(int64(uctx.GetBlockSize())),
		FileSize:  proto.Int64(int64(uctx.GetFileSize())),
	})
	if err != nil {
		return nil, errs.Wrap(errs.ErrMarshal, "encode file extra fail", err)
	}
	_ = os.RemoveAll(dir)

	return &core.FinishFileUploadResponse{
		Key:      filekey,
		Extra:    extra,
		FileSize: int64(uctx.GetFileSize()),
		CheckSum: cks,
	}, nil
}

func (c *TGBot) buildETag(md5s []string) string {
	if len(md5s) == 0 {
		return ""
	}
	if len(md5s) == 1 {
		return md5s[0]
	}
	m := md5.New()
	for _, item := range md5s {
		m.Write([]byte(item))
	}
	return hex.EncodeToString(m.Sum(nil)) + "-" + strconv.FormatInt(int64(len(md5s)), 10)
}

func (c *TGBot) writeMultiPartToBot(ctx context.Context, filesize uint64, blksize uint32, blks []string) (string, error) {
	filemeta, err := utils.EncodeBotUploadContext(&fileinfo.BotUploadContext{
		FileSize:  proto.Int64(int64(filesize)),
		BlockSize: proto.Int64(int64(blksize)),
		Blocks:    blks,
	})
	if err != nil {
		return "", errs.Wrap(errs.ErrMarshal, "encode bot upload ctx fail", err)
	}
	fileid, _, err := c.uploadOne(ctx, bytes.NewReader(filemeta), int64(len(filemeta)))
	if err != nil {
		return "", errs.Wrap(errs.ErrIO, "save multi part meta fail", err)
	}
	return fileid, nil
}
