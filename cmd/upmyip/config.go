package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	LambdaName string `toml:"lambda"`
	AccessKey  string `toml:"access_key"`
	SecretKey  string `toml:"secret_key"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}
	return cfg, nil
}
