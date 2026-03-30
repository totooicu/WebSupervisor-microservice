package main

import (
	configtool "WebSupervisor/MyTool/config"
	"WebSupervisor/model"
)

func loadConfig(configPath string) (*model.ParserConfig, error) {
	var config model.ParserConfig
	if err := configtool.LoadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}