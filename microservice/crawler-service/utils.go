package main

import (
	"WebSupervisor/model"
	"WebSupervisor/MyTool"
)

func loadConfig(configPath string) (*model.CrawlerConfig, error) {
	var config model.CrawlerConfig
	if err := MyTool.LoadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}