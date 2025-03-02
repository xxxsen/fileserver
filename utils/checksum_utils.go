package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func Base64Md52HexMd5(ck string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(ck)
	if err != nil {
		return "", fmt.Errorf("invalid b64 checksum, err:%w", err)
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
		return "", fmt.Errorf("unable to open file, err:%w", err)
	}
	defer file.Close()
	return ReaderMd5(file)
}

func ReaderMd5(r io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", fmt.Errorf("calc file md5 fail, err:%w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
