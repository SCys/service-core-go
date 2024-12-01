package core

import (
	"gopkg.in/ini.v1"
)

type Config struct {
	Path string `json:"path"`

	file *ini.File
}

func (cfg *Config) Reload() {
	if cfg.file != nil {
		cfg.file.Reload()
	}
}

func (cfg *Config) GetSection(section string) *ini.Section {
	return cfg.file.Section(section)
}

func (cfg *Config) GetKey(section, key string) *ini.Key {
	return cfg.file.Section(section).Key(key)
}

var globalConfig *Config

func InitGlobalConfig(path string) error {
	var err error

	globalConfig = &Config{Path: path}

	globalConfig.file, err = ini.Load(path)
	if err != nil {
		E("load config error", err)
		return err
	}

	InitLog()

	I("load config success:%s", path)
	return nil
}

func GetGlobalKey(section, key string) *ini.Key {
	if globalConfig != nil {
		return globalConfig.file.Section(section).Key(key)
	}

	return nil
}

func GetGlobalConfig() *Config {
	return globalConfig
}
