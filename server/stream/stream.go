package stream

import (
	"context"
	"fileserver/filesystem"
	"fileserver/proxyutil"
	"fileserver/service"
	"fileserver/utils"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type OnUploadSuccCallback func(ctx context.Context, fileid uint64) error

func ServeDownload(c *gin.Context, ctx context.Context, fileid uint64) {
	finfo, ok, err := service.FileService.GetFileInfo(ctx, fileid)
	if err != nil {
		proxyutil.Fail(c, http.StatusInternalServerError, fmt.Errorf("read file info failed, err:%w", err))
		return
	}
	if !ok {
		proxyutil.Fail(c, http.StatusNotFound, fmt.Errorf("file not found"))
		return
	}
	sk := filesystem.NewSeeker(ctx, func(ctx context.Context, blkid int32) (string, error) {
		pinfo, ok, err := service.FileService.GetFilePartInfo(ctx, fileid, blkid)
		if err != nil {
			return "", fmt.Errorf("read file part info failed, err:%w", err)
		}
		if !ok {
			return "", fmt.Errorf("partid:%d not found", blkid)
		}
		return pinfo.FileKey, nil
	}, finfo.FileSize)
	defer sk.Close()
	http.ServeContent(c.Writer, c.Request, strconv.Quote(finfo.FileName), time.Unix(int64(finfo.Ctime), 0), sk)
}

func ServeUpload(c *gin.Context, ctx context.Context, reader io.Reader, filename string, fsize int64) (uint64, error) {
	blkcnt := utils.CalcFileBlockCount(uint64(fsize), uint64(filesystem.MaxFileSize()))
	fileid, err := service.FileService.CreateFileDraft(ctx, filename, fsize, int32(blkcnt))
	if err != nil {
		return 0, fmt.Errorf("create file draft failed, err:%w", err)
	}
	for i := 0; i < blkcnt; i++ {
		r := io.LimitReader(reader, filesystem.MaxFileSize())
		fileKey, err := filesystem.Upload(ctx, r)
		if err != nil {
			return 0, fmt.Errorf("upload part failed, err:%w", err)
		}
		if err := service.FileService.CreateFilePart(ctx, fileid, int32(i), fileKey); err != nil {
			return 0, fmt.Errorf("create part record failed, err:%w", err)
		}
	}

	if err := service.FileService.FinishCreateFile(ctx, fileid); err != nil {
		return 0, fmt.Errorf("finish create file failed, err:%w", err)
	}
	return fileid, nil
}
