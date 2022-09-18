package file

import (
	"fileserver/core"
	"fileserver/dao"
	"fileserver/handler/getter"
	"fileserver/model"
	"fileserver/proto/fileserver/fileinfo"
	"fileserver/utils"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/idgen"
	"google.golang.org/protobuf/proto"
)

type BasicFileUploadRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
	MD5  string                `form:"md5"`
}

var ImageUpload = FileUpload
var VideoUpload = FileUpload

func FileUpload(ctx *gin.Context, request interface{}) (int, errs.IError, interface{}) {
	req := request.(*BasicFileUploadRequest)
	header := req.File
	file, err := header.Open()
	if err != nil {
		return http.StatusOK, errs.Wrap(errs.ErrParam, "open file fail", err), nil
	}
	defer file.Close()
	fs := getter.MustGetFsClient(ctx)
	md5 := req.MD5
	if header.Size > fs.MaxFileSize() {
		return http.StatusOK, errs.New(errs.ErrParam, "file size out of limit, should less than:%d", fs.MaxFileSize()), nil
	}
	if header.Size == 0 {
		return http.StatusOK, errs.New(errs.ErrParam, "empty file"), nil
	}

	rsp, err := fs.FileUpload(ctx, &core.FileUploadRequest{
		ReadSeeker: file,
		Size:       header.Size,
		MD5:        md5,
	})
	if err != nil {
		return http.StatusOK, errs.Wrap(errs.ErrStorage, "upload file fail", err), nil
	}
	fileid := idgen.NextId()
	if _, err := dao.FileInfoDao.CreateFile(ctx, &model.CreateFileRequest{
		//TODO: write checksum here
		Item: &model.FileItem{
			FileName:   header.Filename,
			Hash:       rsp.CheckSum,
			FileSize:   uint64(header.Size),
			CreateTime: uint64(time.Now().UnixMilli()),
			DownKey:    fileid,
			FileKey:    rsp.Key,
			Extra:      rsp.Extra,
		},
	}); err != nil {
		return http.StatusOK, errs.Wrap(errs.ErrDatabase, "insert image to db fail", err), nil
	}
	return http.StatusOK, nil, &fileinfo.FileUploadResponse{
		DownKey: proto.String(utils.EncodeFileId(fileid)),
	}
}
