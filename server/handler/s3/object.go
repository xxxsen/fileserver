package s3

import (
	"fileserver/filemgr"
	"fileserver/proxyutil"
	"fileserver/server/handler/s3/s3base"
	"fmt"
	"net/http"

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
	fid, err := filemgr.ResolveLink(ctx, filename)
	if err != nil {
		s3base.WriteError(c, http.StatusInternalServerError, fmt.Errorf("get mapping info fail, err:%w", err))
		return
	}
	finfo, err := filemgr.Stat(ctx, fid)
	if err != nil {
		s3base.WriteError(c, http.StatusInternalServerError, fmt.Errorf("get file info fail, err:%w", err))
		return
	}
	file, err := filemgr.Open(ctx, fid)
	if err != nil {
		s3base.WriteError(c, http.StatusInternalServerError, fmt.Errorf("open file fail, err:%w", err))
		return
	}
	defer file.Close()
	http.ServeContent(c.Writer, c.Request, finfo.Name(), finfo.ModTime(), file)
}

func UploadObject(c *gin.Context) {
	ctx := c.Request.Context()

	sinfo, ok := proxyutil.GetS3Info(ctx)
	if !ok {
		s3base.WriteError(c, http.StatusBadRequest, fmt.Errorf("no s3 info found"))
		return
	}
	filename := fmt.Sprintf("%s/%s", sinfo.Bucket, sinfo.Object)
	fileid, err := filemgr.Create(ctx, filename, c.Request.ContentLength, c.Request.Body)
	if err != nil {
		s3base.WriteError(c, http.StatusInternalServerError, fmt.Errorf("do file upload fail, err:%w", err))
		return
	}
	if err := filemgr.CreateLink(ctx, filename, fileid); err != nil {
		s3base.WriteError(c, http.StatusInternalServerError, fmt.Errorf("create mapping fail, err:%w", err))
		return
	}
	//TODO: 确认下, 不写etag是否会有问题
	c.Writer.WriteHeader(http.StatusOK)
}
