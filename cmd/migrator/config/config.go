package config

import (
	"encoding/json"
	"os"
)

type Source struct {
	LinkList []string `json:"link_list"`
}

type Destination struct {
	Host      string `json:"host"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type Config struct {
	Src Source      `json:"src"`
	Dst Destination `json:"dst"`
}

func Parse(f string) (*Config, error) {
	c := &Config{}
	raw, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, c); err != nil {
		return nil, err
	}
	return c, nil
}
