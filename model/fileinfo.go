package model

type CreateFileRequest struct {
	Item *FileItem
}

type CreateFileResponse struct {
}

type GetFileRequest struct {
	DownKey uint64
}

type FileItem struct {
	Id         uint64
	FileName   string
	Hash       string
	FileSize   uint64
	CreateTime uint64
	DownKey    uint64
	Extra      []byte
	FileKey    string
}

type GetFileResponse struct {
	Item *FileItem
}

type ListFileQuery struct {
	ID      []uint64
	DownKey []uint64
}

type ListFileRequest struct {
	Query  *ListFileQuery
	Offset uint32
	Limit  uint32
}

type ListFileResponse struct {
	List []*FileItem
}
