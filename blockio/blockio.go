package blockio

import (
	"context"
	"io"
)

// IBlockIO 文件系统仅做简单的上传下载操作, 指定位置下载等能力交给外部实现
// 这样后续扩展/调试都会相对容易
type IBlockIO interface {
	MaxFileSize() int64
	Upload(ctx context.Context, r io.Reader) (string, error)
	Download(ctx context.Context, filekey string, pos int64) (io.ReadCloser, error)
}
