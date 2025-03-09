package dao

import (
	"context"
	"tgfile/db"
	"tgfile/entity"
	"time"

	"github.com/xxxsen/common/database/kv"
)

const (
	defaultFileMappingPrefix = "tgfile:mapping:"
)

type IterFileMappingFunc func(ctx context.Context, name string, fileid uint64) (bool, error)

type IFileMappingDao interface {
	GetFileMapping(ctx context.Context, req *entity.GetFileMappingRequest) (*entity.GetFileMappingResponse, bool, error)
	CreateFileMapping(ctx context.Context, req *entity.CreateFileMappingRequest) (*entity.CreateFileMappingResponse, error)
	IterFileMapping(ctx context.Context, cb IterFileMappingFunc) error
}

type fileMappingDao struct {
}

func NewFileMappingDao() IFileMappingDao {
	return &fileMappingDao{}
}

func (f *fileMappingDao) table() string {
	return "tg_file_mapping_tab"
}

func (f *fileMappingDao) buildKey(name string) string {
	return defaultFileMappingPrefix + name
}

func (f *fileMappingDao) GetFileMapping(ctx context.Context, req *entity.GetFileMappingRequest) (*entity.GetFileMappingResponse, bool, error) {
	item, ok, err := kv.GetJsonObject[entity.FileMappingItem](ctx, db.GetClient(), f.table(), f.buildKey(req.FileName))
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	return &entity.GetFileMappingResponse{Item: item}, true, nil
}

func (f *fileMappingDao) CreateFileMapping(ctx context.Context, req *entity.CreateFileMappingRequest) (*entity.CreateFileMappingResponse, error) {
	now := time.Now().UnixMilli()
	item := &entity.FileMappingItem{
		FileName: req.FileName,
		FileId:   req.FileId,
		Ctime:    uint64(now),
		Mtime:    uint64(now),
	}
	if err := kv.SetJsonObject(ctx, db.GetClient(), f.table(), f.buildKey(req.FileName), item); err != nil {
		return nil, err
	}
	return &entity.CreateFileMappingResponse{}, nil
}

func (f *fileMappingDao) IterFileMapping(ctx context.Context, cb IterFileMappingFunc) error {
	return kv.IterJsonObject(ctx, db.GetClient(), f.table(), defaultFileMappingPrefix, func(ctx context.Context, key string, val *entity.FileMappingItem) (bool, error) {
		return cb(ctx, val.FileName, val.FileId)
	})
}
