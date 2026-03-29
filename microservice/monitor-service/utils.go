package main

import (
	"WebSupervisor/model"
	"WebSupervisor/MyTool"
)

func loadConfig(configPath string) (*model.MonitorConfig, error) {
	var config model.MonitorConfig
	if err := MyTool.LoadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}