package dao

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fileserver/db"
	"fileserver/model"
	"hash/fnv"
	"time"

	"github.com/didi/gendry/builder"
	"github.com/xxxsen/common/errs"
)

var (
	mappingDBFields = []string{
		"id", "file_name", "hash_code", "check_sum", "create_time", "modify_time", "file_id",
	}
)

var MappingInfoDao MappingInfoService = &mappingInfoDaoImpl{}

type MappingInfoService interface {
	CreateMappingInfo(ctx context.Context, req *model.CreateMappingInfoRequest) (*model.CreateMappingInfoResponse, error)
	GetMappingInfo(ctx context.Context, req *model.GetMappingInfoRequest) (*model.GetMappingInfoResponse, error)
}

type mappingInfoDaoImpl struct {
}

func (d *mappingInfoDaoImpl) Table() string {
	return "mapping_info_tab"
}

func (d *mappingInfoDaoImpl) genHash(key string) (code uint32, ck string) {
	{
		h := fnv.New32a()
		h.Write([]byte(key))
		code = h.Sum32()
	}
	{
		h := sha1.New()
		h.Write([]byte(key))
		raw := h.Sum(nil)
		ck = base64.StdEncoding.EncodeToString(raw)
	}
	return
}

func (d *mappingInfoDaoImpl) CreateMappingInfo(ctx context.Context, req *model.CreateMappingInfoRequest) (*model.CreateMappingInfoResponse, error) {
	item := req.Item
	code, hash := d.genHash(item.FileName)
	now := time.Now().UnixMilli()
	data := []map[string]interface{}{
		{
			"file_name":   item.FileName,
			"create_time": now,
			"modify_time": now,
			"file_id":     item.FileId,
			"hash_code":   code,
			"check_sum":   hash,
		},
	}
	update := map[string]interface{}{
		"modify_time": now,
		"file_id":     item.FileId,
	}
	sql, args, err := builder.BuildInsertOnDuplicate(d.Table(), data, update)
	if err != nil {
		return nil, errs.Wrap(errs.ErrParam, "build insert sql fail", err)
	}
	client := db.GetFileDB()
	_, err = client.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "insert fail", err)
	}
	return &model.CreateMappingInfoResponse{}, nil
}

func (d *mappingInfoDaoImpl) GetMappingInfo(ctx context.Context, req *model.GetMappingInfoRequest) (*model.GetMappingInfoResponse, error) {
	code, ck := d.genHash(req.FileName)
	where := map[string]interface{}{
		"hash_code": code,
		"check_sum": ck,
	}
	fields := mappingDBFields
	sql, args, err := builder.BuildSelect(d.Table(), where, fields)
	if err != nil {
		return nil, errs.Wrap(errs.ErrParam, "build select fail", err)
	}
	client := db.GetFileDB()
	rows, err := client.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "query fail", err)
	}
	defer rows.Close()

	rsp := &model.GetMappingInfoResponse{}
	for rows.Next() {
		item := &model.MappingInfoItem{}
		if err := rows.Scan(&item.Id, &item.FileName, &item.HashCode, &item.CheckSum,
			&item.CreateTime, &item.ModifyTime, &item.FileId); err != nil {

			return nil, errs.Wrap(errs.ErrDatabase, "scan fail", err)
		}
		rsp.Item = item
		break
	}
	if err := rows.Err(); err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "scan fail", err)
	}

	return rsp, nil
}
