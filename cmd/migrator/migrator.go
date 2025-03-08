package main

import (
	"bytes"
	"context"
	"fileserver/cmd/migrator/config"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/xxxsen/common/logger"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

var conf = flag.String("config", "./config.json", "config")

func main() {
	flag.Parse()
	logkit := logger.Init("", "debug", 0, 0, 0, true)
	c, err := config.Parse(*conf)
	if err != nil {
		logkit.Fatal("parse config failed", zap.Error(err))
	}
	ctx := context.Background()
	client := &http.Client{}
	for _, lnk := range c.Src.LinkList {
		logkit.Info("begin migrate link", zap.String("link", lnk))
		if err := migrateLink(ctx, client, c.Dst, lnk); err != nil {
			logkit.Fatal("migrate link failed", zap.Error(err))
		}
	}
}

func downloadLink(ctx context.Context, client *http.Client, link string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("code not ok, code:%d", rsp.StatusCode)
	}
	raw, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func uploadLink(ctx context.Context, client *http.Client, data []byte, link string, u, p string) error {
	req, err := http.NewRequest(http.MethodPut, link, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.SetBasicAuth(u, p)
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload status code not ok, code:%d", rsp.StatusCode)
	}
	return nil
}

func migrateLink(ctx context.Context, client *http.Client, dst config.Destination, link string) error {
	logger := logutil.GetLogger(ctx).With(zap.String("link", link))
	uri, err := url.Parse(link)
	if err != nil {
		return err
	}
	path := uri.Path
	if len(uri.RawQuery) > 0 {
		return fmt.Errorf("should not contain query, migrate s3 link only")
	}
	data, err := downloadLink(ctx, client, link)
	if err != nil {
		return err
	}
	logger.Info("download link succ", zap.Int("size", len(data)))
	newlink := fmt.Sprintf("%s%s", dst.Host, path)
	if err := uploadLink(ctx, client, data, newlink, dst.AccessKey, dst.SecretKey); err != nil {
		return err
	}
	logger.Info("upload succ", zap.String("dst_link", newlink))
	return nil
}
