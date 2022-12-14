package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/xxxsen/common/errs"
)

const (
	defaultScanCronSpec = "0 */1 * * *" //clean expire file per hour
)

var defaultFsCleaner = NewFileCleaner()

func AddCleanTask(tk *CleanEntry) {
	defaultFsCleaner.StartCleanTask(tk)
}

type CleanEntry struct {
	Dir  string
	Keep time.Duration
}

type FileCleaner struct {
	cr *cron.Cron
}

func NewFileCleaner() *FileCleaner {
	cr := cron.New()
	cr.Start()
	return &FileCleaner{
		cr: cr,
	}
}

func (f *FileCleaner) StartCleanTask(tk *CleanEntry) {
	log.Printf("start clean task at dir:%s", tk.Dir)
	f.cr.AddFunc(defaultScanCronSpec, func() {
		f.worker(tk)
	})
}

func (f *FileCleaner) worker(tk *CleanEntry) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic recover, dir:%s, rec:%v, stack:%s", tk.Dir, err, string(debug.Stack()))
		}
	}()
	if err := f.doWork(tk); err != nil {
		log.Printf("do scan task fail, dir:%s, err:%v", tk.Dir, err)
	}
}

func (f *FileCleaner) doWork(tk *CleanEntry) error {
	enties, err := os.ReadDir(tk.Dir)
	if err != nil {
		return errs.Wrap(errs.ErrIO, "open dir fail", err)
	}
	if len(enties) == 0 {
		return nil
	}
	now := time.Now()
	for _, item := range enties {
		info, err := item.Info()
		if err != nil {
			return errs.Wrap(errs.ErrUnknown, fmt.Sprintf("read file info fail, file:%s", info.Name()), err)
		}
		if !info.ModTime().Add(tk.Keep).Before(now) {
			continue
		}
		log.Printf("file expire, clean it, dir:%s, name:%s", tk.Dir, info.Name())
		_ = os.Remove(fmt.Sprintf("%s%s%s", tk.Dir, string(filepath.Separator), info.Name()))
	}
	return nil
}
