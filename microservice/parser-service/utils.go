package main

import (
	"WebSupervisor/model"
	"WebSupervisor/MyTool"
)

func loadConfig(configPath string) (*model.ParserConfig, error) {
	var config model.ParserConfig
	if err := MyTool.LoadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}