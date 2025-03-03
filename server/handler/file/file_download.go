package file

import (
	"context"
	"fileserver/proxyutil"
	"fileserver/server/model"
	"fileserver/server/stream"
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
	stream.ServeDownload(c, ctx, fileid)
}
