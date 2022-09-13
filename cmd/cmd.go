package main

import (
	"fileserver/config"
	"fileserver/db"
	"fileserver/handler"
	"flag"

	"github.com/xxxsen/common/logger"
	"github.com/xxxsen/common/naivesvr"
	"github.com/xxxsen/common/s3"
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
	if err := db.InitFileDB(&c.FileDBInfo); err != nil {
		logger.With(zap.Error(err)).Fatal("init media db fail")
	}
	if err := s3.InitGlobal(
		s3.WithEndpoint(c.S3Info.Endpoint),
		s3.WithSSL(c.S3Info.UseSSL),
		s3.WithSecret(c.S3Info.SecretId, c.S3Info.SecretKey),
		s3.WithBucket(c.S3Info.Bucket),
	); err != nil {
		logger.With(zap.Error(err)).Fatal("init s3 fail")
	}
	svr, err := naivesvr.NewServer(
		naivesvr.WithAddress(c.ServerInfo.Address),
		naivesvr.WithHandlerRegister(handler.OnRegist),
	)
	if err != nil {
		logger.With(zap.Error(err)).Fatal("init server fail")
	}
	if err := svr.Run(); err != nil {
		logger.With(zap.Error(err)).Fatal("run server fail")
	}
}
