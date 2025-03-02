package bot

import (
	"context"
	"fileserver/utils"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	lru "github.com/hnlq715/golang-lru"
)

var linkCache, _ = lru.New(20000)
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

type PartReader struct {
	ctx     context.Context
	bot     *tgbotapi.BotAPI
	blk     string
	startat int64
	isOpen  bool
	r       io.ReadCloser
	initer  sync.Once
}

func NewPartReader(ctx context.Context, bot *tgbotapi.BotAPI, blk string, at int64) *PartReader {
	return &PartReader{
		ctx:     ctx,
		blk:     blk,
		startat: at,
		bot:     bot,
		isOpen:  false,
	}
}

func (r *PartReader) cacheGetURL(hash string) (string, error) {
	if lnk, ok := linkCache.Get(hash); ok {
		return lnk.(string), nil
	}

	cf := tgbotapi.FileConfig{FileID: hash}
	f, err := r.bot.GetFile(cf)
	if err != nil {
		return "", err
	}
	lnk := f.Link(r.bot.Token)
	//这里应该能1小时有效的...
	linkCache.AddEx(hash, lnk, 30*time.Minute)
	return lnk, nil
}

func (r *PartReader) initReader() error {
	lnk, err := r.cacheGetURL(r.blk)
	if err != nil {
		return fmt.Errorf("get link fail, err:%w", err)
	}
	req, err := http.NewRequestWithContext(r.ctx, http.MethodGet, lnk, nil)
	if err != nil {
		return fmt.Errorf("create http request fail, err:%w", err)
	}
	if r.startat != 0 {
		rangeHeader := fmt.Sprintf("bytes=%d-", r.startat)
		req.Header.Set("Range", rangeHeader)
	}
	rsp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("do http request fail, err:%w", err)
	}
	//caller should close rsp.Body
	if rsp.StatusCode/100 != 2 {
		rsp.Body.Close()
		return fmt.Errorf("status code not ok, code:%d", rsp.StatusCode)
	}
	if r.startat != 0 && len(rsp.Header.Get("Content-Range")) == 0 {
		rsp.Body.Close()
		return fmt.Errorf("not support range")
	}
	r.r = rsp.Body
	r.isOpen = true
	return nil
}

func (r *PartReader) Read(buf []byte) (int, error) {
	var err error
	r.initer.Do(func() {
		err = r.initReader()
	})
	if err != nil {
		return 0, fmt.Errorf("init reader fail, err:%w", err)
	}
	if !r.isOpen {
		return 0, fmt.Errorf("reader:%s closed", r.blk)
	}
	return r.r.Read(buf)
}

func (r *PartReader) Close() error {
	var err error
	if r.isOpen {
		err = r.r.Close()
		r.isOpen = false
	}
	if err != nil {
		return err
	}
	return nil
}

type MultipartMeta struct {
	StartAt  int64
	FileSize int64
	BlkSize  int64
	BlkList  []string
}

type MultipartReader struct {
	ctx       context.Context
	bot       *tgbotapi.BotAPI
	meta      *MultipartMeta
	initer    sync.Once
	partindex int
	reader    *PartReader
	isOpen    bool
}

func NewMultipartReader(ctx context.Context, bot *tgbotapi.BotAPI, meta *MultipartMeta) *MultipartReader {
	return &MultipartReader{
		ctx:       ctx,
		bot:       bot,
		meta:      meta,
		isOpen:    false,
		partindex: 0,
	}
}

func (r *MultipartReader) initReader() error {
	if r.meta.StartAt > r.meta.FileSize {
		return fmt.Errorf("start at:%d > filesize:%d", r.meta.StartAt, r.meta.FileSize)
	}
	blktotal := utils.CalcFileBlockCount(uint64(r.meta.FileSize), uint64(r.meta.BlkSize))
	if blktotal != len(r.meta.BlkList) {
		return fmt.Errorf("file blk size not match, total:%d, calc:%d", len(r.meta.BlkList), blktotal)
	}

	r.partindex = int(r.meta.StartAt) / int(r.meta.BlkSize)
	at := r.meta.StartAt % r.meta.BlkSize
	r.reader = NewPartReader(r.ctx, r.bot, r.meta.BlkList[r.partindex], at)
	r.isOpen = true
	return nil
}

func (r *MultipartReader) Read(buf []byte) (int, error) {
	var err error
	r.initer.Do(func() {
		err = r.initReader()
	})
	if err != nil {
		return 0, fmt.Errorf("init reader fail, err:%w", err)
	}
	if !r.isOpen {
		return 0, fmt.Errorf("reader already closed")
	}
	cnt, err := r.reader.Read(buf)
	if err == io.EOF {
		if cerr := r.reader.Close(); cerr != nil {
			return cnt, fmt.Errorf("close internal part fail, index:%d, err:%w", r.partindex, err)
		}
		if r.partindex+1 < len(r.meta.BlkList) {
			r.partindex++
			r.reader = NewPartReader(r.ctx, r.bot, r.meta.BlkList[r.partindex], 0)
			//multi part, so it should reset io.EOF to nil
			err = nil
		}
	}
	return cnt, err
}

func (r *MultipartReader) Close() error {
	var err error
	if r.isOpen {
		err = r.reader.Close()
		r.isOpen = false
	}
	if err != nil {
		return err
	}
	return nil
}
