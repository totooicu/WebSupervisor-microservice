package main

import (
	"WebSupervisor/model"
	"WebSupervisor/MyTool"
)

func loadConfig(configPath string) (*model.NotifierConfig, error) {
	var config model.NotifierConfig
	if err := MyTool.LoadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
