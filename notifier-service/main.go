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
	log.Println(">>> config",config)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	service := NewNotifierService(config, *debug)
	service.Start()
}

/*

{
   "callback_stream": "notifier-tasks", 
   "consumer_group": "notifier-group",
   "playload": "{ \"subject\": \"邮件主题\",\"content\": \"邮件内容\"}",
   "service_name": "send_email",  
   "message_id": "1",
   "reply_id": "0"
}


"{ \"subject\": \"邮件主题\",\"content\": \"邮件内容\"}"

*/