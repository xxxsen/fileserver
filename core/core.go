package core

import (
	"context"
	"io"
)

type FileUploadRequest struct {
	ReadSeeker io.ReadSeeker
	Size       int64
	MD5        string
}

type FileUploadResponse struct {
	Key      string
	Extra    []byte
	CheckSum string
}

type BeginFileUploadRequest struct {
	FileSize int64
}

type BeginFileUploadResponse struct {
	UploadID string
}

type PartFileUploadRequest struct {
	ReadSeeker io.ReadSeeker
	UploadId   string
	PartId     uint64
	Size       int64
	MD5        string
}

type PartFileUploadResponse struct {
}

type FinishFileUploadRequest struct {
	UploadId string
	FileName string
}

type FinishFileUploadResponse struct {
	Key      string
	Extra    []byte
	FileSize int64
	CheckSum string
}

type FileDownloadRequest struct {
	Key     string
	Extra   []byte
	StartAt int64
}

type FileDownloadResponse struct {
	Reader io.ReadCloser
}

type IFsCore interface {
	BlockSize() int64
	MaxFileSize() int64
	FileUpload(ctx context.Context, uctx *FileUploadRequest) (*FileUploadResponse, error)
	FileDownload(ctx context.Context, fctx *FileDownloadRequest) (*FileDownloadResponse, error)
	BeginFileUpload(ctx context.Context, fctx *BeginFileUploadRequest) (*BeginFileUploadResponse, error)
	PartFileUpload(ctx context.Context, pctx *PartFileUploadRequest) (*PartFileUploadResponse, error)
	FinishFileUpload(ctx context.Context, fctx *FinishFileUploadRequest) (*FinishFileUploadResponse, error)
}
