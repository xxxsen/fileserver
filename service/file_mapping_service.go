package service

import (
	"context"
	"tgfile/dao"
	"tgfile/dao/cache"
	"tgfile/entity"
)

var FileMappingService = newFileMappingService()

type IterMappingFunc func(ctx context.Context, filename string, fileid uint64) (bool, error)

type fileMappingService struct {
	fileMappingDao dao.IFileMappingDao
}

func newFileMappingService() *fileMappingService {
	return &fileMappingService{
		fileMappingDao: cache.NewFileMappingDao(dao.NewFileMappingDao()),
	}
}

func (s *fileMappingService) GetFileMapping(ctx context.Context, filename string) (uint64, bool, error) {
	rsp, ok, err := s.fileMappingDao.GetFileMapping(ctx, &entity.GetFileMappingRequest{
		FileName: filename,
	})
	if err != nil {
		return 0, false, err
	}
	if !ok {
		return 0, false, nil
	}
	return rsp.Item.FileId, true, nil
}

func (s *fileMappingService) CreateFileMapping(ctx context.Context, filename string, fileid uint64) error {
	_, err := s.fileMappingDao.CreateFileMapping(ctx, &entity.CreateFileMappingRequest{
		FileName: filename,
		FileId:   fileid,
	})
	return err
}

func (s *fileMappingService) IterFileMapping(ctx context.Context, cb IterMappingFunc) error {
	return s.fileMappingDao.IterFileMapping(ctx, func(ctx context.Context, name string, fileid uint64) (bool, error) {
		return cb(ctx, name, fileid)
	})
}
