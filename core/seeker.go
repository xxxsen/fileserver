package core

import (
	"context"
	"fmt"
	"io"

	"github.com/xxxsen/common/errs"
)

//fakeReader 由于地层的reader并不是真的可以seek,
// 很多场景下, seek_end只是为了获取文件大小, 所以, 我们可以产生一个假的seeker
type fakeReader struct {
}

func (f *fakeReader) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (f *fakeReader) Close() error {
	return nil
}

var defaultFakeReader = &fakeReader{}

type SeekCore struct {
	ctx    context.Context
	c      IFsCore
	key    string
	extra  []byte
	rc     io.ReadCloser
	cur    int64 //for recording current read pos
	isOpen bool
	fsz    int64
}

func NewSeeker(ctx context.Context, c IFsCore, sz int64, key string, extra []byte) *SeekCore {
	return &SeekCore{
		ctx:    ctx,
		c:      c,
		key:    key,
		extra:  extra,
		cur:    0,
		isOpen: true,
		fsz:    sz,
	}
}

func (s *SeekCore) openStream(at int64) (io.ReadCloser, error) {
	if at == s.fsz {
		return defaultFakeReader, nil
	}
	rsp, err := s.c.FileDownload(s.ctx, &FileDownloadRequest{
		Key:     s.key,
		Extra:   s.extra,
		StartAt: at,
	})
	if err != nil {
		return nil, err
	}
	return rsp.Reader, nil
}

func (s *SeekCore) Read(b []byte) (int, error) {
	if !s.isOpen {
		return 0, fmt.Errorf("file not in open state")
	}
	if s.rc == nil {
		rc, err := s.openStream(int64(s.cur))
		if err != nil {
			return 0, errs.Wrap(errs.ErrIO, "open stream fail", err)
		}
		s.rc = rc
	}

	cnt, err := s.rc.Read(b)
	if cnt > 0 {
		s.cur += int64(cnt)
	}
	if err != nil {
		return cnt, err
	}
	return cnt, nil
}

func (s *SeekCore) Close() error {
	if !s.isOpen {
		return nil
	}
	if s.rc == nil {
		return nil
	}
	s.isOpen = false
	return s.rc.Close()
}

func (s *SeekCore) calcOffset(offset int64, whence int) int64 {
	cur := int64(s.cur)
	switch whence {
	case io.SeekStart:
		cur = offset
	case io.SeekCurrent:
		cur += offset
	case io.SeekEnd:
		cur = s.fsz + offset
	}
	return cur
}

func (s *SeekCore) Seek(offset int64, whence int) (ret int64, err error) {
	if !s.isOpen {
		return 0, fmt.Errorf("file not in open state")
	}
	if s.rc != nil {
		_ = s.rc.Close()
	}
	cur := s.calcOffset(offset, whence)
	if cur < 0 {
		return 0, errs.New(errs.ErrParam, "invalid offset, cur:%d", cur)
	}
	if cur > s.fsz {
		return s.fsz, errs.New(errs.ErrParam, "seek over file size, cur:%d, fsz:%d", cur, s.fsz)
	}
	if cur == 0 { //对于cur == 0的, 延迟到Read的时候才打开流。
		s.rc = nil
		s.cur = 0
		return 0, nil
	}
	rc, err := s.openStream(cur)
	if err != nil {
		return 0, errs.Wrap(errs.ErrIO, "open stream fail", err)
	}
	s.rc = rc
	s.cur = cur
	return cur, nil
}
