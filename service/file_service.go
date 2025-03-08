package service

import (
	"context"
	"fileserver/dao"
	"fileserver/dao/cache"
	"fileserver/entity"
)

var FileService = newFileService()

type fileService struct {
	fileDao     dao.IFileDao
	filePartDao dao.IFilePartDao
}

func newFileService() *fileService {
	return &fileService{
		fileDao:     cache.NewFileDao(dao.FileDao),
		filePartDao: cache.NewFilePartDao(dao.FilePartDao),
	}
}

func (s *fileService) CreateFileDraft(ctx context.Context, filename string, filesize int64, filepartcount int32) (uint64, error) {
	rs, err := s.fileDao.CreateFileDraft(ctx, &entity.CreateFileDraftRequest{
		FileName:      filename,
		FileSize:      filesize,
		FilePartCount: filepartcount,
	})
	if err != nil {
		return 0, err
	}
	return rs.FileId, nil
}

func (s *fileService) CreateFilePart(ctx context.Context, fileid uint64, pidx int32, filekey string) error {
	if _, err := s.filePartDao.CreateFilePart(ctx, &entity.CreateFilePartRequest{
		FileId:     fileid,
		FilePartId: pidx,
		FileKey:    filekey,
	}); err != nil {
		return err
	}
	return nil
}

func (s *fileService) FinishCreateFile(ctx context.Context, fileid uint64) error {
	if _, err := s.fileDao.MarkFileReady(ctx, &entity.MarkFileReadyRequest{
		FileID: fileid,
	}); err != nil {
		return err
	}
	return nil
}

func (s *fileService) GetFileInfo(ctx context.Context, fileid uint64) (*entity.GetFileInfoItem, bool, error) {
	rs, err := s.fileDao.GetFileInfo(ctx, &entity.GetFileInfoRequest{
		FileIds: []uint64{fileid},
	})
	if err != nil {
		return nil, false, err
	}
	if len(rs.List) == 0 {
		return nil, false, nil
	}
	return rs.List[0], true, nil
}

func (s *fileService) BatchGetFileInfo(ctx context.Context, fileids []uint64) (map[uint64]*entity.GetFileInfoItem, error) {
	rs, err := s.fileDao.GetFileInfo(ctx, &entity.GetFileInfoRequest{
		FileIds: fileids,
	})
	if err != nil {
		return nil, err
	}
	mapper := make(map[uint64]*entity.GetFileInfoItem, len(rs.List))
	for _, item := range rs.List {
		mapper[item.FileId] = item
	}
	return mapper, nil
}

func (s *fileService) GetFilePartInfo(ctx context.Context, fileid uint64, partid int32) (*entity.GetFilePartInfoItem, bool, error) {
	rs, err := s.filePartDao.GetFilePartInfo(ctx, &entity.GetFilePartInfoRequest{
		FileId:     fileid,
		FilePartId: []int32{partid},
	})
	if err != nil {
		return nil, false, err
	}
	if len(rs.List) == 0 {
		return nil, false, nil
	}
	return rs.List[0], true, nil
}
