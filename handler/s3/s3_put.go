package s3

import (
	"bytes"
	"fileserver/dao"
	"fileserver/handler/common"
	"fileserver/handler/getter"
	"fileserver/handler/s3base"
	"fileserver/model"
	"fileserver/utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/idgen"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

const (
	maxS3UploadFileLimit = 10 * 1024 * 1024 //
)

func Upload(ctx *gin.Context) {
	md5Base64 := ctx.Request.Header.Get("Content-MD5")
	length := ctx.Request.ContentLength

	if length > maxS3UploadFileLimit {
		s3base.WriteError(ctx, http.StatusBadRequest, errs.New(errs.ErrParam, "size out of limit, s3 file size should less than:%d", maxS3UploadFileLimit))
		return
	}
	raw, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		s3base.WriteError(ctx, http.StatusBadRequest, errs.Wrap(errs.ErrIO, "read body fail", err))
		return
	}

	checksum, err := utils.Base64Md52HexMd5(md5Base64)
	if err != nil {
		s3base.WriteError(ctx, http.StatusBadRequest, errs.Wrap(errs.ErrParam, "invalid checksum", err))
		return
	}
	if len(checksum) == 0 {
		checksum = utils.GetMd5(raw)
	}

	fs := getter.MustGetFsClient(ctx)
	bucket, _ := s3base.GetS3Bucket(ctx)
	obj, _ := s3base.GetS3Object(ctx)
	name := path.Base(obj)
	uploadRequest := common.CommonUploadContext{
		IDG:    idgen.Default(),
		Fs:     fs,
		Dao:    dao.FileInfoDao,
		Reader: bytes.NewReader(raw),
		Size:   length,
		Name:   name,
		Md5Sum: checksum,
	}
	fileid, err := common.Upload(ctx, &uploadRequest)
	if err != nil {
		s3base.WriteError(ctx, http.StatusInternalServerError, errs.Wrap(errs.ErrIO, "upload fail", err))
		return
	}
	filename := fmt.Sprintf("%s/%s", bucket, obj)
	if _, err := dao.MappingInfoDao.CreateMappingInfo(ctx, &model.CreateMappingInfoRequest{
		Item: &model.MappingInfoItem{
			FileName: filename,
			FileId:   fileid,
		},
	}); err != nil {
		s3base.WriteError(ctx, http.StatusInternalServerError, errs.Wrap(errs.ErrDatabase, "create mapping info fail", err))
		return
	}
	ctx.Writer.Header().Set("ETag", `"`+checksum+`"`)
	ctx.Writer.WriteHeader(http.StatusOK)
	logutil.GetLogger(ctx).With(zap.String("bucket", bucket), zap.String("obj", obj)).Info("upload file finish")
}

func S3Put(ctx *gin.Context) {
	_, exist := s3base.GetS3Object(ctx)
	if !exist {
		s3base.WriteError(ctx, http.StatusInternalServerError, errs.New(errs.ErrParam, "no file found"))
		return
	}
	Upload(ctx)
}
