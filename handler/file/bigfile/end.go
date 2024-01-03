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

func End(ctx *gin.Context, request interface{}) (int, interface{}, error) {
	req := request.(*fileinfo.FileUploadEndRequest)
	if len(req.GetFileName()) == 0 ||
		len(req.GetUploadCtx()) == 0 {
		return http.StatusOK, nil, errs.New(errs.ErrParam, "invalid filename/uploadctx")
	}

	fs := getter.MustGetFsClient(ctx)
	rsp, err := fs.FinishFileUpload(ctx, &core.FinishFileUploadRequest{
		UploadId: req.GetUploadCtx(),
		FileName: req.GetFileName(),
	})
	if err != nil {
		return http.StatusOK, nil, errs.Wrap(errs.ErrS3, "complete upload fail", err)
	}
	fileid := idgen.NextId()
	_, err = dao.FileInfoDao.CreateFile(ctx, &model.CreateFileRequest{
		Item: &model.FileItem{
			FileName:   req.GetFileName(),
			Hash:       rsp.CheckSum,
			FileSize:   uint64(rsp.FileSize),
			CreateTime: uint64(time.Now().UnixMilli()),
			FileKey:    rsp.Key,
			Extra:      rsp.Extra,
			DownKey:    fileid,
			StType:     fs.StType(),
		},
	})
	if err != nil {
		return http.StatusOK, nil, errs.Wrap(errs.ErrDatabase, "write file record fail", err)
	}
	return http.StatusOK, &fileinfo.FileUploadEndResponse{
		DownKey: proto.String(utils.EncodeFileId(fileid)),
	}, nil
}
