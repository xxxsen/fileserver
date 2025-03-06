package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

type Config struct {
	LogInfo    logger.LogConfig  `json:"log_info"`
	DBFile     string            `json:"db_file"`
	ServerInfo ServerConfig      `json:"server_info"`
	BotInfo    BotConfig         `json:"bot_config"`
	AuthInfo   map[string]string `json:"auth_info"`
	S3Bucket   []string          `json:"s3_bucket"`
	TempDir    string            `json:"temp_dir"`
	DebugMode  bool              `json:"debug_mode"`
}

func Parse(f string) (*Config, error) {
	raw, err := os.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("read file:%w", err)
	}
	c := &Config{
		TempDir: filepath.Join(os.TempDir(), "tgfile-temp"),
	}
	if err := json.Unmarshal(raw, c); err != nil {
		return nil, fmt.Errorf("decode json:%w", err)
	}
	return c, nil
}
