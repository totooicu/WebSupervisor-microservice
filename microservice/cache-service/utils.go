package main

import (
	configtool "WebSupervisor/MyTool/config"
	"WebSupervisor/model"
)

func loadConfig(configPath string) (*model.CacheConfig, error) {
	var config model.CacheConfig
	if err := configtool.LoadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}