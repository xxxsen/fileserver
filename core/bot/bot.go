package bot

import (
	"context"
	"fileserver/core"
	"fmt"
	"net"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	lru "github.com/hnlq715/golang-lru"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

type TGBot struct {
	c        *config
	bot      *tgbotapi.BotAPI
	client   *http.Client
	cacheLnk *lru.Cache
}

const (
	defaultMaxTGBotFileSize   = 4 * 1024 * 1024 * 1024
	defaultTGBotFileBlockSize = 20 * 1024 * 1024
)

func New(opts ...Option) (*TGBot, error) {
	c := &config{
		fsize:   defaultMaxTGBotFileSize,
		blksize: defaultTGBotFileBlockSize,
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
	cacheLnk, _ := lru.New(20000)

	//init http client
	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).Dial,
			IdleConnTimeout: 20 * time.Second,
			MaxIdleConns:    20,
		},
	}

	return &TGBot{c: c, cacheLnk: cacheLnk, client: httpClient, bot: botClient}, nil
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

func (c *TGBot) FileUpload(ctx context.Context, uctx *core.FileUploadRequest) (*core.FileUploadResponse, error) {
	sname := uuid.NewString()
	freader := tgbotapi.FileReader{
		Name:   sname,
		Reader: uctx.ReadSeeker,
	}
	doc := tgbotapi.NewDocument(c.c.chatid, freader)
	doc.DisableNotification = true
	msg, err := c.bot.Send(doc)
	if err != nil {
		return nil, err
	}
	extra, err := encodeFileExtra(&botFileCtx{
		ChatId:   c.c.chatid,
		FileType: fileTypeOneFile,
	})
	if err != nil {
		return nil, errs.Wrap(errs.ErrMarshal, "encode bot extra fail", err)
	}
	if err != nil {
		return nil, err
	}
	return &core.FileUploadResponse{
		Key:   msg.Document.FileID,
		Extra: extra,
	}, nil

}

func (c *TGBot) cacheGetURL(ctx context.Context, hash string) (string, error) {
	if lnk, ok := c.cacheLnk.Get(hash); ok {
		return lnk.(string), nil
	}

	cf := tgbotapi.FileConfig{FileID: hash}
	f, err := c.bot.GetFile(cf)
	if err != nil {
		return "", err
	}
	lnk := f.Link(c.bot.Token)
	//这里应该能1小时有效的...
	c.cacheLnk.AddEx(hash, lnk, 30*time.Minute)
	return lnk, nil
}

func (c *TGBot) FileDownload(ctx context.Context, fctx *core.FileDownloadRequest) (*core.FileDownloadResponse, error) {
	//TODO: we need to check whether it is a part file
	lnk, err := c.cacheGetURL(ctx, fctx.Key)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, lnk, nil)
	if err != nil {
		return nil, err
	}
	if fctx.StartAt != 0 {
		rangeHeader := fmt.Sprintf("bytes=%d-", fctx.StartAt)
		req.Header.Set("Range", rangeHeader)
	}
	rsp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	//caller should close rsp.Body
	if rsp.StatusCode/100 != 2 {
		rsp.Body.Close()
		return nil, errs.New(errs.ErrServiceInternal, "status code not ok, code:%d", rsp.StatusCode)
	}
	if fctx.StartAt != 0 && len(rsp.Header.Get("Content-Range")) == 0 {
		rsp.Body.Close()
		return nil, errs.New(errs.ErrParam, "not support range")
	}

	return &core.FileDownloadResponse{Reader: rsp.Body}, nil
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
