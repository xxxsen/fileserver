package client

import "io"

type BeginUploadRequest struct {
	FileSize int64 `json:"file_size"`
}

type BeginUploadResponse struct {
	UploadCtx string `json:"upload_ctx"`
	BlockSize uint32 `json:"block_size"`
}

type PartUploadRequest struct {
	PartID    uint32
	PartMD5   string
	UploadCtx string
	Reader    io.Reader
}

type PartUploadResponse struct {
}

type EndUploadRequest struct {
	UploadCtx string `json:"upload_ctx"`
	FileName  string `json:"file_name"`
}

type EndUploadResponse struct {
	DownKey string `json:"down_key"`
}
