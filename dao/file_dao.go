package dao

import (
	"context"
	"fileserver/constant"
	"fileserver/db"
	"fileserver/entity"
	"fmt"
	"time"

	"github.com/didi/gendry/builder"
	"github.com/xxxsen/common/database/dbkit"
	"github.com/xxxsen/common/idgen"
)

var FileDao IFileDao = NewFileDao()

type IFileDao interface {
	CreateFileDraft(ctx context.Context, req *entity.CreateFileDraftRequest) (*entity.CreateFileDraftResponse, error)
	MarkFileReady(ctx context.Context, req *entity.MarkFileReadyRequest) (*entity.MarkFileReadyResponse, error)
	GetFileInfo(ctx context.Context, req *entity.GetFileInfoRequest) (*entity.GetFileInfoResponse, error)
}

type fileDaoImpl struct{}

func NewFileDao() IFileDao {
	return &fileDaoImpl{}
}

func (f *fileDaoImpl) table() string {
	return "tg_file_tab"
}

func (f *fileDaoImpl) CreateFileDraft(ctx context.Context, req *entity.CreateFileDraftRequest) (*entity.CreateFileDraftResponse, error) {
	fileid := idgen.NextId()
	now := time.Now().UnixMilli()
	data := []map[string]interface{}{
		{
			"file_name":       req.FileName,
			"file_size":       req.FileSize,
			"file_part_count": req.FilePartCount,
			"file_id":         fileid,
			"ctime":           now,
			"mtime":           now,
			"file_state":      constant.FileStateInit,
		},
	}
	sql, args, err := builder.BuildInsert(f.table(), data)
	if err != nil {
		return nil, fmt.Errorf("build insert failed, err:%w", err)
	}
	if _, err := db.GetClient().ExecContext(ctx, sql, args...); err != nil {
		return nil, fmt.Errorf("exec insert failed, err:%w", err)
	}
	return &entity.CreateFileDraftResponse{
		FileId: fileid,
	}, nil
}

func (f *fileDaoImpl) MarkFileReady(ctx context.Context, req *entity.MarkFileReadyRequest) (*entity.MarkFileReadyResponse, error) {
	where := map[string]interface{}{
		"file_id":    req.FileID,
		"file_state": constant.FileStateReady,
	}
	update := map[string]interface{}{
		"mtime": time.Now().UnixMilli(),
	}
	sql, args, err := builder.BuildUpdate(f.table(), where, update)
	if err != nil {
		return nil, err
	}
	if _, err := db.GetClient().ExecContext(ctx, sql, args...); err != nil {
		return nil, err
	}
	return &entity.MarkFileReadyResponse{}, nil
}

func (f *fileDaoImpl) GetFileInfo(ctx context.Context, req *entity.GetFileInfoRequest) (*entity.GetFileInfoResponse, error) {
	where := map[string]interface{}{
		"file_id in": req.FileIds,
	}
	rs := make([]*entity.GetFileInfoItem, 0, len(req.FileIds))
	if err := dbkit.SimpleQuery(ctx, db.GetClient(), f.table(), where, &rs, dbkit.ScanWithTagName("json")); err != nil {
		return nil, err
	}
	return &entity.GetFileInfoResponse{List: rs}, nil
}
