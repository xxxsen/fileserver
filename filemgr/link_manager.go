package filemgr

import (
	"context"
	"fileserver/service"
	"fmt"
	"io"
)

type ILinkManager interface {
	CreateLink(ctx context.Context, link string, fileid uint64) error
	OpenLink(ctx context.Context, link string, pos int64) (io.ReadSeekCloser, error)
}

type defaultLinkManager struct {
	fmgr IFileManager
}

func (d *defaultLinkManager) CreateLink(ctx context.Context, link string, fileid uint64) error {
	if err := service.FileMappingService.CreateFileMapping(ctx, link, fileid); err != nil {
		return err
	}
	return nil
}

func (d *defaultLinkManager) OpenLink(ctx context.Context, link string, pos int64) (io.ReadSeekCloser, error) {
	fid, ok, err := service.FileMappingService.GetFileMapping(ctx, link)
	if err != nil {
		return nil, fmt.Errorf("open mapping failed, err:%w", err)
	}
	if !ok {
		return nil, fmt.Errorf("link not found")
	}
	return d.fmgr.Open(ctx, fid, pos)
}

func NewLinkManager(fmgr IFileManager) ILinkManager {
	return &defaultLinkManager{
		fmgr: fmgr,
	}
}
