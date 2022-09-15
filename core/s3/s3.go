package s3

import (
	"context"
	"fileserver/core"
	"fileserver/proto/fileserver/fileinfo"
	"fileserver/utils"
	"fmt"
	"io"

	"github.com/xxxsen/common/errs"
	"google.golang.org/protobuf/proto"
)

const (
	defaultMaxS3FileSize   = 4 * 1024 * 1024 * 1024
	defaultS3FileBlockSize = 200 * 1024 * 1024
)

type S3Core struct {
	c *config
}

func New(opts ...Option) (*S3Core, error) {
	c := &config{
		fsize:   defaultMaxS3FileSize,
		blksize: defaultS3FileBlockSize,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.client == nil {
		return nil, fmt.Errorf("nil client")
	}
	if c.idg == nil {
		return nil, fmt.Errorf("nil idgenerator")
	}
	return &S3Core{c: c}, nil
}

func (c *S3Core) BlockSize() int64 {
	return defaultS3FileBlockSize
}

func (c *S3Core) MaxFileSize() int64 {
	return defaultMaxS3FileSize
}

func (c *S3Core) FileUpload(ctx context.Context, uctx *core.FileUploadRequest) (string, error) {
	fid := c.c.idg.NextId()
	xfid := utils.EncodeFileId(fid)
	if err := c.c.client.Upload(ctx, xfid, uctx.ReadSeeker, uctx.Size, uctx.MD5); err != nil {
		return "", err
	}
	return xfid, nil
}

func (c *S3Core) FileDownload(ctx context.Context, fctx *core.FileDownloadRequest) (io.ReadCloser, error) {
	body, err := c.c.client.DownloadByRange(ctx, fctx.Key, fctx.StartAt)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *S3Core) BeginFileUpload(ctx context.Context, fctx *core.BeginFileUploadRequest) (string, error) {
	fid := c.c.idg.NextId()
	xfid := utils.EncodeFileId(fid)
	uploadid, err := c.c.client.BeginUpload(ctx, xfid)
	if err != nil {
		return "", err
	}
	upid, err := utils.EncodeUploadID(&fileinfo.UploadIdCtx{
		FileSize:  proto.Uint64(uint64(fctx.FileSize)),
		UploadId:  proto.String(uploadid),
		DownKey:   proto.String(xfid),
		BlockSize: proto.Uint32(uint32(c.BlockSize())),
	})
	if err != nil {
		return "", err
	}
	return upid, nil
}

func (c *S3Core) PartFileUpload(ctx context.Context, pctx *core.PartFileUploadRequest) error {
	uctx, err := utils.DecodeUploadID(pctx.UploadId)
	if err != nil {
		return err
	}
	bkcnt := utils.CalcFileBlockCount(uctx.GetFileSize(), uint64(uctx.GetBlockSize()))
	if pctx.PartId == 0 || pctx.PartId > uint64(bkcnt) {
		return errs.New(errs.ErrParam, "invalid partid:%d", pctx.PartId)
	}
	if pctx.PartId != uint64(bkcnt) && pctx.Size != int64(uctx.GetBlockSize()) {
		return errs.New(errs.ErrParam, "invalid part size, partid:%d, blksize:%d", pctx.PartId, uctx.GetBlockSize())
	}
	if pctx.Size == 0 {
		return errs.New(errs.ErrParam, "empty size")
	}
	if pctx.PartId == uint64(bkcnt) {
		if (pctx.PartId-1)*uint64(uctx.GetBlockSize())+uint64(pctx.Size) != uctx.GetFileSize() {
			return errs.New(errs.ErrParam, "invalid file size, calc:%d, real:%d",
				(pctx.PartId-1)*uint64(uctx.GetBlockSize())+uint64(pctx.Size),
				uctx.GetFileSize())
		}
	}
	if err := c.c.client.UploadPart(ctx, uctx.GetDownKey(), uctx.GetUploadId(), int(pctx.PartId), pctx.ReadSeeker, pctx.MD5); err != nil {
		return err
	}
	return nil
}

func (c *S3Core) FinishFileUpload(ctx context.Context, fctx *core.FinishFileUploadRequest) (string, error) {
	uctx, err := utils.DecodeUploadID(fctx.UploadId)
	if err != nil {
		return "", err
	}
	if err := c.c.client.EndUpload(ctx, uctx.GetDownKey(), uctx.GetUploadId(),
		utils.CalcFileBlockCount(uctx.GetFileSize(), uint64(uctx.GetBlockSize()))); err != nil {
		return "", err
	}
	return uctx.GetDownKey(), nil
}
