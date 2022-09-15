package file

import (
	"context"
	"fileserver/core"
	"fileserver/dao"
	"fileserver/handler/getter"
	"fileserver/model"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/cache"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

var fileCache, _ = cache.NewLocalCache(20000)

type BasicFileDownloadRequest struct {
	DownKey uint64 `form:"down_key" binding:"required"`
}

func FileDownload(ctx *gin.Context, request interface{}) (int, errs.IError, interface{}) {
	req := request.(*BasicFileDownloadRequest)
	downKey := req.DownKey
	ifileinfo, exist, err := cacheGetFileMeta(ctx, fileCache, downKey, func() (interface{}, bool, error) {
		daoRsp, exist, err := dao.FileInfoDao.GetFile(ctx, &model.GetFileRequest{
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
		return http.StatusOK, errs.Wrap(errs.ErrStorage, "cache get file meta fail", err), nil
	}
	if !exist {
		return http.StatusOK, errs.New(errs.ErrNotFound, "not found file meta"), nil
	}
	fileinfo := ifileinfo.(*model.FileItem)

	fs := getter.MustGetFsClient(ctx)

	rsp, err := fs.FileDownload(ctx, &core.FileDownloadRequest{
		Key:     fileinfo.FileKey,
		Extra:   fileinfo.Extra,
		StartAt: 0,
	})
	if err != nil {
		return http.StatusOK, errs.Wrap(errs.ErrS3, "create download stream fail", err), nil
	}
	defer rsp.Reader.Close()
	contentType := mime.TypeByExtension(filepath.Ext(fileinfo.FileName))
	writer := ctx.Writer
	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", strconv.Quote(fileinfo.FileName)))
	writer.Header().Set("Content-Length", fmt.Sprintf("%d", fileinfo.FileSize))
	writer.Header().Set("Content-Type", contentType)
	sz, err := io.Copy(ctx.Writer, rsp.Reader)
	if err != nil {
		logutil.GetLogger(ctx).With(zap.Error(err), zap.Uint64("key", req.DownKey)).Error("copy stream fail")
		return http.StatusOK, nil, nil
	}
	if sz != int64(fileinfo.FileSize) {
		logutil.GetLogger(ctx).With(zap.Error(err),
			zap.Uint64("key", req.DownKey), zap.Uint64("need_size", fileinfo.FileSize),
			zap.Int64("write_size", sz)).Error("io size not match")
		return http.StatusOK, nil, nil
	}
	return http.StatusOK, nil, nil
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
