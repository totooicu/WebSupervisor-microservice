package main

import (
	"log"

	"WebSupervisor/new/model"
	"WebSupervisor/new/redis"
	streams_communication "streams-communication"
	streams_model "streams-communication/model"
)

type ParserService struct {
	redisClient *redis.Client
	config      *model.ParserConfig
	debug       bool
}

func NewParserService(config *model.ParserConfig, debug bool) *ParserService {
	return &ParserService{
		redisClient: redis.NewClient(config.Redis.Host, config.Redis.Port, config.Redis.Password, config.Redis.DB),
		config:      config,
		debug:       debug,
	}
}

func (s *ParserService) Start() {
	log.Println("Starting parser service...")
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
		"parser-consumer",
		"",
	)
	
	// 启动流消费者
	consumer.Start()
}

// ProcessTask 实现StreamProcessor接口
func (s *ParserService) ProcessTask(task streams_model.Message) {
	log.Printf("Processing task: %s", task.TaskID)
	
	switch task.ServiceName {
	case "parse_html":
		s.handleParseHTML(task)
	case "parse_json":
		s.handleParseJSON(task)
	default:
		if s.debug {
			log.Printf("Debug - Unsupported service: %s", task.ServiceName)
		}
	}
}