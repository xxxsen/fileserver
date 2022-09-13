package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/xxxsen/common/database"
	"github.com/xxxsen/common/logger"
	"github.com/xxxsen/common/s3"
)

type ServerConfig struct {
	Address string `json:"address"`
}

type IDGenConfig struct {
	WorkerID uint16 `json:"worker_id"`
}

type Config struct {
	LogInfo    logger.LogConfig  `json:"log_info"`
	FileDBInfo database.DBConfig `json:"file_db_info"`
	ServerInfo ServerConfig      `json:"server_info"`
	S3Info     s3.S3Config       `json:"s3_info"`
	IDGenInfo  IDGenConfig       `json:"idgen_info"`
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
