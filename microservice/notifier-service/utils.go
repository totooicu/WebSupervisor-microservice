package main

import (
	configtool "WebSupervisor/MyTool/config"
	"WebSupervisor/model"
)

func loadConfig(configPath string) (*model.NotifierConfig, error) {
	var config model.NotifierConfig
	if err := configtool.LoadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
