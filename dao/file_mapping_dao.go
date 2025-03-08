package dao

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fileserver/db"
	"fileserver/entity"
	"fmt"
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

func (f *fileMappingDao) name2hash(name string) string {
	h := md5.New()
	h.Write([]byte(name))
	return hex.EncodeToString(h.Sum(nil))
}

func (f *fileMappingDao) GetFileMapping(ctx context.Context, req *entity.GetFileMappingRequest) (*entity.GetFileMappingResponse, bool, error) {
	hash := f.name2hash(req.FileName)
	where := map[string]interface{}{
		"file_hash": hash,
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
	fhash := f.name2hash(req.FileName)
	data := []map[string]interface{}{
		{
			"file_name": req.FileName,
			"file_hash": fhash,
			"file_id":   req.FileId,
			"ctime":     now,
			"mtime":     now,
		},
	}
	sql, args, err := builder.BuildInsert(f.table(), data)
	if err != nil {
		return nil, err
	}
	_, insertErr := db.GetClient().ExecContext(ctx, sql, args...)
	if insertErr == nil {
		return &entity.CreateFileMappingResponse{}, nil
	}
	//尝试update
	where := map[string]interface{}{
		"file_hash": fhash,
	}
	update := map[string]interface{}{
		"file_id": req.FileId,
		"mtime":   now,
	}
	sql, args, err = builder.BuildUpdate(f.table(), where, update)
	if err != nil {
		return nil, err
	}
	rs, err := db.GetClient().ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	affect, err := rs.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affect == 0 {
		return nil, fmt.Errorf("insert err and try update failed, insert err:%w", insertErr)
	}
	return &entity.CreateFileMappingResponse{}, nil
}
