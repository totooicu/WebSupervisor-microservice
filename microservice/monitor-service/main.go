package main

import (
	"flag"
	"log"

	configtool "WebSupervisor/MyTool/config"
)

func main() {
	configPath := flag.String("config", "./config.json", "Path to config file")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	if *debug {
		log.Println("Debug mode enabled")
	}

	var config MonitorConfig
	if err := configtool.LoadConfig(*configPath, &config); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	service := NewMonitorService(&config, *debug)
	service.Start()
}
