package main

import (
	"log"

	"WebSupervisor/model"
	"WebSupervisor/MyTool/redis"
	streams_communication "streams-communication"
	streams_model "streams-communication/model"
)

type CrawlerService struct {
	redisClient *redis.Client
	config      *model.CrawlerConfig
	debug       bool
}

func NewCrawlerService(config *model.CrawlerConfig, debug bool) *CrawlerService {
	return &CrawlerService{
		redisClient: redis.NewClient(config.Redis.Host, config.Redis.Port, config.Redis.Password, config.Redis.DB),
		config:      config,
		debug:       debug,
	}
}

// ProcessTask 实现StreamProcessor接口
func (s *CrawlerService) ProcessTask(task streams_model.Message) {
	log.Printf("Processing task: %s", task.TaskID)
	
	switch task.ServiceName {
	case "http_request":
		s.handleHttpRequest(task)
	default:
		if s.debug {
			log.Printf("Debug - Unsupported service: %s", task.ServiceName)
		}
	}
}

func (s *CrawlerService) Start() {
	log.Println("Starting crawler service...")
	if s.debug {
		log.Printf("Debug mode enabled, config: %+v", s.config)
	}

	// 创建流消费者
	consumer := streams_communication.NewStreamConsumer(
		s.redisClient,
		s.config,
		s,
		s.debug,
		s.config.InputStream,
		s.config.ConsumerGroup,
		"crawler-consumer",
		"http_request",
	)
	
	// 启动流消费者
	consumer.Start()
}

