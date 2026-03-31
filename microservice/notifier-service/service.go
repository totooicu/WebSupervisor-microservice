package main

import (
	"log"

	"WebSupervisor/MyTool/redis"
	"WebSupervisor/MyTool/streamtool"
	"WebSupervisor/MyTool/streamtool/models"
)

type NotifierService struct {
	redisClient *redis.Client
	config      *NotifierConfig
	debug       bool
}

func NewNotifierService(config *NotifierConfig, debug bool) *NotifierService {
	return &NotifierService{
		redisClient: redis.NewClient(config.Redis.Host, config.Redis.Port, config.Redis.Password, config.Redis.DB),
		config:      config,
		debug:       debug,
	}
}

// Start 启动通知服务
func (s *NotifierService) Start() {
	log.Println("Starting notifier service...")
	if s.debug {
		log.Printf("Debug mode enabled, config: %+v", s.config)
	}

	// 初始化StreamTool
	streamToolConfig := &models.StreamToolConfig{
		Redis: models.RedisConfig{
			Host:     s.config.Redis.Host,
			Port:     s.config.Redis.Port,
			Password: s.config.Redis.Password,
			DB:       s.config.Redis.DB,
		},
		Services: []models.ServiceConfig{
			{
				Name:          "send_email",
				StreamName:    s.config.InputStream,
				ConsumerGroup: s.config.ConsumerGroup,
				ConsumerID:    "notifier-consumer",
			},
		},
		Debug: s.debug,
	}
	
	streamtool.InitStreamTool(streamToolConfig)
	
	// 启动消息网关
	st := streamtool.GetStreamTool()
	st.StartGateway(s.config.InputStream, s.config.ConsumerGroup, "notifier-gateway")
	
	// 启动服务
	st.StartService("send_email", s.handleStreamMessage)
	
	log.Println("Notifier service started successfully")
	
	// 保持服务运行
	select {}
}

// handleStreamMessage 处理流消息
func (s *NotifierService) handleStreamMessage(msg *models.StreamMessage) {
	log.Printf("Processing task: %s", msg.MessageID)
	
	switch msg.ServiceName {
	case "send_email":
		s.handleSendEmail(msg)
	default:
		if s.debug {
			log.Printf("Debug - Unsupported service: %s", msg.ServiceName)
		}
	}
}