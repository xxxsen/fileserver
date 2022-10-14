package model

type FsItemType int

const (
	FsItemTypeFolder FsItemType = 0
	FsItemTypeFile   FsItemType = 1
)

type ListFsItemQuery struct {
	ChildFileName *string
}

type ListFsItemRequest struct {
	ParentID  uint64
	Query     *ListFsItemQuery
	NeedTotal bool
	Offset    uint64
	Limit     uint32
}

type FsItem struct {
	ID       uint64
	ParentID uint64
	NameCode uint32
	FileName string
	FileType uint32 //refer: FsItemType
	FileSize uint64
	CTime    uint64
	MTime    uint64
	DownKey  uint64
}

type ListFsItemResponse struct {
	Items []*FsItem
	Total uint64
}

type RemoveFsItemRequest struct {
	IDs []uint64
}

type RemoveFsItemResponse struct {
}

type CreateFsItemRequest struct {
	Item *FsItem
}

type CreateFsItemResponse struct {
	ID uint64
}

type ModifyFsItem struct {
	FileName *string
	DownKey  *string
}

type ModifyFsItemRequest struct {
	ID   uint64
	Item *ModifyFsItem
}

type ModifyFsItemResponse struct {
}

type InfoFsItemRequest struct {
	ID uint64
}

type InfoFsItemResponse struct {
	Item  *FsItem
	Exist bool
}

type MoveFsItemRequest struct {
	SrcID      uint64
	ToParentID uint64
}

type MoveFsItemResponse struct {
}
