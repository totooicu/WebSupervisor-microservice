package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"WebSupervisor/model"
	"WebSupervisor/MyTool"
)

func mapToStruct(m map[string]interface{}, s interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, s)
}

func loadConfig(configPath string) (*model.NotifierConfig, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// 先解析为map，应用环境变量替换
	var configMap map[string]interface{}
	if err := json.Unmarshal(data, &configMap); err != nil {
		return nil, err
	}

	log.Printf("Debug - Before env expansion: %+v", configMap)

	// 替换环境变量
	MyTool.ExpandEnvVarsInMap(configMap)

	log.Printf("Debug - After env expansion: %+v", configMap)

	// 将map转换回JSON
	processedData, err := json.Marshal(configMap)
	if err != nil {
		return nil, err
	}

	log.Printf("Debug - Processed JSON: %s", string(processedData))

	// 解析为配置结构体
	var config model.NotifierConfig
	if err := json.Unmarshal(processedData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
