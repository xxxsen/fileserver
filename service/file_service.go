package service

import (
	"context"
	"fileserver/constant"
	"fileserver/dao"
	"fileserver/entity"
	"fmt"
)

var FileService = newFileService()

type fileService struct {
}

func newFileService() *fileService {
	return &fileService{}
}

func (s *fileService) CreateFileDraft(ctx context.Context, filename string, filesize int64, filepartcount int32) (uint64, error) {
	rs, err := dao.FileDao.CreateFileDraft(ctx, &entity.CreateFileDraftRequest{
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
	if _, err := dao.FilePartDao.CreateFilePart(ctx, &entity.CreateFilePartRequest{
		FileId:     fileid,
		FilePartId: pidx,
		FileKey:    filekey,
	}); err != nil {
		return err
	}
	return nil
}

func (s *fileService) FinishCreateFile(ctx context.Context, fileid uint64) error {
	info, ok, err := s.GetFileInfo(ctx, fileid)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("fileid:%d not found", fileid)
	}
	if info.FileState != constant.FileStateInit {
		return fmt.Errorf("file not in init state")
	}
	partCount, err := s.GetFilePartCount(ctx, fileid)
	if err != nil {
		return fmt.Errorf("read file part count failed, err:%w", err)
	}
	if partCount != info.FilePartCount {
		return fmt.Errorf("file part count not match, db count:%d, acquire count:%d", partCount, info.FilePartCount)
	}
	if _, err := dao.FileDao.MarkFileReady(ctx, &entity.MarkFileReadyRequest{
		FileID: fileid,
	}); err != nil {
		return err
	}
	return nil
}

func (s *fileService) GetFileInfo(ctx context.Context, fileid uint64) (*entity.GetFileInfoItem, bool, error) {
	rs, err := dao.FileDao.GetFileInfo(ctx, &entity.GetFileInfoRequest{
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
	rs, err := dao.FileDao.GetFileInfo(ctx, &entity.GetFileInfoRequest{
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
	rs, err := dao.FilePartDao.GetFilePartInfo(ctx, &entity.GetFilePartInfoRequest{
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

func (s *fileService) GetFilePartCount(ctx context.Context, fileid uint64) (int32, error) {
	rs, err := dao.FilePartDao.GetFilePartCount(ctx, &entity.GetFilePartCountRequest{FileId: fileid})
	if err != nil {
		return 0, err
	}
	return rs.Count, nil
}
