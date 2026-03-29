package main

import (
	"flag"
	"log"
)

func main() {
	configPath := flag.String("config", "./config.json", "Path to config file")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	if *debug {
		log.Println("Debug mode enabled")
	}

	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	service := NewCacheService(config, *debug)
	service.Start()
}

/*
{
"callback_stream": "notifier-tasks",
"consumer_group": "notifier-group",
"service_name":"send_email",
"parameter": "{\"app\": \"app_name\",\"key\": \"key\",\"data\": \"要保存的数据\"}"
}
*/