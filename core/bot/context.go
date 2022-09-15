package bot

import (
	"bytes"
	"encoding/gob"
)

const (
	fileTypeOneFile    = 0
	fileTypeMultiBlock = 1
)

type botFileCtx struct {
	ChatId   int64
	FileType int64
}

func encodeFileExtra(fc *botFileCtx) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(fc); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeFileExtra(raw []byte) (*botFileCtx, error) {
	reader := bytes.NewReader(raw)
	dec := gob.NewDecoder(reader)
	fc := &botFileCtx{}
	if err := dec.Decode(fc); err != nil {
		return nil, err
	}
	return fc, nil
}
