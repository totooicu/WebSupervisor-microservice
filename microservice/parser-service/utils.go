package main

import (
	configtool "WebSupervisor/MyTool/config"
)

func loadConfig(configPath string) (*ParserConfig, error) {
	var config ParserConfig
	if err := configtool.LoadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}