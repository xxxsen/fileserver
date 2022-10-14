package fssystem

import (
	"fileserver/dao"
)

type config struct {
	fs dao.FsSystemService
}

type Option func(c *config)

func WithSystem(fs dao.FsSystemService) Option {
	return func(c *config) {
		c.fs = fs
	}
}
