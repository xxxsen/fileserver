package file

import (
	"fileserver/dao"
	"fileserver/handler/common"
	"fileserver/handler/getter"
	"fileserver/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BasicFileDownloadRequest struct {
	DownKey string `form:"down_key" binding:"required"`
}

func FileDownload(ctx *gin.Context, request interface{}) (int, interface{}, error) {
	req := request.(*BasicFileDownloadRequest)
	strDownKey := req.DownKey
	downKey, err := utils.DecodeFileId(strDownKey)
	if err != nil {
		return http.StatusOK, nil, fmt.Errorf("invalid down key, err:%w", err)
	}

	if err := common.Download(ctx, &common.CommonDownloadContext{
		DownKey: downKey,
		Fs:      getter.MustGetFsClient(ctx),
		Dao:     dao.FileInfoDao,
	}); err != nil {
		return http.StatusOK, nil, fmt.Errorf("do download fail, err:%w", err)
	}
	return http.StatusOK, nil, nil
}
