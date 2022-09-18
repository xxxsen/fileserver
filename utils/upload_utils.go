package utils

import (
	"encoding/base64"
	"fileserver/proto/fileserver/fileinfo"

	"github.com/xxxsen/common/errs"

	"google.golang.org/protobuf/proto"
)

func EncodePartPair(pctx *fileinfo.PartPair) ([]byte, error) {
	return encodeMessageRaw(pctx)
}

func DecodePartPair(raw []byte) (*fileinfo.PartPair, error) {
	ctx := &fileinfo.PartPair{}
	if err := decodeMessageRaw(raw, ctx); err != nil {
		return nil, err
	}
	return ctx, nil
}

func EncodeBotUploadContext(bctx *fileinfo.BotUploadContext) ([]byte, error) {
	return encodeMessageRaw(bctx)
}

func DecodeBotUploadContext(raw []byte) (*fileinfo.BotUploadContext, error) {
	ctx := &fileinfo.BotUploadContext{}
	if err := decodeMessageRaw(raw, ctx); err != nil {
		return nil, err
	}
	return ctx, nil
}

func EncodeBotFileExtra(bctx *fileinfo.BotFileExtra) ([]byte, error) {
	return encodeMessageRaw(bctx)
}

func DecodeBotFileExtra(raw []byte) (*fileinfo.BotFileExtra, error) {
	ctx := &fileinfo.BotFileExtra{}
	if err := decodeMessageRaw(raw, ctx); err != nil {
		return nil, err
	}
	return ctx, nil
}

func EncodeUploadID(upload *fileinfo.UploadIdCtx) (string, error) {
	return encodeMessage(upload)
}

func DecodeUploadID(id string) (*fileinfo.UploadIdCtx, error) {
	ctx := &fileinfo.UploadIdCtx{}
	if err := decodeMessage(id, ctx); err != nil {
		return nil, err
	}
	return ctx, nil
}

func encodeMessageRaw(msg proto.Message) ([]byte, error) {
	raw, err := proto.Marshal(msg)
	if err != nil {
		return nil, errs.Wrap(errs.ErrMarshal, "encode pb msg fail", err)
	}
	return raw, nil
}

func encodeMessage(msg proto.Message) (string, error) {
	raw, err := encodeMessageRaw(msg)
	if err != nil {
		return "", errs.Wrap(errs.ErrMarshal, "pb marshal fail", err)
	}
	return base64.StdEncoding.EncodeToString(raw), nil
}

func decodeMessageRaw(raw []byte, dst proto.Message) error {
	if err := proto.Unmarshal(raw, dst); err != nil {
		return errs.Wrap(errs.ErrUnmarshal, "proto decode fail", err)
	}
	return nil
}

func decodeMessage(id string, dst proto.Message) error {
	raw, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return errs.Wrap(errs.ErrUnmarshal, "base64 decode fail", err)
	}
	return decodeMessageRaw(raw, dst)
}

func CalcFileBlockCount(sz uint64, blksz uint64) int {
	return int((sz + blksz - 1) / blksz)
}

func CalcBlockSize(sz uint64, blksz uint64, blkid int) uint64 {
	blkcnt := CalcFileBlockCount(sz, blksz)
	if blkid >= blkcnt || blkid < 0 {
		return 0
	}
	if blkid < blkcnt-1 {
		return blksz
	}
	return sz - blksz*(uint64(blkcnt)-1)
}
