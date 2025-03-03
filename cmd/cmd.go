package main

import (
	_ "fileserver/auth"
	"fileserver/config"
	"fileserver/db"
	"fileserver/filesystem"
	"fileserver/filesystem/telegram"
	"fileserver/server"
	"flag"

	"github.com/xxxsen/common/idgen"
	"github.com/xxxsen/common/logger"
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
	if err := idgen.Init(1); err != nil {
		logger.Fatal("init idgen fail", zap.Error(err))
	}
	logger.Info("recv config", zap.Any("config", c))
	if err := db.InitDB(c.DBFile); err != nil {
		logger.Fatal("init media db fail", zap.Error(err))
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
	tfs, err := telegram.New(int64(c.BotInfo.Chatid), c.BotInfo.Token)
	if err != nil {
		return err
	}
	filesystem.SetFileSystem(tfs)
	return nil
}
