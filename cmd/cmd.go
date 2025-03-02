package main

import (
	"context"
	_ "fileserver/auth"
	"fileserver/config"
	"fileserver/core"
	"fileserver/core/bot"
	"fileserver/core/s3"
	"fileserver/db"
	"fileserver/server"
	"flag"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/xxxsen/common/idgen"
	"github.com/xxxsen/common/logger"
	"github.com/xxxsen/common/logutil"
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
		logger.Fatal("init media db fail", zap.Error(err))
	}
	if err := idgen.Init(c.IDGenInfo.WorkerID); err != nil {
		logger.Fatal("init idgen fail", zap.Error(err))
	}
	fs, err := initStorage(c)
	if err != nil {
		logger.Fatal("init storage fail", zap.Error(err))
	}
	core.SetFsCore(fs)

	svr, err := server.New(c.ServerInfo.Address,
		server.WithS3Buckets(c.FakeS3Info.BucketList),
		server.WithUser(c.AuthInfo),
	)
	if err != nil {
		logger.Fatal("init server fail", zap.Error(err))
	}
	if err := svr.Run(); err != nil {
		logger.Fatal("run server fail", zap.Error(err))
	}
}

func initStorage(c *config.Config) (core.IFsCore, error) {
	names := make([]string, 0, len(c.FsInfo))
	var uploader core.IFsCore
	downloaders := make([]core.IFsCore, 0, len(c.FsInfo))
	uploadfsname := strings.ToLower(c.UploadFs)
	for name, param := range c.FsInfo {
		var c core.IFsCore
		var err error
		name := strings.ToLower(name)
		switch name {
		case "s3":
			c, err = initS3Core(param)
		case "tgbot":
			c, err = initMultiTGBotCore(param)
		}
		if err != nil {
			return nil, fmt.Errorf("init core:%s fail, err:%w", name, err)
		}
		names = append(names, name)
		downloaders = append(downloaders, c)
		if name == uploadfsname {
			uploader = c
		}
	}
	if uploader == nil {
		return nil, fmt.Errorf("upload fs not found, support only:%+v", names)
	}
	return core.NewMultiCore(uploader, downloaders...)
}

func decodeToType(src interface{}, dst interface{}) error {
	c, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  dst,
	})
	if err != nil {
		return fmt.Errorf("create decoder fail, err:%w", err)
	}
	if err := c.Decode(src); err != nil {
		return fmt.Errorf("decode type fail, err:%w", err)
	}
	logutil.GetLogger(context.Background()).Debug("decode type finish", zap.Any("src", src), zap.Any("dst", dst))
	return nil
}

func initS3Core(param interface{}) (core.IFsCore, error) {
	s3info := &s3c.S3Config{}
	if err := decodeToType(param, s3info); err != nil {
		return nil, err
	}
	client, err := s3c.New(
		s3c.WithEndpoint(s3info.Endpoint),
		s3c.WithSSL(s3info.UseSSL),
		s3c.WithSecret(s3info.SecretId, s3info.SecretKey),
		s3c.WithBucket(s3info.Bucket),
	)
	if err != nil {
		return nil, fmt.Errorf("init s3 fail, err:%w", err)
	}
	s3core, err := s3.New(
		s3.WithS3Client(client),
	)
	if err != nil {
		return nil, fmt.Errorf("init s3 core fail, err:%w", err)
	}
	return s3core, nil
}

func initMultiTGBotCore(param interface{}) (core.IFsCore, error) {
	botsInfo := []config.BotConfig{}
	if err := decodeToType(param, &botsInfo); err != nil {
		return nil, err
	}
	cores := make([]*bot.TGBot, 0, len(botsInfo))
	for _, botInfo := range botsInfo {
		botcore, err := bot.New(bot.WithAuth(int64(botInfo.Chatid), botInfo.Token))
		if err != nil {
			return nil, fmt.Errorf("init tg bot fail, chatid:%d, token:%s, err:%w", botInfo.Chatid, botInfo.Token, err)
		}
		logutil.GetLogger(context.Background()).Info("init bot succ", zap.Int64("chatid", botcore.GetChatId()),
			zap.String("token", botcore.GetToken()),
			zap.Uint32("bothash", botcore.GetBotHash()))
		cores = append(cores, botcore)
	}
	return bot.NewMultiBot(cores...)
}
