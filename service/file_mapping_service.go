package service

import (
	"context"
	"fileserver/dao"
	"fileserver/entity"
)

var FileMappingService = newFileMappingService()

type fileMappingService struct {
}

func newFileMappingService() *fileMappingService {
	return &fileMappingService{}
}

func (s *fileMappingService) GetFileMapping(ctx context.Context, filename string) (uint64, bool, error) {
	rsp, ok, err := dao.FileMappingDao.GetFileMapping(ctx, &entity.GetFileMappingRequest{
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
	_, err := dao.FileMappingDao.CreateFileMapping(ctx, &entity.CreateFileMappingRequest{
		FileName: filename,
		FileId:   fileid,
	})
	return err
}
