package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"WebSupervisor/model"
	"WebSupervisor/MyTool/redis"
	"WebSupervisor/MyTool/streams"
)

// ServiceCommunicator 服务通信器
type ServiceCommunicator struct {
	redisClient         *redis.Client
	crawlerStream       string
	parserStream        string
	cacheStream         string
	notifierStream      string
	debug               bool
}

// NewServiceCommunicator 创建服务通信器
func NewServiceCommunicator(redisClient *redis.Client, config *model.MonitorConfig, debug bool) *ServiceCommunicator {
	return &ServiceCommunicator{
		redisClient:         redisClient,
		crawlerStream:       config.CrawlerServiceInputStream,
		parserStream:        config.ParserServiceInputStream,
		cacheStream:         config.CacheServiceInputStream,
		notifierStream:      config.NotifierServiceInputStream,
		debug:               debug,
	}
}

// SendToCrawlerService 发送任务到爬虫服务
func (sc *ServiceCommunicator) SendToCrawlerService(urlConfig URLConfig) error {
	taskID := generateTaskID()
	
	paramData := map[string]interface{}{
		"url":         urlConfig.URL,
		"method":      urlConfig.Method,
		"headers":     urlConfig.Header,
		"body":        urlConfig.Body,
		"str_payload": urlConfig.StringPlayLoad,
		"output":      urlConfig.Output,
		"test":        urlConfig.Test,
	}
	
	paramBytes, err := json.Marshal(paramData)
	if err != nil {
		return fmt.Errorf("failed to marshal crawler params: %w", err)
	}
	
	task := model.Message{
		TaskID:         taskID,
		ConsumerGroup:  "crawler-group",
		CallbackStream: sc.parserStream,
		ServiceName:    "http_request",
		Playload:       string(paramBytes),
	}
	
	if sc.debug {
		log.Printf("Debug - Sending task to crawler service: %+v", task)
	}
	
	if err := sc.redisClient.PublishMessage(sc.crawlerStream, task); err != nil {
		return fmt.Errorf("failed to publish to crawler stream: %w", err)
	}
	
	log.Printf("Successfully sent task %s to crawler service", taskID)
	return nil
}

// SendToParserService 发送任务到解析服务
func (sc *ServiceCommunicator) SendToParserService(rawData string, urlConfig URLConfig) error {
	taskID := generateTaskID()
	
	paramData := map[string]interface{}{
		"raw_data":    rawData,
		"type":        urlConfig.Type,
		"json_keys":   urlConfig.JSONKeys,
		"html_keys":   urlConfig.HTMLKeys,
		"output":      urlConfig.Output,
	}
	
	paramBytes, err := json.Marshal(paramData)
	if err != nil {
		return fmt.Errorf("failed to marshal parser params: %w", err)
	}
	
	task := model.Message{
		TaskID:         taskID,
		ConsumerGroup:  "parser-group",
		CallbackStream: sc.cacheStream,
		ServiceName:    "parse_data",
		Playload:       string(paramBytes),
	}
	
	if sc.debug {
		log.Printf("Debug - Sending task to parser service: %+v", task)
	}
	
	if err := sc.redisClient.PublishMessage(sc.parserStream, task); err != nil {
		return fmt.Errorf("failed to publish to parser stream: %w", err)
	}
	
	log.Printf("Successfully sent task %s to parser service", taskID)
	return nil
}

// SendToCacheService 发送任务到缓存服务
func (sc *ServiceCommunicator) SendToCacheService(parsedData interface{}, urlConfig URLConfig) error {
	taskID := generateTaskID()
	
	paramData := map[string]interface{}{
		"app":          "monitor-service",
		"key":          urlConfig.URL,
		"data":         parsedData,
		"expire":       3600, // 1小时过期
	}
	
	paramBytes, err := json.Marshal(paramData)
	if err != nil {
		return fmt.Errorf("failed to marshal cache params: %w", err)
	}
	
	task := model.Message{
		TaskID:         taskID,
		ConsumerGroup:  "cache-group",
		CallbackStream: sc.notifierStream,
		ServiceName:    "cache_data",
		Playload:       string(paramBytes),
	}
	
	if sc.debug {
		log.Printf("Debug - Sending task to cache service: %+v", task)
	}
	
	if err := sc.redisClient.PublishMessage(sc.cacheStream, task); err != nil {
		return fmt.Errorf("failed to publish to cache stream: %w", err)
	}
	
	log.Printf("Successfully sent task %s to cache service", taskID)
	return nil
}

