package entity

type GetFileMappingRequest struct {
	FileName string
}

type GetFileMappingItem struct {
	Id       uint64 `json:"id"`
	FileName string `json:"file_name"`
	FileHash string `json:"file_hash"`
	FileId   uint64 `json:"file_id"`
	Ctime    uint64 `json:"ctime"`
	Mtime    uint64 `json:"mtime"`
}

type GetFileMappingResponse struct {
	Item *GetFileMappingItem
}

type CreateFileMappingRequest struct {
	FileName string
	FileId   uint64
}

type CreateFileMappingResponse struct {
}
