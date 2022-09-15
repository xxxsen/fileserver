package dao

import (
	"context"
	"database/sql"
	"fileserver/db"
	"fileserver/model"

	"github.com/xxxsen/common/errs"

	"github.com/didi/gendry/builder"
)

var (
	fileDBFields = []string{
		"id", "file_name", "hash", "file_size", "create_time", "down_key", "extra", "file_key",
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
	return db.GetMediaDB()
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
	sql, args, err := builder.BuildSelect(d.Table(), where, d.Fields())
	if err != nil {
		return nil, errs.Wrap(errs.ErrParam, "build select", err)
	}
	rows, err := d.Client().QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "select fail", err)
	}
	defer rows.Close()
	rs := make([]*model.FileItem, 0, req.Limit)
	for rows.Next() {
		item := &model.FileItem{}
		if err := rows.Scan(&item.Id, &item.FileName,
			&item.Hash, &item.FileSize, &item.CreateTime,
			&item.DownKey, &item.Extra, &item.FileKey); err != nil {
			return nil, errs.Wrap(errs.ErrDatabase, "scan fail", err)
		}
		rs = append(rs, item)
	}
	if err := rows.Err(); err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "scan fail", err)
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
		},
	}
	sql, args, err := builder.BuildInsertIgnore(d.Table(), data)
	if err != nil {
		return nil, errs.Wrap(errs.ErrParam, "build insert", err)
	}
	_, err = d.Client().ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "insert fail", err)
	}
	return &model.CreateFileResponse{}, nil
}
