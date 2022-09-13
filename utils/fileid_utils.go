package utils

import (
	"encoding/binary"
	"encoding/hex"
)

func EncodeFileId(fileid uint64) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, fileid)
	return hex.EncodeToString(buf)
}
