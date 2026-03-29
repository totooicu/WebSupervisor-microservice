package main

import (
	"encoding/json"
	"io/ioutil"

	"WebSupervisor/model"
	"WebSupervisor/MyTool"
)

func loadConfig(configPath string) (*model.MonitorConfig, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// 先解析为map，应用环境变量替换
	var configMap map[string]interface{}
	if err := json.Unmarshal(data, &configMap); err != nil {
		return nil, err
	}

	// 替换环境变量
	MyTool.ExpandEnvVarsInMap(configMap)

	// 将map转换回JSON
	processedData, err := json.Marshal(configMap)
	if err != nil {
		return nil, err
	}

	// 解析为配置结构体
	var config model.MonitorConfig
	if err := json.Unmarshal(processedData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}