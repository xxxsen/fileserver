package s3

import (
	"fileserver/dao"
	"fileserver/handler/common"
	"fileserver/handler/getter"
	"fileserver/handler/s3base"
	"fileserver/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

func Download(ctx *gin.Context) {
	bucket, _ := s3base.GetS3Bucket(ctx)
	obj, _ := s3base.GetS3Object(ctx)
	filename := fmt.Sprintf("%s/%s", bucket, obj)
	mappingResponse, err := dao.MappingInfoDao.GetMappingInfo(ctx, &model.GetMappingInfoRequest{
		FileName: filename,
	})
	if err != nil {
		s3base.WriteError(ctx, http.StatusInternalServerError, errs.Wrap(errs.ErrDatabase, "get mapping info fail", err))
		return
	}
	if mappingResponse.Item == nil {
		s3base.WriteError(ctx, http.StatusNotFound, errs.New(errs.ErrNotFound, "data not found"))
		return
	}
	fs := getter.MustGetFsClient(ctx)
	if err := common.Download(ctx, &common.CommonDownloadContext{
		DownKey: mappingResponse.Item.FileId,
		Fs:      fs,
		Dao:     dao.FileInfoDao,
	}); err != nil {
		s3base.WriteError(ctx, http.StatusInternalServerError, errs.Wrap(errs.ErrServiceInternal, "do download fail", err))
		return
	}
	logutil.GetLogger(ctx).Info("download file finish", zap.String("bucket", bucket), zap.String("obj", obj))
}

func GetBucket(ctx *gin.Context) {
	s3base.SimpleReply(ctx)
}
