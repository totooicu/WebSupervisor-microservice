package main

import (
	"log"
	"strconv"

	"WebSupervisor/MyTool/streamtool/models"
)

// ExampleUsage 展示如何使用TemplateService的API
func ExampleUsage(service *TemplateService) {
	log.Println("=== Template Service Usage Examples ===")

	// Example 1: 发送带响应的消息
	log.Println("\n1. Sending message with response (echo service)")
	
	echoMsg := &models.StreamMessage{
		ServiceName:    "echo",
		CallbackStream: "response-stream",
		Playload: map[string]interface{}{
			"message": "Hello, Stream!",
		},
	}
	
	// 使用SendMessageWithResponse发送带响应的消息
	response := service.SendMessageWithResponse(echoMsg, "template-service-input-stream")
	if response != nil {
		// 使用response.Get()阻塞获取响应
		result := response.Get()
		log.Printf("Echo response: %v", result)
	}

	// Example 2: 发送带响应的消息（加法服务）
	log.Println("\n2. Sending message with response (add service)")
	
	addMsg := &models.StreamMessage{
		ServiceName:    "add",
		CallbackStream: "response-stream",
		Playload: map[string]interface{}{
			"num1": 10.5,
			"num2": 20.5,
		},
	}
	
	response = service.SendMessageWithResponse(addMsg, "template-service-input-stream")
	if response != nil {
		result := response.Get()
		log.Printf("Add response: %v", result)
	}

	// Example 3: 发送不带响应的消息
	log.Println("\n3. Sending message without response")
	
	notificationMsg := &models.StreamMessage{
		ServiceName: "notification",
		Playload: map[string]interface{}{
			"type":    "info",
			"message": "This is a notification without response",
		},
	}
	
	// 使用SendMessageWithoutResponse发送不带响应的消息
	success := service.SendMessageWithoutResponse(notificationMsg, "notification-stream")
	if success {
		log.Println("Notification sent successfully")
	} else {
		log.Println("Failed to send notification")
	}
}

// Example of how to register custom handlers
func ExampleCustomHandlers(service *TemplateService) {
	log.Println("\n=== Custom Handler Registration Example ===")
	
	// 注册自定义处理器
	service.RegisterHandler("custom-service", func(msg *models.StreamMessage) {
		log.Printf("Custom service received message: %s", msg.MessageID)
		
		// 处理业务逻辑
		// ...
		
		// 发送响应
		responseMsg := &models.StreamMessage{
			MessageID:      strconv.Itoa(service.stream.GetMessageID()),
			ReplyID:        msg.MessageID,
			ServiceName:    "response",
			CallbackStream: msg.CallbackStream,
			Playload: map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"message": "Custom service processed successfully",
				},
			},
		}
		
		service.stream.StreamPush(responseMsg, msg.CallbackStream)
	})
	
	log.Println("Custom handler registered successfully")
}