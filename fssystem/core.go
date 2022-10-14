package fssystem

import (
	"context"
)

type CreateFileRequest struct {
	Location string
	Name     string
	Size     uint64
	DownKey  uint64
}

type CreateFileResponse struct {
	ID uint64
}

type ListRequest struct {
	Location string
	Offset   uint64
	Limit    uint32
}

type FileItem struct {
	ID       uint64
	ParentID uint64
	FileName string
	FileSize uint64
	FileType uint32
	CTime    uint64
	MTime    uint64
	DownKey  uint64
}

type ListResponse struct {
	List  []*FileItem
	Total uint64
}

type OpenFileRequest struct {
	Location string
	Name     string
}

type OpenFileResponse struct {
	DownKey uint64
	Size    uint64
}

type RenameRequest struct {
	SrcLocation string
	DstLocation string
}

type RenameResponse struct {
}

type InfoRequest struct {
	Location string
}

type InfoResponse struct {
	Item  *FileItem
	Exist bool
}

type MkdirRequest struct {
	Location string
	Name     string
}

type MkdirResponse struct {
	ID uint64
}

type IFsSystemCore interface {
	Mkdir(ctx context.Context, req *MkdirRequest) (*MkdirResponse, error)
	List(ctx context.Context, req *ListRequest) (*ListResponse, error)
	CreateFile(ctx context.Context, req *CreateFileRequest) (*CreateFileResponse, error)
	OpenFile(ctx context.Context, req *OpenFileRequest) (*OpenFileResponse, error)
	Info(ctx context.Context, req *InfoRequest) (*InfoResponse, error)
	Rename(ctx context.Context, req *RenameRequest) (*RenameResponse, error)
}
