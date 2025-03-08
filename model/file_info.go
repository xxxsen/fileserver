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
	Id         uint64 `json:"id"`
	FileName   string `json:"file_name"`
	Hash       string `json:"hash"`
	FileSize   uint64 `json:"file_size"`
	CreateTime uint64 `json:"create_time"`
	DownKey    uint64 `json:"down_key"`
	Extra      []byte `json:"extra"`
	FileKey    string `json:"file_key"`
	StType     uint8  `json:"st_type"`
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
