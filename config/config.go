package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/xxxsen/common/database"
	"github.com/xxxsen/common/logger"
)

type ServerConfig struct {
	Address string `json:"address"`
}

type IDGenConfig struct {
	WorkerID uint16 `json:"worker_id"`
}

type BotConfig struct {
	Chatid uint64 `json:"chatid"`
	Token  string `json:"token"`
}

type IOConfig struct {
	MaxUploadThread   int `json:"max_upload_thread"`
	MaxDownloadThread int `json:"max_download_thread"`
}

type FakeS3Config struct {
	Enable     bool     `json:"enable"`
	BucketList []string `json:"bucket_list"`
}

type RefererConfig struct {
	Enable  bool     `json:"enable"`
	Referer []string `json:"referer"`
}

type Config struct {
	LogInfo     logger.LogConfig       `json:"log_info"`
	FileDBInfo  database.DBConfig      `json:"file_db_info"`
	ServerInfo  ServerConfig           `json:"server_info"`
	IDGenInfo   IDGenConfig            `json:"idgen_info"`
	FsInfo      map[string]interface{} `json:"fs_info"`
	UploadFs    string                 `json:"upload_fs"`
	AuthInfo    map[string]string      `json:"auth_info"`
	IOInfo      IOConfig               `json:"io_info"`
	FakeS3Info  FakeS3Config           `json:"fake_s3_info"`
	RefererInfo RefererConfig          `json:"referer_info"`
}

func Parse(f string) (*Config, error) {
	raw, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("read file:%w", err)
	}
	c := &Config{}
	if err := json.Unmarshal(raw, c); err != nil {
		return nil, fmt.Errorf("decode json:%w", err)
	}
	return c, nil
}
