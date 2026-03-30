package main

import (
	configtool "WebSupervisor/MyTool/config"
	"WebSupervisor/model"
)

func loadConfig(configPath string) (*model.MonitorConfig, error) {
	var config model.MonitorConfig
	if err := configtool.LoadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}