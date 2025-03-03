package file

import (
	"context"
	"fileserver/filesystem"
	"fileserver/proto/fileserver/fileinfo"
	"fileserver/proxyutil"
	"fileserver/server/model"
	"fileserver/service"
	"fileserver/utils"
	"fmt"
	"io"
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
	blkcnt := utils.CalcFileBlockCount(uint64(header.Size), uint64(filesystem.MaxFileSize()))
	fileid, err := service.FileService.CreateFileDraft(ctx, header.Filename, header.Size, int32(blkcnt))
	if err != nil {
		proxyutil.Fail(c, http.StatusInternalServerError, fmt.Errorf("create file draft failed, err:%w", err))
		return
	}
	for i := 0; i < blkcnt; i++ {
		r := io.LimitReader(file, filesystem.MaxFileSize())
		fileKey, err := filesystem.Upload(ctx, r)
		if err != nil {
			proxyutil.Fail(c, http.StatusInternalServerError, fmt.Errorf("upload part failed, err:%w", err))
			return
		}
		if err := service.FileService.CreateFilePart(ctx, fileid, int32(i), fileKey); err != nil {
			proxyutil.Fail(c, http.StatusInternalServerError, fmt.Errorf("create part record failed, err:%w", err))
			return
		}
	}

	if err := service.FileService.FinishCreateFile(ctx, fileid); err != nil {
		proxyutil.Fail(c, http.StatusInternalServerError, fmt.Errorf("finish create file failed, err:%w", err))
		return
	}
	proxyutil.Success(c, &fileinfo.FileUploadResponse{
		DownKey: proto.String(utils.EncodeFileId(fileid)),
	})
}
