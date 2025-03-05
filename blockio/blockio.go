package blockio

import (
	"context"
	"io"
)

var (
	defaultInst IBlockIO
)

// IBlockIO 文件系统仅做简单的上传下载操作, 指定位置下载等能力交给外部实现
// 这样后续扩展/调试都会相对容易
type IBlockIO interface {
	MaxFileSize() int64
	Upload(ctx context.Context, r io.Reader) (string, error)
	Download(ctx context.Context, filekey string, pos int64) (io.ReadCloser, error)
}

func SetBlockIO(fs IBlockIO) {
	defaultInst = fs
}

func MaxFileSize() int64 {
	return defaultInst.MaxFileSize()
}

func Upload(ctx context.Context, r io.Reader) (string, error) {
	return defaultInst.Upload(ctx, r)
}

func Download(ctx context.Context, filekey string, pos int64) (io.ReadCloser, error) {
	return defaultInst.Download(ctx, filekey, pos)
}
