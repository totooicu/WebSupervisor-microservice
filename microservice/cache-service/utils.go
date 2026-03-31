package main

import (
	configtool "WebSupervisor/MyTool/config"
)

func loadConfig(configPath string) (*CacheConfig, error) {
	var config CacheConfig
	if err := configtool.LoadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}