package cache

import (
	"context"
	"fmt"
	"tgfile/cache"
	"tgfile/dao"
	"tgfile/entity"
	"time"
)

const (
	defaultFileDaoCacheExpireTime = 7 * 24 * time.Hour
)

type fileDao struct {
	impl dao.IFileDao
}

func NewFileDao(impl dao.IFileDao) dao.IFileDao {
	return &fileDao{
		impl: impl,
	}
}

func (f *fileDao) CreateFileDraft(ctx context.Context, req *entity.CreateFileDraftRequest) (*entity.CreateFileDraftResponse, error) {
	return f.impl.CreateFileDraft(ctx, req)
}

func (f *fileDao) MarkFileReady(ctx context.Context, req *entity.MarkFileReadyRequest) (*entity.MarkFileReadyResponse, error) {
	return f.impl.MarkFileReady(ctx, req)
}

func (f *fileDao) buildCacheKey(fid uint64) string {
	return fmt.Sprintf("tgfile:cache:fileid:%d", fid)
}

func (f *fileDao) GetFileInfo(ctx context.Context, req *entity.GetFileInfoRequest) (*entity.GetFileInfoResponse, error) {
	keys := make([]string, 0, len(req.FileIds))
	mapping := make(map[string]uint64, len(req.FileIds))
	for _, fid := range req.FileIds {
		key := f.buildCacheKey(fid)
		keys = append(keys, key)
		mapping[key] = fid
	}
	caacheRs, err := cache.LoadMany(ctx, keys, func(ctx context.Context, c cache.ICache, ks []string) (map[string]interface{}, error) {
		fids := make([]uint64, 0, len(ks))
		for _, k := range ks {
			fids = append(fids, mapping[k])
		}
		rs, err := f.impl.GetFileInfo(ctx, &entity.GetFileInfoRequest{
			FileIds: fids,
		})
		if err != nil {
			return nil, err
		}
		ret := make(map[string]interface{}, len(rs.List))
		for _, item := range rs.List {
			k := f.buildCacheKey(item.FileId)
			ret[k] = item
			_ = c.Set(ctx, k, item, defaultFileDaoCacheExpireTime)
		}
		return ret, nil
	})
	if err != nil {
		return nil, err
	}
	rsp := &entity.GetFileInfoResponse{}
	for _, fid := range req.FileIds {
		v, ok := caacheRs[f.buildCacheKey(fid)]
		if !ok {
			continue
		}
		rsp.List = append(rsp.List, v.(*entity.GetFileInfoItem))
	}
	return rsp, nil
}