// SendNotification 发送通知到通知服务
func (sc *ServiceCommunicator) SendNotification(notificationType, message string) error {
	taskID := generateTaskID()
	
	paramData := map[string]interface{}{
		"type":    notificationType,
		"message": message,
		"time":    time.Now().Format(time.RFC3339),
	}
	
	paramBytes, err := json.Marshal(paramData)
	if err != nil {
		return fmt.Errorf("failed to marshal notification params: %w", err)
	}
	
	task := model.Message{
		TaskID:         taskID,
		ConsumerGroup:  "notifier-group",
		CallbackStream: "",
		ServiceName:    "send_notification",
		Playload:       string(paramBytes),
	}
	
	if sc.debug {
		log.Printf("Debug - Sending notification: %+v", task)
	}
	
	if err := sc.redisClient.PublishMessage(sc.notifierStream, task); err != nil {
		return fmt.Errorf("failed to publish to notifier stream: %w", err)
	}
	
	log.Printf("Successfully sent notification task %s", taskID)
	return nil
}

// StartEventListener 启动事件监听器
func (sc *ServiceCommunicator) StartEventListener() {
	go func() {
		log.Println("Starting event listener...")
		
		// 监听缓存服务的回调
		cacheConsumer := streams.NewStreamConsumer(
			sc.redisClient,
			nil,
			&EventListener{communicator: sc},
			sc.debug,
			sc.cacheStream,
			"monitor-listener",
			"monitor-listener-1",
			"cache_callback",
		)
		
		cacheConsumer.Start()
	}()
}

// EventListener 事件监听器
type EventListener struct {
	communicator *ServiceCommunicator
}

// ProcessTask 处理事件任务
func (el *EventListener) ProcessTask(task model.Message) {
	if el.communicator.debug {
		log.Printf("Debug - Received event task: %+v", task)
	}
	
	switch task.ServiceName {
	case "cache_callback":
		el.handleCacheCallback(task)
	case "parser_callback":
		el.handleParserCallback(task)
	case "crawler_callback":
		el.handleCrawlerCallback(task)
	default:
		log.Printf("Unknown event type: %s", task.ServiceName)
	}
}

// handleCacheCallback 处理缓存回调
func (el *EventListener) handleCacheCallback(task model.Message) {
	var callbackData struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	
	if err := json.Unmarshal([]byte(task.Playload), &callbackData); err != nil {
		log.Printf("Error unmarshalling cache callback: %v", err)
		return
	}
	
	if callbackData.Status == "success" {
		log.Printf("Cache operation successful: %s", callbackData.Message)
		
		// 发送通知
		notificationMsg := fmt.Sprintf("Data cached successfully: %s", callbackData.Message)
		if err := el.communicator.SendNotification("info", notificationMsg); err != nil {
			log.Printf("Failed to send notification: %v", err)
		}
	} else {
		log.Printf("Cache operation failed: %s", callbackData.Message)
		
		// 发送错误通知
		errorMsg := fmt.Sprintf("Cache operation failed: %s", callbackData.Message)
		if err := el.communicator.SendNotification("error", errorMsg); err != nil {
			log.Printf("Failed to send error notification: %v", err)
		}
	}
}

// handleParserCallback 处理解析回调
func (el *EventListener) handleParserCallback(task model.Message) {
	log.Printf("Received parser callback: %s", task.TaskID)
	// 可以根据需要实现解析回调处理逻辑
}

// handleCrawlerCallback 处理爬虫回调
func (el *EventListener) handleCrawlerCallback(task model.Message) {
	log.Printf("Received crawler callback: %s", task.TaskID)
	// 可以根据需要实现爬虫回调处理逻辑
}