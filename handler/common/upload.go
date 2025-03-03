package common

import (
	"context"
	"fileserver/dao"
	"fileserver/model"
	"fileserver/tgfile"
	"fmt"
	"io"
	"time"

	"github.com/xxxsen/common/idgen"
)

type CommonUploadContext struct {
	Name   string
	Size   int64
	Reader io.ReadSeeker
	Md5Sum string
}

func Upload(ctx context.Context, fctx *CommonUploadContext) (uint64, error) {
	var (
		file = fctx.Reader
		md5  = fctx.Md5Sum
		size = fctx.Size
		name = fctx.Name
	)
	if size > tgfile.MaxFileSize() {
		return 0, fmt.Errorf("file size out of limit, should less than:%d", tgfile.MaxFileSize())
	}
	if size == 0 {
		return 0, fmt.Errorf("empty file")
	}

	rsp, err := tgfile.FileUpload(ctx, &tgfile.FileUploadRequest{
		ReadSeeker: file,
		Size:       size,
		MD5:        md5,
	})
	if err != nil {
		return 0, fmt.Errorf("upload file fail, err:%w", err)
	}
	fileid := idgen.NextId()
	if _, err := dao.FileInfoDao.CreateFile(ctx, &model.CreateFileRequest{
		Item: &model.FileItem{
			FileName:   name,
			Hash:       rsp.CheckSum,
			FileSize:   uint64(size),
			CreateTime: uint64(time.Now().UnixMilli()),
			DownKey:    fileid,
			FileKey:    rsp.Key,
			Extra:      rsp.Extra,
			StType:     tgfile.StType(),
		},
	}); err != nil {
		return 0, fmt.Errorf("insert file to db fail, err:%w", err)
	}
	return fileid, nil
}
