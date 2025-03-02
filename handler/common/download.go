package common

import (
	"context"
	"fileserver/core"
	"fileserver/dao"
	"fileserver/model"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/qingstor/go-mime"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/cache"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

const (
	miniEnableRangeDownloadSize = 1 //enable range download by default
)

var fileCache, _ = cache.NewLocalCache(20000)

type CommonDownloadContext struct {
	DownKey uint64
	Fs      core.IFsCore
	Dao     dao.FileInfoService
}

func streamDownload(ctx *gin.Context, downKey uint64, fs core.IFsCore, fileinfo *model.FileItem) error {
	rsp, err := fs.FileDownload(ctx, &core.FileDownloadRequest{
		Key:     fileinfo.FileKey,
		Extra:   fileinfo.Extra,
		StartAt: 0,
		StType:  fileinfo.StType,
	})
	if err != nil {
		return fmt.Errorf("create download stream fail, err:%w", err)
	}
	defer rsp.Reader.Close()
	contentType := mime.DetectFilePath(fileinfo.FileName)
	writer := ctx.Writer
	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", strconv.Quote(fileinfo.FileName)))
	writer.Header().Set("Content-Length", fmt.Sprintf("%d", fileinfo.FileSize))
	writer.Header().Set("Content-Type", contentType)
	sz, err := io.Copy(writer, rsp.Reader)
	//shall not return err when write data to client
	if err != nil {
		logutil.GetLogger(ctx).With(zap.Error(err), zap.Uint64("key", downKey)).Error("copy stream fail")
		return nil
	}
	if sz != int64(fileinfo.FileSize) {
		logutil.GetLogger(ctx).Error("io size not match", zap.Error(err),
			zap.Uint64("key", downKey), zap.Uint64("need_size", fileinfo.FileSize),
			zap.Int64("write_size", sz))
		return nil
	}
	return nil
}

func cacheGetFileMeta(ctx context.Context, c cache.ICache, key interface{},
	cb func() (interface{}, bool, error)) (interface{}, bool, error) {

	ival, exist, _ := c.Get(ctx, key)
	if exist {
		return ival, true, nil
	}
	val, exist, err := cb()
	if err != nil {
		return nil, false, err
	}
	if exist {
		c.Set(ctx, key, val, 10*time.Minute)
	}
	return val, exist, nil
}

func Download(ctx *gin.Context, fctx *CommonDownloadContext) error {
	downKey := fctx.DownKey
	fs := fctx.Fs

	ifileinfo, exist, err := cacheGetFileMeta(ctx, fileCache, downKey, func() (interface{}, bool, error) {
		daoRsp, exist, err := fctx.Dao.GetFile(ctx, &model.GetFileRequest{
			DownKey: downKey,
		})
		if err != nil {
			return nil, false, err
		}
		if !exist {
			return nil, false, nil
		}
		return daoRsp.Item, true, nil
	})
	if err != nil {
		return fmt.Errorf("cache get file meta fail, err:%w", err)
	}
	if !exist {
		return errs.New(errs.ErrNotFound, "not found file meta")
	}
	fileinfo := ifileinfo.(*model.FileItem)

	if r := ctx.GetHeader("range"); len(r) == 0 || fileinfo.FileSize < miniEnableRangeDownloadSize { //filesize < 200MB will not enable range download
		return streamDownload(ctx, downKey, fs, fileinfo)
	}

	file := core.NewSeeker(ctx, fs, int64(fileinfo.FileSize), fileinfo.FileKey, fileinfo.Extra, fileinfo.StType)
	defer file.Close()
	http.ServeContent(ctx.Writer, ctx.Request, strconv.Quote(fileinfo.FileName), time.Unix(int64(fileinfo.CreateTime), 0), file)
	return nil
}
