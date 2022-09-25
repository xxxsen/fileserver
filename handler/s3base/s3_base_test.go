package s3base

import (
	"bytes"
	"context"
	"fileserver/utils"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	s3c "github.com/xxxsen/common/s3"
)

func TestS3UploadDownload(t *testing.T) {
	client, err := s3c.New(s3c.WithBucket("s3"), s3c.WithEndpoint("http://127.0.0.1:9901"), s3c.WithSSL(false), s3c.WithSecret("abc", "123"))
	assert.NoError(t, err)
	ctx := context.Background()
	content := "hello world, this is a test file xxxx for s3 upload"
	fileid := "aaaa/1234567.txt"
	cks, err := client.Upload(ctx, fileid, bytes.NewReader([]byte(content)), int64(len(content)), utils.GetMd5([]byte(content)))
	assert.NoError(t, err)
	t.Logf("upload suss, cks:%s", cks)

	r, err := client.Download(ctx, fileid)
	assert.NoError(t, err)
	raw, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	t.Logf("download succ, data:%s", string(raw))
	assert.Equal(t, content, string(raw))
}
