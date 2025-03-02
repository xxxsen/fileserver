package s3

import (
	"bytes"
	"fileserver/dao"
	"fileserver/handler/common"
	"fileserver/handler/s3base"
	"fileserver/model"
	"fileserver/utils"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/idgen"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

func smallFileUpload(ctx *gin.Context) (uint64, string, error) {
	md5Base64 := ctx.Request.Header.Get("Content-MD5")
	length := ctx.Request.ContentLength

	raw, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return 0, "", fmt.Errorf("read body fail, err:%w", err)
	}
	checksum, err := utils.Base64Md52HexMd5(md5Base64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid checksum, err:%w", err)
	}
	if len(checksum) == 0 {
		checksum = utils.GetMd5(raw)
	}

	obj, _ := s3base.GetS3Object(ctx)
	name := path.Base(obj)
	uploadRequest := common.CommonUploadContext{
		IDG:    idgen.Default(),
		Reader: bytes.NewReader(raw),
		Size:   length,
		Name:   name,
		Md5Sum: checksum,
	}
	fileid, err := common.Upload(ctx, &uploadRequest)
	if err != nil {
		return 0, "", fmt.Errorf("upload fail, err:%w", err)
	}
	return fileid, checksum, nil
}

func Upload(ctx *gin.Context) {
	caller := smallFileUpload
	length := ctx.Request.ContentLength
	fileid, checksum, err := caller(ctx)
	if err != nil {
		s3base.WriteError(ctx, http.StatusInternalServerError, fmt.Errorf("do file upload fail, err:%w", err))
		return
	}
	bucket, _ := s3base.GetS3Bucket(ctx)
	obj, _ := s3base.GetS3Object(ctx)
	filename := fmt.Sprintf("%s/%s", bucket, obj)
	if _, err := dao.MappingInfoDao.CreateMappingInfo(ctx, &model.CreateMappingInfoRequest{
		Item: &model.MappingInfoItem{
			FileName: filename,
			FileId:   fileid,
		},
	}); err != nil {
		s3base.WriteError(ctx, http.StatusInternalServerError, fmt.Errorf("create mapping fail, err:%w", err))
		return
	}
	ctx.Writer.Header().Set("ETag", `"`+checksum+`"`)
	ctx.Writer.WriteHeader(http.StatusOK)
	logutil.GetLogger(ctx).Info("upload file finish", zap.Int64("size", length), zap.String("bucket", bucket), zap.String("obj", obj))
}
