package main

import (
	"log"

	"WebSupervisor/MyTool/redis"
	"WebSupervisor/MyTool/streamtool"
	"WebSupervisor/MyTool/streamtool/models"
	"WebSupervisor/model"
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

// Start 启动解析服务
func (s *ParserService) Start() {
	log.Println("Starting parser service...")
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
				Name:          "parse_html",
				StreamName:    s.config.InputStream,
				ConsumerGroup: s.config.ConsumerGroup,
				ConsumerID:    "parser-consumer",
			},
			{
				Name:          "parse_json",
				StreamName:    s.config.InputStream,
				ConsumerGroup: s.config.ConsumerGroup,
				ConsumerID:    "parser-consumer",
			},
		},
		Debug: s.debug,
	}

	streamtool.InitStreamTool(streamToolConfig)

	// 启动消息网关
	st := streamtool.GetStreamTool()
	st.StartGateway(s.config.InputStream, s.config.ConsumerGroup, "parser-gateway")

	// 启动服务
	st.StartService("parse_html", s.handleStreamMessage)
	st.StartService("parse_json", s.handleStreamMessage)

	log.Println("Parser service started successfully")

	// 保持服务运行
	select {}
}

// handleStreamMessage 处理流消息
func (s *ParserService) handleStreamMessage(msg *models.StreamMessage) {
	log.Printf("Processing task: %s", msg.MessageID)

	switch msg.ServiceName {
	case "parse_html":
		s.handleParseHTML(msg)
	case "parse_json":
		s.handleParseJSON(msg)
	default:
		if s.debug {
			log.Printf("Debug - Unsupported service: %s", msg.ServiceName)
		}
	}
}
