package file

import (
	"context"
	"fileserver/handler/common"
	"fileserver/proto/fileserver/fileinfo"
	"fileserver/proxyutil"
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

func FileUpload(c *gin.Context, ctx context.Context, request interface{}) {
	req := request.(*BasicFileUploadRequest)
	header := req.File
	file, err := header.Open()
	if err != nil {
		proxyutil.Fail(c, http.StatusBadRequest, fmt.Errorf("open file fail, err:%w", err))
		return
	}
	defer file.Close()
	md5 := req.MD5

	fileid, err := common.Upload(ctx, &common.CommonUploadContext{
		IDG:    idgen.Default(),
		Name:   header.Filename,
		Size:   header.Size,
		Reader: file,
		Md5Sum: md5,
	})
	if err != nil {
		proxyutil.Fail(c, http.StatusInternalServerError, fmt.Errorf("do file upload fail, err:%w", err))
		return

	}
	proxyutil.Success(c, &fileinfo.FileUploadResponse{
		DownKey: proto.String(utils.EncodeFileId(fileid)),
	})
}
