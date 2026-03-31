package main

import (
	"fmt"
	"log"

	"encoding/json"

	"WebSupervisor/MyTool/streamtool"
	"WebSupervisor/MyTool/streamtool/models"
)

// ServiceCommunicator 服务通信器
type ServiceCommunicator struct {
	crawlerStream  string
	parserStream   string
	cacheStream    string
	notifierStream string
	inputStream    string
	debug          bool
}

// NewServiceCommunicator 创建服务通信器
func NewServiceCommunicator(config *MonitorConfig, debug bool) *ServiceCommunicator {
	return &ServiceCommunicator{
		crawlerStream:  config.CrawlerServiceInputStream,
		parserStream:   config.ParserServiceInputStream,
		cacheStream:    config.CacheServiceInputStream,
		notifierStream: config.NotifierServiceInputStream,
		inputStream:    config.InputStream,
		debug:          debug,
	}
}

// SendToCrawlerService 发送任务到爬虫服务
func (sc *ServiceCommunicator) SendToCrawlerService(urlConfig URLConfig) error {
	paramData := map[string]interface{}{
		"url":         urlConfig.URL,
		"method":      urlConfig.Method,
		"headers":     urlConfig.Header,
		"body":        urlConfig.Body,
		"str_payload": urlConfig.StringPlayLoad,
		"output":      urlConfig.Output,
		"test":        urlConfig.Test,
	}

	streamMsg := &models.StreamMessage{
		ServiceName:    "http_request",
		CallbackStream: sc.inputStream,
		Playload:       paramData,
	}

	st := streamtool.GetStreamTool()
	if !st.StreamPush(streamMsg, sc.crawlerStream) {
		return fmt.Errorf("failed to publish to crawler stream")
	}

	if sc.debug {
		log.Printf("Debug - Sent task to crawler service: %+v", streamMsg)
	}

	return nil
}

// SendToParserService 发送任务到解析服务
func (sc *ServiceCommunicator) SendToParserService(rawData string, urlConfig URLConfig) error {
	paramData := map[string]interface{}{
		"raw_data":  rawData,
		"type":      urlConfig.Type,
		"json_keys": urlConfig.JSONKeys,
		"html_keys": urlConfig.HTMLKeys,
		"output":    urlConfig.Output,
	}

	streamMsg := &models.StreamMessage{
		ServiceName:    "parse_data",
		CallbackStream: sc.inputStream,
		Playload:       paramData,
	}

	st := streamtool.GetStreamTool()
	if !st.StreamPush(streamMsg, sc.parserStream) {
		return fmt.Errorf("failed to publish to parser stream")
	}

	if sc.debug {
		log.Printf("Debug - Sent task to parser service: %+v", streamMsg)
	}

	return nil
}

// SendToCacheService 发送任务到缓存服务
func (sc *ServiceCommunicator) SendToCacheService(parsedData interface{}, urlConfig URLConfig) error {
	paramData := map[string]interface{}{
		"app":    "monitor-service",
		"key":    urlConfig.URL,
		"data":   parsedData,
		"expire": 3600, // 1小时过期
	}

	streamMsg := &models.StreamMessage{
		ServiceName:    "cache_data",
		CallbackStream: sc.inputStream,
		Playload:       paramData,
	}

	st := streamtool.GetStreamTool()
	if !st.StreamPush(streamMsg, sc.cacheStream) {
		return fmt.Errorf("failed to publish to cache stream")
	}

	if sc.debug {
		log.Printf("Debug - Sent task to cache service: %+v", streamMsg)
	}

	return nil
}

// SendNotification 发送通知到通知服务
func (sc *ServiceCommunicator) SendNotification(subject string, content string) error {
	paramData := map[string]interface{}{
		"subject": subject,
		"content": content,
	}

	streamMsg := &models.StreamMessage{
		ServiceName:    "send_email",
		CallbackStream: "",
		Playload:       paramData,
	}

	st := streamtool.GetStreamTool()
	if !st.StreamPush(streamMsg, sc.notifierStream) {
		return fmt.Errorf("failed to publish to notifier stream")
	}

	if sc.debug {
		log.Printf("Debug - Sent notification: %+v", streamMsg)
	}

	return nil
}

// SendMessageWithResponse 发送消息并等待响应
func (sc *ServiceCommunicator) SendMessageWithResponse(stream string, task Message) (*streamtool.Response, error) {
	// 转换为StreamMessage
	playload := map[string]interface{}{}
	if err := json.Unmarshal([]byte(task.Playload), &playload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal playload: %w", err)
	}
	streamMsg := &models.StreamMessage{
		ServiceName:    task.ServiceName,
		CallbackStream: sc.inputStream,
		Playload:       playload,
	}

	st := streamtool.GetStreamTool()
	response := st.Send(streamMsg, stream)
	if response == nil {
		return nil, fmt.Errorf("failed to send message")
	}

	if sc.debug {
		log.Printf("Debug - Message sent, response object created")
	}

	return response, nil
}

// StartEventListener 启动事件监听器（使用streamtool后不再需要）
func (sc *ServiceCommunicator) StartEventListener() {
	// streamtool已经处理了事件监听，此方法保留以保持兼容性
	log.Println("Event listener already handled by streamtool")
}
