package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

type LogConfig struct {
	File  string
	Level string
}

type ServerConfig struct {
	Host string
	Port string
}

type DBConfig struct {
	ConnectionString string `toml:"connection_string"`
}

type Config struct {
	Logger      LogConfig
	Server      ServerConfig
	Database    DBConfig
	StorageType string `toml:"storage_type"`
}

func NewConfig(configPath string) (c Config, err error) {
	config, err := os.ReadFile(configPath)
	if err != nil {
		return c, err
	}

	_, err = toml.Decode(string(config), &c)
	return c, err
}
