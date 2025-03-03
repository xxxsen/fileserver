package main

import (
	"context"
	_ "fileserver/auth"
	"fileserver/config"
	"fileserver/cron"
	"fileserver/db"
	"fileserver/server"
	"fileserver/tgfile"
	"flag"
	"fmt"
	"time"

	"github.com/xxxsen/common/logger"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

var file = flag.String("config", "./config.json", "config file path")

func main() {
	flag.Parse()

	c, err := config.Parse(*file)
	if err != nil {
		panic(err)
	}
	logitem := c.LogInfo
	logger := logger.Init(logitem.File, logitem.Level, int(logitem.FileCount), int(logitem.FileSize), int(logitem.KeepDays), logitem.Console)

	logger.Info("recv config", zap.Any("config", c))
	if err := db.InitDB(c.DBFile); err != nil {
		logger.Fatal("init media db fail", zap.Error(err))
	}
	if err := initCron(c); err != nil {
		logger.Fatal("init cron job failed", zap.Error(err))
	}
	if err := initStorage(c); err != nil {
		logger.Fatal("init storage fail", zap.Error(err))
	}

	svr, err := server.New(c.ServerInfo.Address,
		server.WithS3Buckets(c.S3Bucket),
		server.WithUser(c.AuthInfo),
	)
	if err != nil {
		logger.Fatal("init server fail", zap.Error(err))
	}
	logger.Info("init server succ, start it...")
	if err := svr.Run(); err != nil {
		logger.Fatal("run server fail", zap.Error(err))
	}
}

func initStorage(c *config.Config) error {
	botfs, err := tgfile.New(
		tgfile.WithAuth(int64(c.BotInfo.Chatid), c.BotInfo.Token),
		tgfile.WithTmpDir(c.TempDir),
	)
	if err != nil {
		return err
	}
	tgfile.SetFileSystem(botfs)
	return nil
}

func initCron(c *config.Config) error {
	cr := cron.New()
	jobs := []cron.ICronJob{
		cron.NewCleanTempFileCron(c.TempDir, 7*24*time.Hour),
	}
	for _, job := range jobs {
		if err := cr.AddJob(job); err != nil {
			return fmt.Errorf("init cron job failed, name:%s, err:%w", job.Name(), err)
		}
		logutil.GetLogger(context.Background()).Info("init job succ", zap.String("name", job.Name()))
	}
	cr.Start()
	return nil
}
