package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Directory DirectoryConfig
	Volume    VolumeConfig
	Mail      MailConfig
	Rank      RankConfig
}

type DirectoryConfig struct {
	ROOT_DIRECTORY     string `toml:"ROOT_DIRECTORY"`
	TARGET_DIRECTORIES string `toml:"TARGET_DIRECTORIES"`
	IGNORE_DIRECTORIES string `toml:"IGNORE_DIRECTORIES"`
}

type VolumeConfig struct {
	FREE_BYTE_TH int `toml:"FREE_BYTE_TH"`
}

type MailConfig struct {
	USE_SEND_MAIL bool   `toml:"USE_SEND_MAIL"`
	USER_NAME     string `toml:"USER_NAME"`
	PASSWORD      string `toml:"PASSWORD"`
	TO            string `toml:"TO"`
	CC            string `toml:"CC"`
}

type RankConfig struct {
	MAX int `toml:"MAX"`
}

// 読み込み
func LoadConfig(path string) (*Config, error) {
	c := Config{}
	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
