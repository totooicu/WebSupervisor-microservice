package main

import (
	"flag"
	"log"

	configtool "WebSupervisor/MyTool/config"
	"WebSupervisor/model"
)

func main() {
	configPath := flag.String("config", "./config.json", "Path to config file")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	if *debug {
		log.Println("Debug mode enabled")
	}

	var config model.TemplateConfig
	if err := configtool.LoadConfig(*configPath, &config); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	service := NewTemplateService(&config, *debug)
	
	// 手动注册服务处理器
	service.RegisterHandler("echo", service.HandleEcho)
	service.RegisterHandler("add", service.HandleAdd)
	
	// 启动服务
	service.Start()
	
	// 服务启动后会持续运行
	select {}
}