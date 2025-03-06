package file

import (
	"context"
	"fileserver/filemgr"
	"fileserver/proxyutil"
	"fileserver/server/model"
	"fileserver/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FileUpload(c *gin.Context, ctx context.Context, request interface{}) {
	req := request.(*model.UploadFileRequest)
	header := req.File
	file, err := header.Open()
	if err != nil {
		proxyutil.Fail(c, http.StatusBadRequest, fmt.Errorf("open file fail, err:%w", err))
		return
	}
	defer file.Close()
	fileid, err := filemgr.Create(ctx, header.Filename, header.Size, file)
	if err != nil {
		proxyutil.Fail(c, http.StatusInternalServerError, fmt.Errorf("upload file fail, err:%w", err))
		return
	}
	proxyutil.Success(c, &model.UploadFileResponse{
		Key: utils.EncodeFileId(fileid),
	})
}
