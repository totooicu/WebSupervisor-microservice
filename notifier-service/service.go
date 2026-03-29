package main

import (
	"log"

	"WebSupervisor/new/model"
	"WebSupervisor/new/redis"
	streams_communication "streams-communication"
	streams_model "streams-communication/model"
)

type NotifierService struct {
	redisClient *redis.Client
	config      *model.NotifierConfig
	debug       bool
}

func NewNotifierService(config *model.NotifierConfig, debug bool) *NotifierService {
	return &NotifierService{
		redisClient: redis.NewClient(config.Redis.Host, config.Redis.Port, config.Redis.Password, config.Redis.DB),
		config:      config,
		debug:       debug,
	}
}

func (s *NotifierService) Start() {
	log.Println("Starting notifier service...")
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
		"notifier-consumer",
		"send_email",
	)
	
	// 启动流消费者
	consumer.Start()
}

// ProcessTask 实现StreamProcessor接口
func (s *NotifierService) ProcessTask(task streams_model.Message) {
	log.Printf("Processing task: %s", task.TaskID)
	
	switch task.ServiceName {
	case "send_email":
		s.handleSendEmail(task)
	default:
		if s.debug {
			log.Printf("Debug - Unsupported service: %s", task.ServiceName)
		}
	}
}