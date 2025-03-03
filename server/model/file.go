package model

import "mime/multipart"

type DownloadFileRequest struct {
	DownKey string `form:"down_key" binding:"required"`
}

type UploadFileRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type GetFileInfoRequest struct {
	DownKey []string `json:"down_key"`
}

type GetFileInfoItem struct {
	DownKey       string `json:"down_key"`
	FileName      string `json:"file_name"`
	FileSize      int64  `json:"file_size"`
	Ctime         int64  `json:"ctime"`
	Mtime         int64  `json:"mtime"`
	FilePartCount int32  `json:"file_part_count"`
}

type GetFileInfoResponse struct {
	List []*GetFileInfoItem `json:"list"`
}
