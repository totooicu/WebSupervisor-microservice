package main

import (
	"log"

	"WebSupervisor/MyTool/streamtool"
	"WebSupervisor/MyTool/streamtool/models"
	"WebSupervisor/model"
)

type TemplateService struct {
	config *model.TemplateConfig
	debug  bool
	stream *streamtool.StreamTool
}

func NewTemplateService(config *model.TemplateConfig, debug bool) *TemplateService {
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
		"echo": s.HandleEcho,
		"add":s.HandleAdd,
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