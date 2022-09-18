package bot

import (
	"crypto/md5"
	"hash"
	"io"
)

type md5Reader struct {
	ck hash.Hash
	r  io.Reader
}

func MD5Reader(r io.Reader) *md5Reader {
	return &md5Reader{r: r, ck: md5.New()}
}

func (r *md5Reader) Read(out []byte) (int, error) {
	cnt, err := r.r.Read(out)
	if cnt > 0 {
		r.ck.Write(out[:cnt])
	}
	if err != nil {
		return cnt, err
	}
	return cnt, nil
}

func (r *md5Reader) GetSum() []byte {
	return r.ck.Sum(nil)
}

type countReader struct {
	r   io.Reader
	cnt int
}

func CountReader(r io.Reader) *countReader {
	return &countReader{r: r}
}

func (r *countReader) Read(out []byte) (int, error) {
	cnt, err := r.r.Read(out)
	if cnt > 0 {
		r.cnt += cnt
	}
	if err != nil {
		return cnt, err
	}
	return cnt, nil
}

func (r *countReader) GetCount() int {
	return r.cnt
}
