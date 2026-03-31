package main

import (
	"log"

	"WebSupervisor/MyTool/streamtool"
	"WebSupervisor/MyTool/streamtool/models"
)

type CrawlerService struct {
	config *CrawlerConfig
	debug  bool
}

func NewCrawlerService(config *CrawlerConfig, debug bool) *CrawlerService {
	return &CrawlerService{
		config: config,
		debug:  debug,
	}
}

// Start 启动爬虫服务
func (s *CrawlerService) Start() {
	log.Println("Starting crawler service...")
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
				Name:          "http_request",
				StreamName:    s.config.InputStream,
				ConsumerGroup: s.config.ConsumerGroup,
				ConsumerID:    "crawler-consumer",
			},
		},
		Debug: s.debug,
	}

	streamtool.InitStreamTool(streamToolConfig)

	// 启动消息网关
	st := streamtool.GetStreamTool()
	st.StartGateway(s.config.InputStream, s.config.ConsumerGroup, "crawler-gateway")

	// 启动服务
	st.StartService("http_request", s.handleStreamMessage)

	log.Println("Crawler service started successfully")

	// 保持服务运行
	select {}
}

// handleStreamMessage 处理流消息
func (s *CrawlerService) handleStreamMessage(msg *models.StreamMessage) {
	log.Printf("Processing task: %s", msg.MessageID)

	switch msg.ServiceName {
	case "http_request":
		s.handleHttpRequest(msg)
	default:
		if s.debug {
			log.Printf("Debug - Unsupported service: %s", msg.ServiceName)
		}
	}
}
