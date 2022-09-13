package bigfile

import (
	"fileserver/constants"
	"fileserver/model"
	"fileserver/proto/fileserver/fileinfo"
	"fileserver/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/idgen"
	"github.com/xxxsen/common/s3"
	"google.golang.org/protobuf/proto"
)

func Begin(ctx *gin.Context, request interface{}) (int, errs.IError, interface{}) {
	req := request.(*fileinfo.FileUploadBeginRequest)
	if req.GetFileSize() == 0 {
		return http.StatusOK, errs.New(errs.ErrParam, "zero size file"), nil
	}
	if req.GetFileSize() > constants.MaxFileSize {
		return http.StatusOK, errs.New(errs.ErrParam, "file size out of limit"), nil
	}
	downkey := fmt.Sprintf("%d_%s", model.FileTypeAny, utils.EncodeFileId(idgen.NextId()))
	key, err := s3.Client.BeginUpload(ctx, downkey)
	if err != nil {
		return http.StatusOK, errs.Wrap(errs.ErrS3, "begin upload fail", err), nil
	}
	uploadidctx := &fileinfo.UploadIdCtx{
		FileSize: req.FileSize,
		UploadId: proto.String(key),
		DownKey:  proto.String(downkey),
	}
	uploadctx, err := utils.EncodeUploadID(uploadidctx)
	if err != nil {
		return http.StatusOK, errs.Wrap(errs.ErrServiceInternal, "build upload id fail", err), nil
	}
	return http.StatusOK, nil, &fileinfo.FileUploadBeginResponse{
		UploadCtx: proto.String(uploadctx),
		BlockSize: proto.Uint32(uint32(constants.BlockSize)),
	}
}
