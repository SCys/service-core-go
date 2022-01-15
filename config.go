package core

import (
	"gopkg.in/ini.v1"
)

var Config *ini.File

func LoadConfig(filename string) (*ini.File, error) {
	var err error

	if Config == nil {
		Config, err = ini.Load(filename)
		if err != nil {
			E("load config error", err)
			return nil, err
		}

		I("load config success", H{"path": filename})
	} else {
		err := Config.Reload()
		if err != nil {
			W("reload config error", err)
		}
	}

	return Config, nil
}
