package main

import (
	"encoding/json"
	"log"

	"WebSupervisor/model"
	"WebSupervisor/MyTool"
	streams_model "streams-communication/model"
)

func (s *NotifierService) handleSendEmail(task streams_model.Message) {
	var params model.NotifierParameter
	if err := json.Unmarshal([]byte(task.Playload), &params); err != nil {
		log.Printf("Error unmarshalling playload: %v", err)
		return
	}

	userName := s.config.Mail.UserName
	password := s.config.Mail.Password
	tos := s.config.Mail.Tos

	if params.UserName != "" {
		userName = params.UserName
	}

	if params.Password != "" {
		password = params.Password
	}

	if len(params.Tos) > 0 {
		tos = params.Tos
	}

	// 校验发送者邮箱格式
	if !isValidEmail(userName) {
		log.Printf("Error: Invalid sender email format: %s", userName)
		return
	}

	// 校验接收者邮箱格式
	for _, to := range tos {
		if !isValidEmail(to) {
			log.Printf("Error: Invalid recipient email format: %s", to)
			return
		}
	}

	if s.debug {
		log.Printf("Debug - Email parameters: userName=%s, tos=%v, subject=%s", userName, tos, params.Subject)
	}

	log.Printf("Sending email to: %v", tos)

	if err := MyTool.SendEmail(userName, password, tos, params.Subject, params.Content); err != nil {
		log.Printf("Error sending email: %v", err)
		return
	}

	log.Printf("Email sent successfully")

	if s.debug {
		log.Printf("Debug - Email sent successfully to: %v", tos)
	}

	paramData := map[string]interface{}{
		"success": true,
		"tos":     tos,
	}
	paramBytes, _ := json.Marshal(paramData)

	// 构建新的响应消息格式
	result := map[string]interface{}{
		"callback_stream": s.config.InputStream, // 本微服务的streams
		"consumer_group":  task.CallbackStream,  // 请求消息的callback_stream
		"playload":        string(paramBytes),   // 处理结果
		"service_name":    "response",           // 响应消息统一为response
		"message_id":      task.ReplyID,         // 上一次消息编号+1
		"reply_id":        task.TaskID,          // 请求消息的id
	}

	if err := s.redisClient.PublishMessage(task.CallbackStream, result); err != nil {
		log.Printf("Error publishing result: %v", err)
	} else if s.debug {
		log.Printf("Debug - Published email result to stream: %s", task.CallbackStream)
	}
}

// isValidEmail 校验邮箱格式
func isValidEmail(email string) bool {
	// 简单的邮箱格式校验
	if email == "" {
		return false
	}

	// 检查是否包含@符号
	atIndex := -1
	for i, c := range email {
		if c == '@' {
			atIndex = i
			break
		}
	}

	if atIndex == -1 || atIndex == 0 || atIndex == len(email)-1 {
		return false
	}

	// 检查@后面是否有域名部分
	domain := email[atIndex+1:]
	if domain == "" {
		return false
	}

	// 检查域名是否包含点号
	hasDot := false
	for _, c := range domain {
		if c == '.' {
			hasDot = true
			break
		}
	}

	return hasDot
}
