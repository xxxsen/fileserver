package filemgr

import (
	"context"
	"fileserver/blockio"
	"fileserver/service"
	"fileserver/utils"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"time"
)

var defaultFileMgr IFileManager

type IFileManager interface {
	Stat(ctx context.Context, fileid uint64) (fs.FileInfo, error)
	Open(ctx context.Context, fileid uint64) (io.ReadSeekCloser, error)
	Create(ctx context.Context, name string, size int64, r io.Reader) (uint64, error)
}

func SetFileManagerImpl(mgr IFileManager) {
	defaultFileMgr = mgr
}

func Stat(ctx context.Context, fileid uint64) (fs.FileInfo, error) {
	return defaultFileMgr.Stat(ctx, fileid)
}

func Open(ctx context.Context, fileid uint64) (io.ReadSeekCloser, error) {
	return defaultFileMgr.Open(ctx, fileid)
}

func Create(ctx context.Context, name string, size int64, r io.Reader) (uint64, error) {
	return defaultFileMgr.Create(ctx, name, size, r)
}

type defaultFileManager struct {
	bkio blockio.IBlockIO
}

func (d *defaultFileManager) Open(ctx context.Context, fileid uint64) (io.ReadSeekCloser, error) {
	finfo, ok, err := service.FileService.GetFileInfo(ctx, fileid)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("file not found")
	}
	rsc := newReadSeekCloser(ctx, d.bkio, func(ctx context.Context, blkid int32) (string, error) {
		pinfo, ok, err := service.FileService.GetFilePartInfo(ctx, fileid, blkid)
		if err != nil {
			return "", fmt.Errorf("read file part info failed, err:%w", err)
		}
		if !ok {
			return "", fmt.Errorf("partid:%d not found", blkid)
		}
		return pinfo.FileKey, nil
	}, finfo.FileSize)
	return rsc, nil
}

func (d *defaultFileManager) Create(ctx context.Context, filename string, size int64, reader io.Reader) (uint64, error) {
	blkcnt := utils.CalcFileBlockCount(uint64(size), uint64(d.bkio.MaxFileSize()))
	fileid, err := service.FileService.CreateFileDraft(ctx, filename, size, int32(blkcnt))
	if err != nil {
		return 0, fmt.Errorf("create file draft failed, err:%w", err)
	}
	for i := 0; i < blkcnt; i++ {
		r := io.LimitReader(reader, d.bkio.MaxFileSize())
		fileKey, err := d.bkio.Upload(ctx, r)
		if err != nil {
			return 0, fmt.Errorf("upload part failed, err:%w", err)
		}
		if err := service.FileService.CreateFilePart(ctx, fileid, int32(i), fileKey); err != nil {
			return 0, fmt.Errorf("create part record failed, err:%w", err)
		}
	}

	if err := service.FileService.FinishCreateFile(ctx, fileid); err != nil {
		return 0, fmt.Errorf("finish create file failed, err:%w", err)
	}
	return fileid, nil
}

func (d *defaultFileManager) Stat(ctx context.Context, fileid uint64) (fs.FileInfo, error) {
	finfo, ok, err := service.FileService.GetFileInfo(ctx, fileid)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("file not found")
	}
	return newFileInfo(filepath.Base(finfo.FileName), finfo.FileSize, finfo.Mtime), nil
}

func newFileInfo(name string, sz int64, mtime int64) fs.FileInfo {
	return &defaultFileInfo{
		name:  name,
		sz:    sz,
		mtime: time.UnixMilli(mtime),
	}
}

type defaultFileInfo struct {
	name  string
	sz    int64
	mtime time.Time
}

func (d *defaultFileInfo) Name() string {
	return d.name
}

func (d *defaultFileInfo) Size() int64 {
	return d.sz
}

func (d *defaultFileInfo) Mode() fs.FileMode {
	return 0644
}

func (d *defaultFileInfo) ModTime() time.Time {
	return d.mtime
}

func (d *defaultFileInfo) IsDir() bool {
	return false
}

func (d *defaultFileInfo) Sys() interface{} {
	return nil

}

func NewFileManager(bkio blockio.IBlockIO) IFileManager {
	return &defaultFileManager{
		bkio: bkio,
	}
}
