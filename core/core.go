package core

import (
	"context"
	"io"
)

type FileUploadContext struct {
	ReadSeeker io.ReadSeeker
	Size       int64
	MD5        string
}

type FileMeta struct {
	Name string
	Size int64
}

type BeginFileUploadContext struct {
	FileSize int64
}

type PartFileUploadContext struct {
	ReadSeeker io.ReadSeeker
	UploadId   string
	PartId     uint64
	Size       int64
	MD5        string
}

type FinishFileUploadContext struct {
	UploadId string
	FileMd5  string
	FileName string
}

type FileDownloadContext struct {
	Key   string
	Range *string
}

type IFsCore interface {
	BlockSize() int64
	MaxFileSize() int64
	FileUpload(ctx context.Context, uctx *FileUploadContext) (string, error)
	FileDownload(ctx context.Context, fctx *FileDownloadContext) (io.ReadCloser, error)
	BeginFileUpload(ctx context.Context, fctx *BeginFileUploadContext) (string, error)
	PartFileUpload(ctx context.Context, pctx *PartFileUploadContext) error
	FinishFileUpload(ctx context.Context, fctx *FinishFileUploadContext) (string, error)
}
