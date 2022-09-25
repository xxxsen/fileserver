package s3

import (
	"fileserver/dao"
	"fileserver/handler/common"
	"fileserver/handler/getter"
	"fileserver/handler/middlewares"
	"fileserver/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

func Download(ctx *gin.Context) {
	bucket, _ := middlewares.GetS3Bucket(ctx)
	obj, _ := middlewares.GetS3Object(ctx)
	filename := fmt.Sprintf("%s/%s", bucket, obj)
	mappingResponse, err := dao.MappingInfoDao.GetMappingInfo(ctx, &model.GetMappingInfoRequest{
		FileName: filename,
	})
	if err != nil {
		WriteError(ctx, http.StatusInternalServerError, errs.Wrap(errs.ErrDatabase, "get mapping info fail", err))
		return
	}
	if mappingResponse.Item == nil {
		WriteError(ctx, http.StatusNotFound, errs.New(errs.ErrNotFound, "data not found"))
		return
	}
	fs := getter.MustGetFsClient(ctx)
	if err := common.Download(ctx, &common.CommonDownloadContext{
		DownKey: mappingResponse.Item.FileId,
		Fs:      fs,
		Dao:     dao.FileInfoDao,
	}); err != nil {
		WriteError(ctx, http.StatusInternalServerError, errs.Wrap(errs.ErrServiceInternal, "do download fail", err))
		return
	}
	logutil.GetLogger(ctx).With(zap.String("bucket", bucket), zap.String("obj", obj)).Info("download file finish")
}
