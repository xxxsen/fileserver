package entity

type GetFilePartInfoRequest struct {
	FileId     uint64
	FilePartId []int32
}

type GetFilePartInfoItem struct {
	Id         uint64 `json:"id"`
	FileId     uint64 `json:"file_id"`
	FilePartId int32  `json:"file_part_id"`
	FileKey    string `json:"file_key"`
	Ctime      int64  `json:"ctime"`
	Mtime      int64  `json:"mtime"`
}

type GetFilePartInfoResponse struct {
	List []*GetFilePartInfoItem
}

type GetFilePartCountRequest struct {
	FileId uint64
}

type GetFilePartCountResponse struct {
	Count int32
}
