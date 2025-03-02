package core

import (
	"context"
	"fmt"
)

type MultiCore struct {
	IFsCore
	downloaders map[int]IFsCore
}

func NewMultiCore(basic IFsCore, downloaders ...IFsCore) (*MultiCore, error) {
	if len(downloaders) == 0 {
		downloaders = []IFsCore{basic}
	}
	m := make(map[int]IFsCore)
	for _, c := range downloaders {
		if _, ok := m[int(c.StType())]; ok {
			return nil, fmt.Errorf("multi core with same st type found, st:%d", c.StType())
		}
		m[int(c.StType())] = c
	}
	return &MultiCore{
		IFsCore:     basic,
		downloaders: m,
	}, nil
}

func (c *MultiCore) FileDownload(ctx context.Context, fctx *FileDownloadRequest) (*FileDownloadResponse, error) {
	fs, ok := c.downloaders[int(fctx.StType)]
	if !ok {
		return nil, fmt.Errorf("core not found, type:%d", fctx.StType)
	}
	return fs.FileDownload(ctx, fctx)
}
