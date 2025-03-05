package filemgr

import (
	"context"
	"fileserver/blockio"
	"fmt"
	"io"
)

// BlockIdToFileKeyConvertFunc 实现blockid到文件key的转换, 之后seeker会使用filesystem再去获取文件流
type BlockIdToFileKeyConvertFunc func(ctx context.Context, blkid int32) (string, error)

// fakeReader 由于底层的reader并不是真的可以seek,
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

type defaultReadSeekCloser struct {
	ctx    context.Context
	fs     blockio.IBlockIO
	b2f    BlockIdToFileKeyConvertFunc
	fsize  int64
	cur    int64 //for recording current read pos
	isOpen bool
	rc     io.ReadCloser
}

func newReadSeekCloser(ctx context.Context, bkio blockio.IBlockIO, b2f BlockIdToFileKeyConvertFunc, fsize int64) io.ReadSeekCloser {
	return &defaultReadSeekCloser{
		ctx:    ctx,
		fs:     bkio,
		b2f:    b2f,
		fsize:  fsize,
		cur:    0,
		isOpen: true,
	}
}

func (s *defaultReadSeekCloser) openStream(at int64) (io.ReadCloser, error) {
	if at == s.fsize {
		return defaultFakeReader, nil
	}
	//计算出在哪个块
	//计算出在块内的位置
	blkid := at / s.fs.MaxFileSize()
	pos := at % s.fs.MaxFileSize()
	filekey, err := s.b2f(s.ctx, int32(blkid))
	if err != nil {
		return nil, fmt.Errorf("unable to convert blockid:%d to fileid, err:%w", blkid, err)
	}
	stream, err := s.fs.Download(s.ctx, filekey, pos)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (s *defaultReadSeekCloser) Seek(offset int64, whence int) (ret int64, err error) {
	if !s.isOpen {
		return 0, fmt.Errorf("file not in open state")
	}
	if s.rc != nil {
		_ = s.rc.Close()
	}
	cur := s.calcOffset(offset, whence)
	if cur < 0 {
		return 0, fmt.Errorf("invalid offset, cur:%d", cur)
	}
	if cur > s.fsize {
		return s.fsize, fmt.Errorf("seek over file size, cur:%d, fsz:%d", cur, s.fsize)
	}
	if cur == 0 { //对于cur == 0的, 延迟到Read的时候才打开流。
		s.rc = nil
		s.cur = 0
		return 0, nil
	}
	rc, err := s.openStream(cur)
	if err != nil {
		return 0, fmt.Errorf("open stream fail, err:%w", err)
	}
	s.rc = rc
	s.cur = cur
	return cur, nil
}

func (s *defaultReadSeekCloser) calcOffset(offset int64, whence int) int64 {
	cur := int64(s.cur)
	switch whence {
	case io.SeekStart:
		cur = offset
	case io.SeekCurrent:
		cur += offset
	case io.SeekEnd:
		cur = s.fsize + offset
	}
	return cur
}

func (s *defaultReadSeekCloser) Close() error {
	if !s.isOpen {
		return nil
	}
	if s.rc == nil {
		return nil
	}
	s.isOpen = false
	return s.rc.Close()
}

func (s *defaultReadSeekCloser) Read(b []byte) (int, error) {
	if !s.isOpen {
		return 0, fmt.Errorf("file not in open state")
	}
	if s.rc == nil {
		rc, err := s.openStream(s.cur)
		if err != nil {
			return 0, fmt.Errorf("open stream fail, err:%w", err)
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
