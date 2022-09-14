package s3

import (
	"github.com/xxxsen/common/idgen"
	s3c "github.com/xxxsen/common/s3"
)

type config struct {
	client  *s3c.S3Client
	idg     idgen.IDGenerator
	fsize   int64
	blksize int64
}

type Option func(c *config)

func WithS3Client(client *s3c.S3Client) Option {
	return func(c *config) {
		c.client = client
	}
}

func WithIDGen(idg idgen.IDGenerator) Option {
	return func(c *config) {
		c.idg = idg
	}
}

func WithSizeLimit(fsize int64, blksize int64) Option {
	return func(c *config) {
		c.fsize = fsize
		c.blksize = blksize
	}
}
