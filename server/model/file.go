package model

import "mime/multipart"

type DownloadFileRequest struct {
	Key string `form:"key" binding:"required"`
}

type UploadFileRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type UploadFileResponse struct {
	Key string `json:"key"`
}

type GetFileInfoRequest struct {
	Key string `form:"key"  binding:"required"`
}

type GetFileInfoItem struct {
	Key           string `json:"key"`
	Exist         bool   `json:"exist"`
	FileName      string `json:"file_name"`
	FileSize      int64  `json:"file_size"`
	Ctime         int64  `json:"ctime"`
	Mtime         int64  `json:"mtime"`
	FilePartCount int32  `json:"file_part_count"`
}

type GetFileInfoResponse struct {
	Item *GetFileInfoItem `json:"item"`
}
