package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"

	"github.com/xxxsen/common/errs"
)

func Base64Md52HexMd5(ck string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(ck)
	if err != nil {
		return "", errs.Wrap(errs.ErrParam, "invalid b64 checksum", err)
	}
	return hex.EncodeToString(raw), nil
}

func GetMd5(raw []byte) string {
	h := md5.New()
	h.Write(raw)
	return hex.EncodeToString(h.Sum(nil))
}
