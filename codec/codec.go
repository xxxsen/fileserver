package codec

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func EncodeID(id uint64) (string, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, id)
	return hex.EncodeToString(buf), nil
}

func MustEncodeID(id uint64) string {
	hid, err := EncodeID(id)
	if err != nil {
		panic(err)
	}
	return hid
}

func DecodeID(id string) (uint64, error) {
	if len(id) != 16 {
		return 0, fmt.Errorf("invalid id")
	}
	raw, err := hex.DecodeString(id)
	if err != nil {
		return 0, fmt.Errorf("decode id failed, err:%w", err)
	}
	return binary.BigEndian.Uint64(raw), nil
}
