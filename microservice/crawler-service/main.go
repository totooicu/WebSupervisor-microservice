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

	service := NewCrawlerService(config, *debug)
	service.Start()
}


/*
{
  "callback_stream": "notifier-tasks",
  "consumer_group": "notifier-group",
  "message_id": "1",
  "playload": "{\"url\": \"https://docs.mai-mai.org/manual/\",   \"method\": \"GET\",   \"headers\": {},   \"body\": {},\"str_payload\": \"\"}",
  "reply_id": "-1",
  "service_name": "http_request"
}

"{\"url\": \"https://docs.mai-mai.org/manual/\",   \"method\": \"GET\",   \"headers\": {},   \"body\": {},\"str_payload\": \"\"}"

*/