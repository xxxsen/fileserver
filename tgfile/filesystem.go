package tgfile

import "context"

var defaultFs IFileSystem

type IFileSystem interface {
	StType() uint8
	BlockSize() int64
	MaxFileSize() int64
	FileUpload(ctx context.Context, uctx *FileUploadRequest) (*FileUploadResponse, error)
	FileDownload(ctx context.Context, fctx *FileDownloadRequest) (*FileDownloadResponse, error)
	BeginFileUpload(ctx context.Context, fctx *BeginFileUploadRequest) (*BeginFileUploadResponse, error)
	PartFileUpload(ctx context.Context, pctx *PartFileUploadRequest) (*PartFileUploadResponse, error)
	FinishFileUpload(ctx context.Context, fctx *FinishFileUploadRequest) (*FinishFileUploadResponse, error)
}

func SetFileSystem(fs IFileSystem) {
	defaultFs = fs
}

func GetFileSystem() IFileSystem {
	return defaultFs
}

func StType() uint8 {
	return defaultFs.StType()
}

func BlockSize() int64 {
	return defaultFs.BlockSize()
}

func MaxFileSize() int64 {
	return defaultFs.MaxFileSize()
}

func FileUpload(ctx context.Context, uctx *FileUploadRequest) (*FileUploadResponse, error) {
	return defaultFs.FileUpload(ctx, uctx)
}

func FileDownload(ctx context.Context, fctx *FileDownloadRequest) (*FileDownloadResponse, error) {
	return defaultFs.FileDownload(ctx, fctx)
}

func BeginFileUpload(ctx context.Context, fctx *BeginFileUploadRequest) (*BeginFileUploadResponse, error) {
	return defaultFs.BeginFileUpload(ctx, fctx)
}

func PartFileUpload(ctx context.Context, pctx *PartFileUploadRequest) (*PartFileUploadResponse, error) {
	return defaultFs.PartFileUpload(ctx, pctx)
}

func FinishFileUpload(ctx context.Context, fctx *FinishFileUploadRequest) (*FinishFileUploadResponse, error) {
	return defaultFs.FinishFileUpload(ctx, fctx)
}
