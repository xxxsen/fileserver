package s3base

import (
	"bytes"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignature(t *testing.T) {
	body := "hello world, this is a test file xxxx for s3 upload"
	req, err := http.NewRequest(http.MethodPut, "http://127.0.0.1:9901/s3/aaaa/1234567.txt", bytes.NewReader([]byte(body)))
	assert.NoError(t, err)
	req.Header.Set("X-Amz-Content-Sha256", "54579d254b79513a2fe0b977af1a5afdd030dfeaf2c87ad5b66465ae265891c1")
	req.Header.Set("Content-Length", strconv.FormatInt(int64(len(body)), 10))
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential=abc/20221004/cn/s3/aws4_request, SignedHeaders=content-length;content-md5;host;x-amz-content-sha256;x-amz-date, Signature=cf3ad29c3243d6935386fc60cdf35b1d4d593de8639875d14d9ffacfc7b77851")
	req.Header.Set("Content-Md5", "usc6jTToHiJAUYOfOGnbGQ==")
	req.Header.Set("X-Amz-Date", "20221004T155011Z")
	req.Header.Set("Host", "127.0.0.1:9901")

	isV4Sign := IsRequestSignatureV4(req)
	assert.True(t, isV4Sign)
	sign, exist, err := ParseV4Signature(req)
	assert.NoError(t, err)
	assert.True(t, exist)
	assert.Equal(t, sign.AKey, "abc")
	assert.Equal(t, sign.Algorithm, v4SignAlgorithm)
	assert.Equal(t, sign.Contentmd5, "usc6jTToHiJAUYOfOGnbGQ==")
	assert.Equal(t, sign.Contentsha256, "54579d254b79513a2fe0b977af1a5afdd030dfeaf2c87ad5b66465ae265891c1")
	assert.Equal(t, sign.Date, "20221004T155011Z")
	assert.Equal(t, sign.Region, "cn")
	assert.Equal(t, sign.RequestType, "aws4_request")
	assert.Equal(t, sign.Service, "s3")
	assert.Equal(t, sign.SignedHeaders, []string{"content-length", "content-md5", "host", "x-amz-content-sha256", "x-amz-date"})
	assert.Equal(t, sign.Signature, "cf3ad29c3243d6935386fc60cdf35b1d4d593de8639875d14d9ffacfc7b77851")
	pass, err := S3AuthV4(req, "abc", "123456", sign)
	assert.NoError(t, err)
	assert.True(t, pass)
}
