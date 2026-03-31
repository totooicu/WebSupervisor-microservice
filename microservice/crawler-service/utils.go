package main

import (
	configtool "WebSupervisor/MyTool/config"
)

func loadConfig(configPath string) (*CrawlerConfig, error) {
	var config CrawlerConfig
	if err := configtool.LoadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}