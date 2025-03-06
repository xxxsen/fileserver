package file

import (
	"context"
	"fileserver/filemgr"
	"fileserver/proxyutil"
	"fileserver/server/model"
	"fileserver/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FileDownload(c *gin.Context, ctx context.Context, request interface{}) {
	req := request.(*model.DownloadFileRequest)
	fileid, err := utils.DecodeFileId(req.DownKey)
	if err != nil {
		proxyutil.Fail(c, http.StatusBadRequest, fmt.Errorf("invalid down key, err:%w", err))
		return
	}
	file, err := filemgr.Open(ctx, fileid)
	if err != nil {
		proxyutil.Fail(c, http.StatusInternalServerError, fmt.Errorf("open file failed, err:%w", err))
		return
	}
	defer file.Close()
	finfo, err := filemgr.Stat(ctx, fileid)
	if err != nil {
		proxyutil.Fail(c, http.StatusInternalServerError, fmt.Errorf("stat file failed, err:%w", err))
		return
	}
	http.ServeContent(c.Writer, c.Request, finfo.Name(), finfo.ModTime(), file)
}
