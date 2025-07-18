package core

import (
	"sync"

	"gopkg.in/ini.v1"
)

type Config struct {
	Path string `json:"path"`

	file *ini.File
}

func (cfg *Config) Reload() error {
	if cfg.file != nil {
		return cfg.file.Reload()
	}
	return nil
}

func (cfg *Config) GetSection(section string) *ini.Section {
	if cfg.file == nil {
		return nil
	}
	return cfg.file.Section(section)
}

func (cfg *Config) GetKey(section, key string) *ini.Key {
	if cfg.file == nil {
		return nil
	}
	return cfg.file.Section(section).Key(key)
}

// 使用sync.Once确保线程安全的单例模式
var (
	globalConfig *Config
	configOnce   sync.Once
)

func InitGlobalConfig(path string) error {
	var err error

	configOnce.Do(func() {
		globalConfig = &Config{Path: path}

		globalConfig.file, err = ini.Load(path)
		if err != nil {
			E("load config error", err)
			return
		}

		InitLog()

		I("load config success:%s", path)
	})

	return err
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
