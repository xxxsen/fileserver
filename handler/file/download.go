package file

import (
	"fileserver/dao"
	"fileserver/handler/common"
	"fileserver/handler/getter"
	"fileserver/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/errs"
)

type BasicFileDownloadRequest struct {
	DownKey string `form:"down_key" binding:"required"`
}

func FileDownload(ctx *gin.Context, request interface{}) (int, errs.IError, interface{}) {
	req := request.(*BasicFileDownloadRequest)
	strDownKey := req.DownKey
	downKey, err := utils.DecodeFileId(strDownKey)
	if err != nil {
		return http.StatusOK, errs.Wrap(errs.ErrParam, "invalid down key", err), nil
	}

	if err := common.Download(ctx, &common.CommonDownloadContext{
		DownKey: downKey,
		Fs:      getter.MustGetFsClient(ctx),
		Dao:     dao.FileInfoDao,
	}); err != nil {
		return http.StatusOK, errs.Wrap(errs.ErrServiceInternal, "do download fail", err), nil
	}
	return http.StatusOK, nil, nil
}
