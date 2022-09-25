package model

type GetMappingInfoRequest struct {
	FileName string
}

type GetMappingInfoResponse struct {
	Item *MappingInfoItem
}

type CreateMappingInfoRequest struct {
	Item *MappingInfoItem
}

type MappingInfoItem struct {
	Id         uint64
	FileName   string
	FileId     uint64
	CreateTime uint64
	ModifyTime uint64
	HashCode   uint32
	CheckSum   string
}

type CreateMappingInfoResponse struct {
}
