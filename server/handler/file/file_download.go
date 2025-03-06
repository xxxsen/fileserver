package file

import (
	"context"
	"fileserver/filemgr"
	"fileserver/proxyutil"
	"fileserver/server/model"
	"fileserver/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FileDownload(c *gin.Context, ctx context.Context, request interface{}) {
	req := request.(*model.DownloadFileRequest)
	fileid, err := utils.DecodeFileId(req.Key)
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
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", strconv.Quote(finfo.Name())))
	http.ServeContent(c.Writer, c.Request, strconv.Quote(finfo.Name()), finfo.ModTime(), file)
}
