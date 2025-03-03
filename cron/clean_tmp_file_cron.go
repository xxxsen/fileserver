package cron

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"context"
)

const (
	defaultCleanTempFileCronSpec = "0 */1 * * * *" //clean expire file per hour
)

type cleanTempFileCron struct {
	dir  string
	keep time.Duration
}

func NewCleanTempFileCron(dir string, keep time.Duration) ICronJob {
	return &cleanTempFileCron{
		dir:  dir,
		keep: keep,
	}
}

func (c *cleanTempFileCron) Name() string {
	return "clean_temp_file"
}

func (c *cleanTempFileCron) Expression() string {
	return defaultCleanTempFileCronSpec
}

func (c *cleanTempFileCron) Run(ctx context.Context) error {
	enties, err := os.ReadDir(c.dir)
	if err != nil {
		return fmt.Errorf("open dir fail, err:%w", err)
	}
	if len(enties) == 0 {
		return nil
	}
	now := time.Now()
	for _, item := range enties {
		info, err := item.Info()
		if err != nil {
			return fmt.Errorf("read file info fail, file:%s, err:%w", info.Name(), err)
		}
		if !info.ModTime().Add(c.keep).Before(now) {
			continue
		}
		log.Printf("file expire, clean it, dir:%s, name:%s", c.dir, info.Name())
		_ = os.Remove(fmt.Sprintf("%s%s%s", c.dir, string(filepath.Separator), info.Name()))
	}
	return nil
}
