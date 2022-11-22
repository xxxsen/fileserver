package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"

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

func FileMd5(f string) (string, error) {
	file, err := os.Open(f)
	if err != nil {
		return "", errs.Wrap(errs.ErrServiceInternal, "unable to open file", err)
	}
	defer file.Close()
	return ReaderMd5(file)
}

func ReaderMd5(r io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", errs.Wrap(errs.ErrIO, "calc file md5 fail", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
