package bigfile

import (
	"fileserver/core"
	"fileserver/handler/getter"
	"fileserver/proto/fileserver/fileinfo"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/errs"
	"google.golang.org/protobuf/proto"
)

func Begin(ctx *gin.Context, request interface{}) (int, errs.IError, interface{}) {
	fs := getter.MustGetFsClient(ctx)
	req := request.(*fileinfo.FileUploadBeginRequest)
	if req.GetFileSize() == 0 {
		return http.StatusOK, errs.New(errs.ErrParam, "zero size file"), nil
	}
	if req.GetFileSize() > uint64(fs.MaxFileSize()) {
		return http.StatusOK, errs.New(errs.ErrParam, "file size out of limit"), nil
	}
	uploadRsp, err := fs.BeginFileUpload(ctx, &core.BeginFileUploadRequest{
		FileSize: int64(req.GetFileSize()),
	})
	if err != nil {
		return http.StatusOK, errs.Wrap(errs.ErrStorage, "begin upload fail", err), nil
	}
	return http.StatusOK, nil, &fileinfo.FileUploadBeginResponse{
		UploadCtx: proto.String(uploadRsp.UploadID),
		BlockSize: proto.Uint32(uint32(fs.BlockSize())),
	}
}
