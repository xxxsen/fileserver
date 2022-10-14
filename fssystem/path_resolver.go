package fssystem

import (
	"context"
	"fileserver/dao"
	"fileserver/model"
	"strings"

	"github.com/xxxsen/common/errs"
	"google.golang.org/protobuf/proto"
)

type PathResolver struct {
	fs dao.FsSystemService
}

func NewPathResolver(fs dao.FsSystemService) *PathResolver {
	return &PathResolver{
		fs: fs,
	}
}

func (p *PathResolver) IsSubPath(origin []string, detect []string) bool {
	if len(detect) == 0 {
		return true
	}
	for i := 0; i < len(detect); i++ {
		if i >= len(origin) {
			return false
		}
		if detect[i] != origin[i] {
			return false
		}
	}
	return true
}

func (p *PathResolver) PathAsArray(location string) []string {
	arrs := strings.Split(location, "/")
	if len(arrs) == 0 {
		return []string{}
	}
	index := 0
	for i := 0; i < len(arrs); i++ {
		if arrs[i] == "." || arrs[i] == "" {
			continue
		}
		if arrs[i] == ".." && index > 0 {
			index--
		}
		arrs[index] = arrs[i]
		index++
	}
	return arrs[:index]
}

func (p *PathResolver) ResolveByArray(ctx context.Context, paths []string) (*model.FsItem, bool, error) {
	root := p.fs.Root()
	if len(paths) == 0 {
		return &model.FsItem{
			ID:       root,
			ParentID: root,
			FileType: uint32(model.FsItemTypeFolder),
		}, true, nil
	}
	index := 0
	var pid = root
	var info *model.FsItem
	for index < len(paths) {
		name := paths[index]

		rsp, err := p.fs.List(ctx, &model.ListFsItemRequest{
			ParentID: pid,
			Query: &model.ListFsItemQuery{
				ChildFileName: proto.String(name),
			},
			NeedTotal: false,
			Offset:    0,
			Limit:     1,
		})
		if err != nil {
			return nil, false, errs.Wrap(errs.ErrDatabase, "read info by search id/name fail", err)
		}
		if len(rsp.Items) == 0 {
			return nil, false, nil
		}
		item := rsp.Items[0]
		pid = item.ID
		info = item
		index++
	}
	return info, true, nil
}

func (p *PathResolver) Resolve(ctx context.Context, location string) (*model.FsItem, bool, error) {
	paths := p.PathAsArray(location)
	return p.ResolveByArray(ctx, paths)
}
