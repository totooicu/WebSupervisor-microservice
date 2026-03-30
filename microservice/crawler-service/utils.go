package main

import (
	configtool "WebSupervisor/MyTool/config"
	"WebSupervisor/model"
)

func loadConfig(configPath string) (*model.CrawlerConfig, error) {
	var config model.CrawlerConfig
	if err := configtool.LoadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}