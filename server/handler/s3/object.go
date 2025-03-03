package s3

import (
	"fileserver/proxyutil"
	"fileserver/server/handler/s3/s3base"
	"fileserver/server/stream"
	"fileserver/service"
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

func DownloadObject(c *gin.Context) {
	ctx := c.Request.Context()

	sinfo, ok := proxyutil.GetS3Info(ctx)
	if !ok {
		s3base.WriteError(c, http.StatusBadRequest, fmt.Errorf("no s3 info found"))
		return
	}
	filename := fmt.Sprintf("%s/%s", sinfo.Bucket, sinfo.Object)
	fid, ok, err := service.FileMappingService.GetFileMapping(ctx, filename)
	if err != nil {
		s3base.WriteError(c, http.StatusInternalServerError, fmt.Errorf("get mapping info fail, err:%w", err))
		return
	}
	if !ok {
		s3base.WriteError(c, http.StatusNotFound, fmt.Errorf("data not found"))
		return
	}
	stream.ServeDownload(c, ctx, fid)
}

func UploadObject(c *gin.Context) {
	ctx := c.Request.Context()

	sinfo, ok := proxyutil.GetS3Info(ctx)
	if !ok {
		s3base.WriteError(c, http.StatusBadRequest, fmt.Errorf("no s3 info found"))
		return
	}
	filename := fmt.Sprintf("%s/%s", sinfo.Bucket, sinfo.Object)
	name := fmt.Sprintf("s3:%s", path.Base(filename))
	fileid, err := stream.ServeUpload(c, ctx, c.Request.Body, name, c.Request.ContentLength)
	if err != nil {
		s3base.WriteError(c, http.StatusInternalServerError, fmt.Errorf("do file upload fail, err:%w", err))
		return
	}
	if err := service.FileMappingService.CreateFileMapping(ctx, filename, fileid); err != nil {
		s3base.WriteError(c, http.StatusInternalServerError, fmt.Errorf("create mapping fail, err:%w", err))
		return
	}
	//TODO: 确认下, 不写etag是否会有问题
	c.Writer.WriteHeader(http.StatusOK)
}
