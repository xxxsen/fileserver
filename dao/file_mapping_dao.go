package dao

import (
	"context"
	"fileserver/db"
	"fileserver/entity"
	"time"

	"github.com/didi/gendry/builder"
	"github.com/xxxsen/common/database/dbkit"
)

var FileMappingDao IFileMappingDao = NewFileMappingDao()

type IFileMappingDao interface {
	GetFileMapping(ctx context.Context, req *entity.GetFileMappingRequest) (*entity.GetFileMappingResponse, bool, error)
	CreateFileMapping(ctx context.Context, req *entity.CreateFileMappingRequest) (*entity.CreateFileMappingResponse, error)
}

type fileMappingDao struct {
}

func NewFileMappingDao() IFileMappingDao {
	return &fileMappingDao{}
}

func (f *fileMappingDao) table() string {
	return "tg_file_mapping_tab"
}

func (f *fileMappingDao) GetFileMapping(ctx context.Context, req *entity.GetFileMappingRequest) (*entity.GetFileMappingResponse, bool, error) {
	where := map[string]interface{}{
		"file_name": req.FileName,
	}
	rs := make([]*entity.GetFileMappingItem, 0, 1)
	dbkit.SimpleQuery(ctx, db.GetClient(), f.table(), where, &rs)
	if len(rs) == 0 {
		return nil, false, nil
	}
	return &entity.GetFileMappingResponse{
		Item: rs[0],
	}, true, nil
}

func (f *fileMappingDao) CreateFileMapping(ctx context.Context, req *entity.CreateFileMappingRequest) (*entity.CreateFileMappingResponse, error) {
	now := time.Now().UnixMilli()
	data := []map[string]interface{}{
		{
			"file_name": req.FileName,
			"file_id":   req.FileId,
			"ctime":     now,
			"mtime":     now,
		},
	}
	sql, args, err := builder.BuildInsert(f.table(), data)
	if err != nil {
		return nil, err
	}
	if _, err := db.GetClient().ExecContext(ctx, sql, args...); err != nil {
		return nil, err
	}
	return &entity.CreateFileMappingResponse{}, nil
}
