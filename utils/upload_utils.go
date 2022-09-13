package utils

import (
	"encoding/base64"
	"fileserver/proto/fileserver/fileinfo"

	"github.com/xxxsen/common/errs"

	"google.golang.org/protobuf/proto"
)

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

func encodeMessage(msg proto.Message) (string, error) {
	raw, err := proto.Marshal(msg)
	if err != nil {
		return "", errs.Wrap(errs.ErrMarshal, "pb marshal fail", err)
	}
	return base64.StdEncoding.EncodeToString(raw), nil
}

func decodeMessage(id string, dst proto.Message) error {
	raw, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return errs.Wrap(errs.ErrUnmarshal, "base64 decode fail", err)
	}
	if err := proto.Unmarshal(raw, dst); err != nil {
		return errs.Wrap(errs.ErrUnmarshal, "proto decode fail", err)
	}
	return nil
}

func CalcFileBlockCount(sz uint64, blksz uint64) int {
	return int((sz + blksz - 1) / blksz)
}
