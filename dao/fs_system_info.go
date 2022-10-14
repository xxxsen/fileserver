package dao

import (
	"context"
	"database/sql"
	"errors"
	"fileserver/db"
	"fileserver/model"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/didi/gendry/builder"
	"github.com/xxxsen/common/errs"
)

const (
	idRoot = 0
)

var fsSystemFields = []string{
	"id", "parent_id", "name_code", "file_name", "file_type", "file_size", "ctime", "mtime", "down_key",
}

var FsSystemDao FsSystemService = &fsSystemImpl{}

type FsSystemService interface {
	Root() uint64
	List(ctx context.Context, req *model.ListFsItemRequest) (*model.ListFsItemResponse, error)
	Remove(ctx context.Context, req *model.RemoveFsItemRequest) (*model.RemoveFsItemResponse, error)
	TxCreate(ctx context.Context, tx db.IExecutor, req *model.CreateFsItemRequest) (*model.CreateFsItemResponse, error)
	Create(ctx context.Context, req *model.CreateFsItemRequest) (*model.CreateFsItemResponse, error)
	Modify(ctx context.Context, req *model.ModifyFsItemRequest) (*model.ModifyFsItemResponse, error)
	TxInfo(ctx context.Context, tx db.IQueryer, req *model.InfoFsItemRequest) (*model.InfoFsItemResponse, error)
	Info(ctx context.Context, req *model.InfoFsItemRequest) (*model.InfoFsItemResponse, error)
	TxMove(ctx context.Context, tx db.IExecutor, req *model.MoveFsItemRequest) (*model.MoveFsItemResponse, error)
	Move(ctx context.Context, req *model.MoveFsItemRequest) (*model.MoveFsItemResponse, error)
}

type fsSystemImpl struct {
}

func (f *fsSystemImpl) Table() string {
	return "fs_system_tab"
}

func (f *fsSystemImpl) Client() *sql.DB {
	return db.GetFileDB()
}

func (f *fsSystemImpl) Root() uint64 {
	return idRoot
}

func (f *fsSystemImpl) nameCode(name string) uint32 {
	h := fnv.New32()
	_, _ = h.Write([]byte(fmt.Sprintf("s:%s:%d:e", name, len(name))))
	return h.Sum32()
}

func (f *fsSystemImpl) queryTotal(ctx context.Context, where map[string]interface{}) (uint64, error) {
	delete(where, "_limit")
	delete(where, "_orderby")
	query, args, err := builder.BuildSelect(f.Table(), where, []string{"count(*)"})
	if err != nil {
		return 0, errs.Wrap(errs.ErrParam, "build select total fail", err)
	}
	row := f.Client().QueryRowContext(ctx, query, args...)
	var total sql.NullInt64
	err = row.Scan(&total)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, errs.Wrap(errs.ErrDatabase, "read total fail", err)
	}
	if !total.Valid {
		return 0, nil
	}
	return uint64(total.Int64), nil
}

func (f *fsSystemImpl) List(ctx context.Context, req *model.ListFsItemRequest) (*model.ListFsItemResponse, error) {
	where := map[string]interface{}{
		"parent_id": req.ParentID,
		"_limit":    []uint{uint(req.Offset), uint(req.Limit)},
	}
	if req.Query != nil {
		if req.Query.ChildFileName != nil && len(*req.Query.ChildFileName) > 0 {
			where["file_name"] = *req.Query.ChildFileName
			where["name_code"] = f.nameCode(*req.Query.ChildFileName)
		}
	}
	fields := fsSystemFields
	sql, args, err := builder.BuildSelect(f.Table(), where, fields)
	if err != nil {
		return nil, errs.Wrap(errs.ErrParam, "build select fail", err)
	}
	rows, err := f.Client().QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "query fail", err)
	}
	defer rows.Close()
	rsp := &model.ListFsItemResponse{}
	for rows.Next() {
		item := &model.FsItem{}
		if err := rows.Scan(&item.ID, &item.ParentID, &item.NameCode,
			&item.FileName, &item.FileType, &item.FileSize, &item.CTime,
			&item.MTime, &item.DownKey); err != nil {

			return nil, errs.Wrap(errs.ErrDatabase, "scan item fail", err)
		}
		rsp.Items = append(rsp.Items, item)
	}
	if !req.NeedTotal {
		return rsp, nil
	}
	total, err := f.queryTotal(ctx, where)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "query total fail", err)
	}
	rsp.Total = total
	return rsp, nil
}

func (f *fsSystemImpl) Remove(ctx context.Context, req *model.RemoveFsItemRequest) (*model.RemoveFsItemResponse, error) {
	where := map[string]interface{}{
		"id in": req.IDs,
	}
	sql, args, err := builder.BuildDelete(f.Table(), where)
	if err != nil {
		return nil, errs.Wrap(errs.ErrParam, "build delete fail", err)
	}
	if _, err := f.Client().ExecContext(ctx, sql, args...); err != nil {
		return nil, errs.Wrap(errs.ErrParam, "delete rec fail", err)
	}
	return &model.RemoveFsItemResponse{}, nil
}

