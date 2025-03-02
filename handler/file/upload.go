package file

import (
	"fileserver/dao"
	"fileserver/handler/common"
	"fileserver/handler/getter"
	"fileserver/proto/fileserver/fileinfo"
	"fileserver/utils"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/idgen"
	"google.golang.org/protobuf/proto"
)

type BasicFileUploadRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
	MD5  string                `form:"md5"`
}

var ImageUpload = FileUpload
var VideoUpload = FileUpload

func FileUpload(ctx *gin.Context, request interface{}) (int, interface{}, error) {
	req := request.(*BasicFileUploadRequest)
	header := req.File
	file, err := header.Open()
	if err != nil {
		return http.StatusOK, fmt.Errorf("open file fail, err:%w", err), nil
	}
	defer file.Close()
	fs := getter.MustGetFsClient(ctx)
	md5 := req.MD5

	fileid, err := common.Upload(ctx, &common.CommonUploadContext{
		IDG:    idgen.Default(),
		Fs:     fs,
		Dao:    dao.FileInfoDao,
		Name:   header.Filename,
		Size:   header.Size,
		Reader: file,
		Md5Sum: md5,
	})
	if err != nil {
		return http.StatusOK, fmt.Errorf("do file upload fail, err:%w", err), nil
	}
	return http.StatusOK, &fileinfo.FileUploadResponse{
		DownKey: proto.String(utils.EncodeFileId(fileid)),
	}, nil
}
