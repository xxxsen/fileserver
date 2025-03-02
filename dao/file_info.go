package dao

import (
	"context"
	"database/sql"
	"fileserver/db"
	"fileserver/model"
	"fmt"

	"github.com/didi/gendry/builder"
	"github.com/xxxsen/common/database/dbkit"
)

var (
	fileDBFields = []string{
		"id", "file_name", "hash", "file_size", "create_time", "down_key", "extra", "file_key", "st_type",
	}
)

var FileInfoDao FileInfoService = &fileInfoDaoImpl{}

type FileInfoService interface {
	GetFile(ctx context.Context, req *model.GetFileRequest) (*model.GetFileResponse, bool, error)
	ListFile(ctx context.Context, req *model.ListFileRequest) (*model.ListFileResponse, error)
	CreateFile(ctx context.Context, req *model.CreateFileRequest) (*model.CreateFileResponse, error)
}

type fileInfoDaoImpl struct {
}

func (d *fileInfoDaoImpl) Table() string {
	return "file_info_tab"
}

func (d *fileInfoDaoImpl) Client() *sql.DB {
	return db.GetClient()
}

func (d *fileInfoDaoImpl) Fields() []string {
	return fileDBFields
}

func (d *fileInfoDaoImpl) ListFile(ctx context.Context, req *model.ListFileRequest) (*model.ListFileResponse, error) {
	where := map[string]interface{}{
		"_limit": []uint{uint(req.Offset), uint(req.Limit)},
	}
	if len(req.Query.DownKey) > 0 {
		where["down_key in"] = req.Query.DownKey
	}
	if len(req.Query.ID) > 0 {
		where["id in"] = req.Query.ID
	}
	rs := make([]*model.FileItem, 0, req.Limit)
	if err := dbkit.SimpleQuery(ctx, d.Client(), d.Table(), where, &rs, dbkit.ScanWithTagName("json")); err != nil {
		return nil, err
	}
	return &model.ListFileResponse{List: rs}, nil
}

func (d *fileInfoDaoImpl) GetFile(ctx context.Context, req *model.GetFileRequest) (*model.GetFileResponse, bool, error) {
	listReq := &model.ListFileRequest{
		Query: &model.ListFileQuery{
			DownKey: []uint64{req.DownKey},
		},
		Offset: 0,
		Limit:  1,
	}
	listRsp, err := d.ListFile(ctx, listReq)
	if err != nil {
		return nil, false, err
	}
	if len(listRsp.List) != 1 {
		return nil, false, nil
	}
	return &model.GetFileResponse{Item: listRsp.List[0]}, true, nil
}

func (d *fileInfoDaoImpl) CreateFile(ctx context.Context, req *model.CreateFileRequest) (*model.CreateFileResponse, error) {
	data := []map[string]interface{}{
		{
			"file_name":   req.Item.FileName,
			"hash":        req.Item.Hash,
			"file_size":   req.Item.FileSize,
			"create_time": req.Item.CreateTime,
			"down_key":    req.Item.DownKey,
			"extra":       req.Item.Extra,
			"file_key":    req.Item.FileKey,
			"st_type":     req.Item.StType,
		},
	}
	sql, args, err := builder.BuildInsert(d.Table(), data)
	if err != nil {
		return nil, fmt.Errorf("build insert, err:%w", err)
	}
	_, err = d.Client().ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("insert fail, err:%w", err)
	}
	return &model.CreateFileResponse{}, nil
}
