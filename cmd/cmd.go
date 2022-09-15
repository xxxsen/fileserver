package main

import (
	"fileserver/config"
	"fileserver/constants"
	"fileserver/core"
	"fileserver/core/s3"
	"fileserver/db"
	"fileserver/handler"
	"flag"
	"strings"

	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/idgen"
	"github.com/xxxsen/common/logger"
	"github.com/xxxsen/common/naivesvr"
	s3c "github.com/xxxsen/common/s3"
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
	if err := idgen.Init(c.IDGenInfo.WorkerID); err != nil {
		logger.With(zap.Error(err)).Fatal("init idgen fail")
	}
	fs, err := initStorage(c)
	if err != nil {
		logger.With(zap.Error(err)).Fatal("init storage fail")
	}
	svr, err := naivesvr.NewServer(
		naivesvr.WithAddress(c.ServerInfo.Address),
		naivesvr.WithHandlerRegister(handler.OnRegist),
		naivesvr.WithAttach(constants.KeyStorageClient, fs),
	)
	if err != nil {
		logger.With(zap.Error(err)).Fatal("init server fail")
	}
	if err := svr.Run(); err != nil {
		logger.With(zap.Error(err)).Fatal("run server fail")
	}
}

func initStorage(c *config.Config) (core.IFsCore, error) {
	switch strings.ToLower(c.StorageType) {
	case "s3":
		client, err := s3c.New(
			s3c.WithEndpoint(c.S3Info.Endpoint),
			s3c.WithSSL(c.S3Info.UseSSL),
			s3c.WithSecret(c.S3Info.SecretId, c.S3Info.SecretKey),
			s3c.WithBucket(c.S3Info.Bucket),
		)
		if err != nil {
			return nil, errs.Wrap(errs.ErrStorage, "init s3 fail", err)
		}
		s3core, err := s3.New(
			s3.WithS3Client(client),
			s3.WithIDGen(idgen.Default()),
		)
		if err != nil {
			return nil, errs.Wrap(errs.ErrStorage, "init s3 core fail", err)
		}
		return s3core, nil
	}
	return nil, errs.New(errs.ErrParam, "unsupport storage type:%s", c.StorageType)
}
