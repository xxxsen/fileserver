package service

import (
	"context"
	"fileserver/dao"
	"fileserver/dao/cache"
	"fileserver/entity"
)

var FileMappingService = newFileMappingService()

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
