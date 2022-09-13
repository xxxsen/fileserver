package bigfile

import (
	"fileserver/constants"
	"fileserver/dao"
	"fileserver/model"
	"fileserver/proto/fileserver/fileinfo"
	"fileserver/utils"
	"net/http"
	"time"

	"github.com/xxxsen/common/s3"

	"github.com/xxxsen/common/errs"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func End(ctx *gin.Context, request interface{}) (int, errs.IError, interface{}) {
	req := request.(*fileinfo.FileUploadEndRequest)
	if len(req.GetHash()) == 0 || len(req.GetFileName()) == 0 ||
		len(req.GetUploadCtx()) == 0 {
		return http.StatusOK, errs.New(errs.ErrParam, "invalid hash/filename/uploadctx"), nil
	}
	uploadctx, err := utils.DecodeUploadID(req.GetUploadCtx())
	if err != nil {
		return http.StatusOK, errs.Wrap(errs.ErrParam, "decode upload ctx fail", err), nil
	}
	maxpartsz := utils.CalcFileBlockCount(uploadctx.GetFileSize(), constants.BlockSize)
	err = s3.Client.EndUpload(ctx, *uploadctx.DownKey, uploadctx.GetUploadId(), maxpartsz)
	if err != nil {
		return http.StatusOK, errs.Wrap(errs.ErrS3, "complete upload fail", err), nil
	}
	_, err = dao.FileInfoDao.CreateFile(ctx, &model.CreateFileRequest{
		Item: &model.FileItem{
			FileName:   req.GetFileName(),
			Hash:       req.GetHash(),
			FileSize:   uploadctx.GetFileSize(),
			CreateTime: uint64(time.Now().UnixMilli()),
			DownKey:    uploadctx.GetDownKey(),
		},
	})
	if err != nil {
		return http.StatusOK, errs.Wrap(errs.ErrDatabase, "write file record fail", err), nil
	}
	return http.StatusOK, nil, &fileinfo.FileUploadEndResponse{
		DownKey: proto.String(uploadctx.GetDownKey()),
	}
}
