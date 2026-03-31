package main

import (
	"log"

	"WebSupervisor/MyTool/streamtool"
	"WebSupervisor/MyTool/streamtool/models"
)

// TemplateConfig 配置结构体
type TemplateConfig struct {
	Redis struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
	ConsumerGroup string `json:"consumer_group"`
	InputStream   string `json:"input_stream"`
}

type TemplateService struct {
	config *TemplateConfig
	debug  bool
	stream *streamtool.StreamTool
}

func NewTemplateService(config *TemplateConfig, debug bool) *TemplateService {
	// 初始化StreamTool
	streamToolConfig := &models.StreamToolConfig{
		Redis: models.RedisConfig{
			Host:     config.Redis.Host,
			Port:     config.Redis.Port,
			Password: config.Redis.Password,
			DB:       config.Redis.DB,
		},
		Services: []models.ServiceConfig{
			{
				Name:          "echo",
				StreamName:    config.InputStream,
				ConsumerGroup: config.ConsumerGroup,
				ConsumerID:    "template-consumer",
			},
			{
				Name:          "add",
				StreamName:    config.InputStream,
				ConsumerGroup: config.ConsumerGroup,
				ConsumerID:    "template-consumer",
			},
			{
				Name:          "subtract",
				StreamName:    config.InputStream,
				ConsumerGroup: config.ConsumerGroup,
				ConsumerID:    "template-consumer",
			},
			{
				Name:          "multiply",
				StreamName:    config.InputStream,
				ConsumerGroup: config.ConsumerGroup,
				ConsumerID:    "template-consumer",
			},
			{
				Name:          "divide",
				StreamName:    config.InputStream,
				ConsumerGroup: config.ConsumerGroup,
				ConsumerID:    "template-consumer",
			},
		},
		Debug: debug,
	}

	streamtool.InitStreamTool(streamToolConfig)

	return &TemplateService{
		config: config,
		debug:  debug,
		stream: streamtool.GetStreamTool(),
	}
}

// handleStreamMessage 处理流消息
func (s *TemplateService) handleStreamMessage(msg *models.StreamMessage) {
	log.Printf("Processing task: %s", msg.MessageID)

	switch msg.ServiceName {
	case "echo":
		s.HandleEcho(msg)
	case "add":
		s.HandleAdd(msg)
	default:
		if s.debug {
			log.Printf("Debug - Unsupported service: %s", msg.ServiceName)
		}
	}
}
func (s *TemplateService) Start() {
	if s.debug {
		log.Println("Starting Template Service...")
	}

	// 启动消息网关
	s.stream.StartGateway(s.config.InputStream, s.config.ConsumerGroup, "template-service")

	s.stream.StartAllServices(map[string]func(msg *models.StreamMessage){
		"echo":     s.HandleEcho,
		"add":      s.HandleAdd,
		"subtract": s.HandleSubtract,
		"multiply": s.HandleMultiply,
		"divide":   s.HandleDivide,
	})
	if s.debug {
		log.Println("Template Service started successfully")
	}
}

// RegisterHandler 手动注册服务处理器
func (s *TemplateService) RegisterHandler(serviceName string, handler func(msg *models.StreamMessage)) {
	if s.debug {
		log.Printf("Registering handler for service: %s", serviceName)
	}

	s.stream.StartService(serviceName, handler)
}

// SendMessageWithResponse 发送带响应的消息
func (s *TemplateService) SendMessageWithResponse(msg *models.StreamMessage, stream string) *streamtool.Response {
	if s.debug {
		log.Printf("Sending message with response to stream: %s", stream)
	}

	return s.stream.Send(msg, stream)
}

// SendMessageWithoutResponse 发送不带响应的消息
func (s *TemplateService) SendMessageWithoutResponse(msg *models.StreamMessage, stream string) bool {
	if s.debug {
		log.Printf("Sending message without response to stream: %s", stream)
	}

	return s.stream.StreamPush(msg, stream)
}
