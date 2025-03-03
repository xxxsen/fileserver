package file

import (
	"context"
	"fileserver/filesystem"
	"fileserver/proxyutil"
	"fileserver/server/model"
	"fileserver/service"
	"fileserver/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func FileDownload(c *gin.Context, ctx context.Context, request interface{}) {
	req := request.(*model.DownloadFileRequest)
	fileid, err := utils.DecodeFileId(req.DownKey)
	if err != nil {
		proxyutil.Fail(c, http.StatusBadRequest, fmt.Errorf("invalid down key, err:%w", err))
		return
	}
	finfo, ok, err := service.FileService.GetFileInfo(ctx, fileid)
	if err != nil {
		proxyutil.Fail(c, http.StatusInternalServerError, fmt.Errorf("read file info failed, err:%w", err))
		return
	}
	if !ok {
		proxyutil.Fail(c, http.StatusNotFound, fmt.Errorf("file not found"))
		return
	}
	sk := filesystem.NewSeeker(ctx, func(ctx context.Context, blkid int32) (string, error) {
		pinfo, ok, err := service.FileService.GetFilePartInfo(ctx, fileid, blkid)
		if err != nil {
			return "", fmt.Errorf("read file part info failed, err:%w", err)
		}
		if !ok {
			return "", fmt.Errorf("partid:%d not found", blkid)
		}
		return pinfo.FileKey, nil
	}, finfo.FileSize)
	defer sk.Close()
	http.ServeContent(c.Writer, c.Request, strconv.Quote(finfo.FileName), time.Unix(int64(finfo.Ctime), 0), sk)
}
