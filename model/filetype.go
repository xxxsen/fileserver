package model

type FileType uint32

const (
	FileTypeVideo FileType = 1
	FileTypeImage FileType = 2
	FileTypeFile  FileType = 3
	FileTypeAny   FileType = 10
)
