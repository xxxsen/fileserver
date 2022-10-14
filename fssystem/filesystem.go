package fssystem

import (
	"context"
	"fileserver/model"
	"fmt"
	"time"

	"github.com/xxxsen/common/errs"
	"google.golang.org/protobuf/proto"
)

type DBFileSystem struct {
	c  *config
	ps *PathResolver
}

func NewDBFileSystem(opts ...Option) (*DBFileSystem, error) {
	c := &config{}
	for _, opt := range opts {
		opt(c)
	}
	if c.fs == nil {
		return nil, errs.New(errs.ErrParam, "no core/fs found")
	}
	return &DBFileSystem{
		c:  c,
		ps: NewPathResolver(c.fs),
	}, nil
}

func (f *DBFileSystem) existNode(ctx context.Context, id uint64, name string) (bool, error) {
	rsp, err := f.c.fs.List(ctx, &model.ListFsItemRequest{
		ParentID: id,
		Query: &model.ListFsItemQuery{
			ChildFileName: proto.String(name),
		},
		NeedTotal: false,
		Offset:    0,
		Limit:     1,
	})
	if err != nil {
		return false, errs.Wrap(errs.ErrDatabase, "check file node fail", err)
	}
	if len(rsp.Items) == 0 {
		return false, nil
	}
	return true, nil
}

func (f *DBFileSystem) Mkdir(ctx context.Context, req *MkdirRequest) (*MkdirResponse, error) {
	if len(req.Name) == 0 {
		return nil, errs.New(errs.ErrParam, "nil name")
	}
	folderinfo, exist, err := f.ps.Resolve(ctx, req.Location)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "resolve location fail", err)
	}
	if !exist {
		return nil, errs.New(errs.ErrParam, "path not found, path:%s", req.Location)
	}
	if folderinfo.FileType != uint32(model.FsItemTypeFolder) {
		return nil, errs.New(errs.ErrParam, "parent should be a folder")
	}
	itemexist, err := f.existNode(ctx, folderinfo.ID, req.Name)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "check node fail", err)
	}
	if itemexist {
		return nil, errs.New(errs.ErrParam, "node exist")
	}
	now := uint64(time.Now().UnixMilli())
	rsp, err := f.c.fs.Create(ctx, &model.CreateFsItemRequest{
		Item: &model.FsItem{
			ParentID: folderinfo.ParentID,
			FileName: req.Name,
			FileType: uint32(model.FsItemTypeFolder),
			FileSize: 0,
			CTime:    now,
			MTime:    now,
			DownKey:  0,
		},
	})
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "create fail", err)
	}
	return &MkdirResponse{ID: rsp.ID}, nil
}

func (f *DBFileSystem) fsItem2fileItem(fsitem *model.FsItem) *FileItem {
	return &FileItem{
		ID:       fsitem.ID,
		ParentID: fsitem.ParentID,
		FileName: fsitem.FileName,
		FileSize: fsitem.FileSize,
		FileType: fsitem.FileType,
		CTime:    fsitem.CTime,
		MTime:    fsitem.MTime,
		DownKey:  fsitem.DownKey,
	}
}

func (f *DBFileSystem) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	info, exist, err := f.ps.Resolve(ctx, req.Location)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "resolve path fail", err)
	}
	if !exist {
		return nil, errs.New(errs.ErrParam, "not found")
	}
	innerRsp, err := f.c.fs.List(ctx, &model.ListFsItemRequest{
		ParentID:  info.ID,
		Query:     &model.ListFsItemQuery{},
		NeedTotal: true,
		Offset:    req.Offset,
		Limit:     req.Limit,
	})
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "list fail", err)
	}
	rsp := &ListResponse{
		Total: innerRsp.Total,
	}
	for _, item := range innerRsp.Items {
		rsp.List = append(rsp.List, f.fsItem2fileItem(item))
	}
	return rsp, nil
}

func (f *DBFileSystem) CreateFile(ctx context.Context, req *CreateFileRequest) (*CreateFileResponse, error) {
	folderinfo, exist, err := f.ps.Resolve(ctx, req.Location)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "resolve path fail", err)
	}
	if !exist {
		return nil, errs.New(errs.ErrParam, "location not found", err)
	}
	if folderinfo.FileType != uint32(model.FsItemTypeFolder) {
		return nil, errs.New(errs.ErrParam, "could not create file, parent is not folder")
	}
	itemexist, err := f.existNode(ctx, folderinfo.ID, req.Name)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "check node fail", err)
	}
	if itemexist {
		return nil, errs.New(errs.ErrParam, "node exist")
	}

	now := uint64(time.Now().UnixMilli())
	rsp, err := f.c.fs.Create(ctx, &model.CreateFsItemRequest{
		Item: &model.FsItem{
			ParentID: folderinfo.ID,
			FileName: req.Name,
			FileType: uint32(model.FsItemTypeFile),
			FileSize: req.Size,
			CTime:    now,
			MTime:    now,
			DownKey:  req.DownKey,
		},
	})
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "create file item fail", err)
	}
	return &CreateFileResponse{
		ID: rsp.ID,
	}, nil
}

func (f *DBFileSystem) OpenFile(ctx context.Context, req *OpenFileRequest) (*OpenFileResponse, error) {
	if len(req.Name) == 0 {
		return nil, errs.New(errs.ErrParam, "nil name")
	}
	location := fmt.Sprintf("%s/%s", req.Location, req.Name)
	item, exist, err := f.ps.Resolve(ctx, location)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "resolve path fail", err)
	}
	if !exist {
		return nil, errs.New(errs.ErrNotFound, "file not found")
	}
	return &OpenFileResponse{
		DownKey: item.DownKey,
		Size:    item.FileSize,
	}, nil
}

func (f *DBFileSystem) Info(ctx context.Context, req *InfoRequest) (*InfoResponse, error) {
	info, exist, err := f.ps.Resolve(ctx, req.Location)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "resolve path fail", err)
	}
	if !exist {
		return nil, errs.New(errs.ErrNotFound, "path not found")
	}
	innerRsp, err := f.c.fs.Info(ctx, &model.InfoFsItemRequest{
		ID: info.ID,
	})
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "info fail", err)
	}
	rsp := &InfoResponse{
		Exist: innerRsp.Exist,
	}
	if !innerRsp.Exist {
		return rsp, nil
	}
	rsp.Item = f.fsItem2fileItem(innerRsp.Item)
	return rsp, nil
}

func (f *DBFileSystem) Rename(ctx context.Context, req *RenameRequest) (*RenameResponse, error) {
	srcpts := f.ps.PathAsArray(req.SrcLocation)
	if len(srcpts) == 0 {
		return nil, errs.New(errs.ErrParam, "not allow to move root node")
	}
	dstpts := f.ps.PathAsArray(req.DstLocation)
	if f.ps.IsSubPath(dstpts, srcpts) {
		return nil, errs.New(errs.ErrParam, "not allow to move parent to child")
	}
	srcinfo, exist, err := f.ps.ResolveByArray(ctx, srcpts)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "resolve src fail", err)
	}
	if !exist {
		return nil, errs.New(errs.ErrNotFound, "src path not found", err)
	}
	dstinfo, exist, err := f.ps.ResolveByArray(ctx, dstpts)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "resolve dst fail", err)
	}
	if !exist {
		return nil, errs.New(errs.ErrNotFound, "dst path not found")
	}
	_, err = f.c.fs.Move(ctx, &model.MoveFsItemRequest{
		SrcID:      srcinfo.ID,
		ToParentID: dstinfo.ID,
	})
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "move fail", err)
	}
	return &RenameResponse{}, nil
}
