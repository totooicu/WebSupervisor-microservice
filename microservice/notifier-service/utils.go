package main

import (
	configtool "WebSupervisor/MyTool/config"
)

func loadConfig(configPath string) (*NotifierConfig, error) {
	var config NotifierConfig
	if err := configtool.LoadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
