package bigfile

import (
	"fileserver/core"
	"fileserver/handler/getter"
	"fileserver/proto/fileserver/fileinfo"
	"net/http"

	"github.com/xxxsen/common/errs"

	"github.com/gin-gonic/gin"
)

type PartUploadRequest struct {
	PartId    uint64 `form:"part_id" binding:"required"`
	MD5       string `form:"md5" binding:"required"`
	UploadCtx string `form:"upload_ctx" binding:"required"`
}

func Part(ctx *gin.Context, request interface{}) (int, errs.IError, interface{}) {
	req := request.(*PartUploadRequest)
	var (
		partid     = req.PartId
		md5        = req.MD5
		suploadctx = req.UploadCtx
	)

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		return http.StatusOK, errs.Wrap(errs.ErrParam, "get file fail", err), nil
	}
	defer file.Close()

	fs := getter.MustGetFsClient(ctx)
	_, err = fs.PartFileUpload(ctx, &core.PartFileUploadRequest{
		ReadSeeker: file,
		UploadId:   suploadctx,
		PartId:     partid,
		Size:       header.Size,
		MD5:        md5,
	})
	if err != nil {
		return http.StatusOK, errs.Wrap(errs.ErrS3, "upload part fail", err), nil
	}
	return http.StatusOK, nil, &fileinfo.FileUploadPartResponse{}
}
