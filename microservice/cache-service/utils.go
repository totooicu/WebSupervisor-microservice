package main

import (
	"WebSupervisor/model"
	"WebSupervisor/MyTool"
)

func loadConfig(configPath string) (*model.CacheConfig, error) {
	var config model.CacheConfig
	if err := MyTool.LoadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}