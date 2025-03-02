package common

import (
	"context"
	"fileserver/core"
	"fileserver/dao"
	"fileserver/model"
	"fmt"
	"io"
	"time"

	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/idgen"
)

type CommonUploadContext struct {
	IDG idgen.IDGenerator

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
		fs   = core.GetFsCore()
	)
	if size > fs.MaxFileSize() {
		return 0, fmt.Errorf("file size out of limit, should less than:%d", fs.MaxFileSize())
	}
	if size == 0 {
		return 0, errs.New(errs.ErrParam, "empty file")
	}

	rsp, err := fs.FileUpload(ctx, &core.FileUploadRequest{
		ReadSeeker: file,
		Size:       size,
		MD5:        md5,
	})
	if err != nil {
		return 0, fmt.Errorf("upload file fail, err:%w", err)
	}
	fileid := fctx.IDG.NextId()
	if _, err := dao.FileInfoDao.CreateFile(ctx, &model.CreateFileRequest{
		Item: &model.FileItem{
			FileName:   name,
			Hash:       rsp.CheckSum,
			FileSize:   uint64(size),
			CreateTime: uint64(time.Now().UnixMilli()),
			DownKey:    fileid,
			FileKey:    rsp.Key,
			Extra:      rsp.Extra,
			StType:     fs.StType(),
		},
	}); err != nil {
		return 0, fmt.Errorf("insert file to db fail, err:%w", err)
	}
	return fileid, nil
}