func (f *fsSystemImpl) TxMove(ctx context.Context, tx db.IExecutor, req *model.MoveFsItemRequest) (*model.MoveFsItemResponse, error) {
	where := map[string]interface{}{
		"id": req.SrcID,
	}
	update := map[string]interface{}{
		"parent_id": req.ToParentID,
	}
	sql, args, err := builder.BuildUpdate(f.Table(), where, update)
	if err != nil {
		return nil, errs.Wrap(errs.ErrParam, "build update fail", err)
	}
	if _, err := tx.ExecContext(ctx, sql, args...); err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "exec move fail", err)
	}
	return &model.MoveFsItemResponse{}, nil
}

func (f *fsSystemImpl) Move(ctx context.Context, req *model.MoveFsItemRequest) (*model.MoveFsItemResponse, error) {
	tx, err := f.Client().Begin()
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "begin tx fail", err)
	}
	defer tx.Rollback()

	if err := f.isExistFolder(ctx, tx, req.ToParentID); err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "check id fail", err)
	}
	rsp, err := f.TxMove(ctx, tx, req)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "do move fail", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "commit move fail", err)
	}
	return rsp, nil
}

func (f *fsSystemImpl) TxCreate(ctx context.Context, tx db.IExecutor, req *model.CreateFsItemRequest) (*model.CreateFsItemResponse, error) {
	data := []map[string]interface{}{
		{
			"parent_id": req.Item.ParentID,
			"file_name": req.Item.FileName,
			"name_code": f.nameCode(req.Item.FileName),
			"file_size": req.Item.FileSize,
			"file_type": req.Item.FileType,
			"ctime":     req.Item.CTime,
			"mtime":     req.Item.MTime,
			"down_key":  req.Item.DownKey,
		},
	}
	sql, args, err := builder.BuildInsert(f.Table(), data)
	if err != nil {
		return nil, errs.Wrap(errs.ErrParam, "insert record fail", err)
	}
	rs, err := tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "exec insert fail", err)
	}
	lstid, err := rs.LastInsertId()
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "get last id fail", err)
	}
	return &model.CreateFsItemResponse{
		ID: uint64(lstid),
	}, nil
}

func (f *fsSystemImpl) Create(ctx context.Context, req *model.CreateFsItemRequest) (*model.CreateFsItemResponse, error) {
	tx, err := f.Client().Begin()
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "begin tx fail", err)
	}
	defer tx.Rollback()
	//check parent exist and parent should be a folder
	//shall we do this logic in dao?
	if err := f.isExistFolder(ctx, tx, req.Item.ParentID); err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "check id fail", err)
	}
	rsp, err := f.TxCreate(ctx, tx, req)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "do create fail", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "commit create fail", err)
	}
	return rsp, nil
}

func (f *fsSystemImpl) Modify(ctx context.Context, req *model.ModifyFsItemRequest) (*model.ModifyFsItemResponse, error) {
	where := map[string]interface{}{
		"id": req.ID,
	}
	update := map[string]interface{}{
		"mtime": time.Now().UnixMilli(),
	}
	if req.Item != nil {
		if req.Item.FileName != nil && len(*req.Item.FileName) > 0 {
			update["file_name"] = *req.Item.FileName
			update["name_code"] = f.nameCode(*req.Item.FileName)
		}
		if req.Item.DownKey != nil && len(*req.Item.DownKey) > 0 {
			update["down_key"] = *req.Item.DownKey
		}
	}
	sql, args, err := builder.BuildUpdate(f.Table(), where, update)
	if err != nil {
		return nil, errs.Wrap(errs.ErrParam, "build sql fail", err)
	}
	if _, err := f.Client().ExecContext(ctx, sql, args...); err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "exec update fail", err)
	}
	return &model.ModifyFsItemResponse{}, nil
}

func (f *fsSystemImpl) Info(ctx context.Context, req *model.InfoFsItemRequest) (*model.InfoFsItemResponse, error) {
	return f.TxInfo(ctx, f.Client(), req)
}

func (f *fsSystemImpl) isExistFolder(ctx context.Context, tx db.IQueryer, id uint64) error {
	infoRsp, err := f.TxInfo(ctx, tx, &model.InfoFsItemRequest{
		ID: id,
	})
	if err != nil {
		return errs.Wrap(errs.ErrDatabase, "read id info fail", err)
	}
	if !infoRsp.Exist {
		return errs.New(errs.ErrParam, "id not found")
	}
	if infoRsp.Item.FileType != uint32(model.FsItemTypeFolder) {
		return errs.New(errs.ErrParam, "id not folder")
	}
	return nil
}

func (f *fsSystemImpl) TxInfo(ctx context.Context, tx db.IQueryer, req *model.InfoFsItemRequest) (*model.InfoFsItemResponse, error) {
	where := map[string]interface{}{
		"id": req.ID,
	}
	fields := fsSystemFields
	qSql, args, err := builder.BuildSelect(f.Table(), where, fields)
	if err != nil {
		return nil, errs.Wrap(errs.ErrParam, "build select fail", err)
	}
	row := tx.QueryRowContext(ctx, qSql, args...)
	item := &model.FsItem{}
	err = row.Scan(&item.ID, &item.ParentID, &item.NameCode, &item.FileName, &item.FileType, &item.FileSize, &item.CTime, &item.MTime, &item.DownKey)
	if err == sql.ErrNoRows {
		return &model.InfoFsItemResponse{Exist: false}, nil
	}
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "scan fail", err)
	}
	return &model.InfoFsItemResponse{Item: item, Exist: true}, nil
}
