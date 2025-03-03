package telegram

import (
	"context"
	"fileserver/filesystem"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	lru "github.com/hnlq715/golang-lru"
)

const (
	defaultMaxFileSize = 20 * 1024 * 1024
)

var defaultHTTPClient = &http.Client{
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).Dial,
		IdleConnTimeout: 20 * time.Second,
		MaxIdleConns:    20,
	},
}

type telegramFileSystem struct {
	chatid    int64
	token     string
	bot       *tgbotapi.BotAPI
	linkCache *lru.Cache
}

func New(chatid int64, token string) (filesystem.IFileSystem, error) {
	cache, err := lru.New(1000)
	if err != nil {
		return nil, err
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("init bot fail, err:%w", err)
	}
	return &telegramFileSystem{
		chatid:    chatid,
		token:     token,
		bot:       bot,
		linkCache: cache,
	}, nil
}

func (t *telegramFileSystem) MaxFileSize() int64 {
	return defaultMaxFileSize
}

func (t *telegramFileSystem) Upload(ctx context.Context, r io.Reader) (string, error) {
	sname := uuid.NewString()
	freader := tgbotapi.FileReader{
		Name:   sname,
		Reader: r,
	}
	doc := tgbotapi.NewDocument(t.chatid, freader)
	doc.DisableNotification = true
	msg, err := t.bot.Send(doc)
	if err != nil {
		return "", fmt.Errorf("send document fail, err:%w", err)
	}

	return msg.Document.FileID, nil
}

func (t *telegramFileSystem) cacheGetDownloadLink(filekey string) (string, error) {
	if lnk, ok := t.linkCache.Get(filekey); ok {
		return lnk.(string), nil
	}
	cf := tgbotapi.FileConfig{FileID: filekey}
	f, err := t.bot.GetFile(cf)
	if err != nil {
		return "", err
	}
	lnk := f.Link(t.bot.Token)
	//这里应该能1小时有效的...
	t.linkCache.AddEx(filekey, lnk, 30*time.Minute)
	return lnk, nil
}

func (t *telegramFileSystem) Download(ctx context.Context, filekey string, pos int64) (io.ReadCloser, error) {
	link, err := t.cacheGetDownloadLink(filekey)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return nil, fmt.Errorf("create http request fail, err:%w", err)
	}
	if pos != 0 {
		rangeHeader := fmt.Sprintf("bytes=%d-", pos)
		req.Header.Set("Range", rangeHeader)
	}
	rsp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do http request fail, err:%w", err)
	}
	//caller should close rsp.Body
	if rsp.StatusCode/100 != 2 {
		rsp.Body.Close()
		return nil, fmt.Errorf("status code not ok, code:%d", rsp.StatusCode)
	}
	if pos != 0 && len(rsp.Header.Get("Content-Range")) == 0 {
		rsp.Body.Close()
		return nil, fmt.Errorf("not support range")
	}
	return rsp.Body, nil
}
