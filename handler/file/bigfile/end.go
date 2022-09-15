package bigfile

import (
	"fileserver/core"
	"fileserver/dao"
	"fileserver/handler/getter"
	"fileserver/model"
	"fileserver/proto/fileserver/fileinfo"
	"fileserver/utils"
	"net/http"
	"time"

	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/idgen"
	"google.golang.org/protobuf/proto"

	"github.com/gin-gonic/gin"
)

func End(ctx *gin.Context, request interface{}) (int, errs.IError, interface{}) {
	req := request.(*fileinfo.FileUploadEndRequest)
	if len(req.GetHash()) == 0 || len(req.GetFileName()) == 0 ||
		len(req.GetUploadCtx()) == 0 {
		return http.StatusOK, errs.New(errs.ErrParam, "invalid hash/filename/uploadctx"), nil
	}

	fs := getter.MustGetFsClient(ctx)
	rsp, err := fs.FinishFileUpload(ctx, &core.FinishFileUploadRequest{
		UploadId: req.GetUploadCtx(),
		FileMd5:  req.GetHash(),
		FileName: req.GetFileName(),
	})
	if err != nil {
		return http.StatusOK, errs.Wrap(errs.ErrS3, "complete upload fail", err), nil
	}
	fileid := idgen.NextId()
	_, err = dao.FileInfoDao.CreateFile(ctx, &model.CreateFileRequest{
		Item: &model.FileItem{
			FileName:   req.GetFileName(),
			Hash:       req.GetHash(),
			FileSize:   uint64(rsp.FileSize),
			CreateTime: uint64(time.Now().UnixMilli()),
			FileKey:    rsp.Key,
			Extra:      rsp.Extra,
			DownKey:    fileid,
		},
	})
	if err != nil {
		return http.StatusOK, errs.Wrap(errs.ErrDatabase, "write file record fail", err), nil
	}
	return http.StatusOK, nil, &fileinfo.FileUploadEndResponse{
		DownKey: proto.String(utils.EncodeFileId(fileid)),
	}
}
