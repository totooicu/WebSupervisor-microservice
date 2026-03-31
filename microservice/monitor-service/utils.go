package main

import (
	configtool "WebSupervisor/MyTool/config"
)

func loadConfig(configPath string) (*MonitorConfig, error) {
	var config MonitorConfig
	if err := configtool.LoadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}