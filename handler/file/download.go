package file

import (
	"context"
	"fileserver/handler/common"
	"fileserver/proxyutil"
	"fileserver/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BasicFileDownloadRequest struct {
	DownKey string `form:"down_key" binding:"required"`
}

func FileDownload(c *gin.Context, ctx context.Context, request interface{}) {
	req := request.(*BasicFileDownloadRequest)
	strDownKey := req.DownKey
	downKey, err := utils.DecodeFileId(strDownKey)
	if err != nil {
		proxyutil.Fail(c, http.StatusBadRequest, fmt.Errorf("invalid down key, err:%w", err))
		return
	}

	if err := common.Download(c, &common.CommonDownloadContext{
		DownKey: downKey,
	}); err != nil {
		//TODO: 优化这里的逻辑
		proxyutil.Fail(c, http.StatusInternalServerError, fmt.Errorf("do download fail, err:%w", err))
		return
	}
}
