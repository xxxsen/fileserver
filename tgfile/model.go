package tgfile

import "io"

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
	StType  uint8
}

type FileDownloadResponse struct {
	Reader io.ReadCloser
}
