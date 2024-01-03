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

func Begin(ctx *gin.Context, request interface{}) (int, interface{}, error) {
	fs := getter.MustGetFsClient(ctx)
	req := request.(*fileinfo.FileUploadBeginRequest)
	if req.GetFileSize() == 0 {
		return http.StatusOK, nil, errs.New(errs.ErrParam, "zero size file")
	}
	if req.GetFileSize() > uint64(fs.MaxFileSize()) {
		return http.StatusOK, nil, errs.New(errs.ErrParam, "file size out of limit")
	}
	uploadRsp, err := fs.BeginFileUpload(ctx, &core.BeginFileUploadRequest{
		FileSize: int64(req.GetFileSize()),
	})
	if err != nil {
		return http.StatusOK, nil, errs.Wrap(errs.ErrStorage, "begin upload fail", err)
	}
	return http.StatusOK, &fileinfo.FileUploadBeginResponse{
		UploadCtx: proto.String(uploadRsp.UploadID),
		BlockSize: proto.Uint32(uint32(fs.BlockSize())),
	}, nil
}
