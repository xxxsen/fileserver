package dao

import (
	"context"
	"fmt"
	"tgfile/db"
	"tgfile/entity"
	"time"

	"github.com/xxxsen/common/database/kv"
)

type IFilePartDao interface {
	CreateFilePart(ctx context.Context, req *entity.CreateFilePartRequest) (*entity.CreateFilePartResponse, error)
	GetFilePartInfo(ctx context.Context, req *entity.GetFilePartInfoRequest) (*entity.GetFilePartInfoResponse, error)
}

type filePartDaoImpl struct {
}

func NewFilePartDao() IFilePartDao {
	return &filePartDaoImpl{}
}

func (f *filePartDaoImpl) table() string {
	return "tg_file_part_tab"
}

func (f *filePartDaoImpl) buildKey(fileid uint64, idx int32) string {
	return fmt.Sprintf("tgfile:filepart:%d:%d", fileid, idx)
}

func (f *filePartDaoImpl) CreateFilePart(ctx context.Context, req *entity.CreateFilePartRequest) (*entity.CreateFilePartResponse, error) {
	now := time.Now().UnixMilli()
	item := &entity.FilePartInfoItem{
		FileId:     req.FileId,
		FilePartId: req.FilePartId,
		FileKey:    req.FileKey,
		Ctime:      now,
		Mtime:      now,
	}
	if err := kv.SetJsonObject(ctx, db.GetClient(), f.table(), f.buildKey(req.FileId, req.FilePartId), item); err != nil {
		return nil, err
	}
	return &entity.CreateFilePartResponse{}, nil
}

func (f *filePartDaoImpl) GetFilePartInfo(ctx context.Context, req *entity.GetFilePartInfoRequest) (*entity.GetFilePartInfoResponse, error) {
	keys := make([]string, 0, len(req.FilePartId))
	mapping := make(map[int32]string, len(req.FilePartId))
	for _, idx := range req.FilePartId {
		key := f.buildKey(req.FileId, idx)
		keys = append(keys, key)
		mapping[idx] = key
	}
	rs, err := kv.MultiGetJsonObject[entity.FilePartInfoItem](ctx, db.GetClient(), f.table(), keys)
	if err != nil {
		return nil, err
	}
	rsp := &entity.GetFilePartInfoResponse{}
	for _, idx := range req.FilePartId {
		val, ok := rs[mapping[idx]]
		if !ok {
			continue
		}
		rsp.List = append(rsp.List, val)
	}
	return rsp, nil
}
