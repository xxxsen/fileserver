package proxyutil

import "context"

type S3Info struct {
	Bucket string
	FileID string
}

type s3InfoType struct{}

var (
	defaultS3InfoKey = s3InfoType{}
)

func GetS3Info(ctx context.Context) (*S3Info, bool) {
	c, ok := ctx.Value(defaultS3InfoKey).(*S3Info)
	if !ok {
		return nil, false
	}
	return c, true
}

func SetS3Info(ctx context.Context, info *S3Info) context.Context {
	return context.WithValue(ctx, defaultS3InfoKey, info)
}
