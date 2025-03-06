package mem

import (
	"bytes"
	"context"
	"fileserver/blockio"
	"fmt"
	"io"
	"sync"

	"github.com/google/uuid"
)

type memBlockIO struct {
	m sync.Map
}

func (m *memBlockIO) MaxFileSize() int64 {
	return 4 * 1024
}

func (m *memBlockIO) Upload(ctx context.Context, r io.Reader) (string, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	key := uuid.NewString()
	m.m.Store(key, raw)
	return key, nil
}

func (m *memBlockIO) Download(ctx context.Context, filekey string, pos int64) (io.ReadCloser, error) {
	raw, ok := m.m.Load(filekey)
	if !ok {
		return nil, fmt.Errorf("key not found")
	}
	data := raw.([]byte)
	if pos > int64(len(data)) {
		pos = int64(len(data))
	}
	return io.NopCloser(bytes.NewReader(data[pos:])), nil
}

func New() blockio.IBlockIO {
	return &memBlockIO{}
}
