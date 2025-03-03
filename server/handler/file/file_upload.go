package file

import (
	"context"
	"fileserver/proto/fileserver/fileinfo"
	"fileserver/proxyutil"
	"fileserver/server/model"
	"fileserver/server/stream"
	"fileserver/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
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
	fileid, err := stream.ServeUpload(c, ctx, file, header.Filename, header.Size)
	if err != nil {
		proxyutil.Fail(c, http.StatusInternalServerError, fmt.Errorf("upload file fail, err:%w", err))
		return
	}
	proxyutil.Success(c, &fileinfo.FileUploadResponse{
		DownKey: proto.String(utils.EncodeFileId(fileid)),
	})
}
