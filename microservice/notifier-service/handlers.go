package main

import (
	"log"
	"strconv"

	email "WebSupervisor/MyTool/email"
	"WebSupervisor/MyTool/streamtool"
	"WebSupervisor/MyTool/streamtool/models"
)

func (s *NotifierService) handleSendEmail(msg *models.StreamMessage) {
	// 解析参数
	var params NotifierParameter
	email_content := msg.Playload
	params.Subject = email_content["subject"].(string)
	params.Content = email_content["content"].(string)
	if email_content["userName"] != nil {
		params.UserName = email_content["userName"].(string)
	}
	if email_content["password"] != nil {
		params.Password = email_content["password"].(string)
	}
	if email_content["tos"] != nil {
		params.Tos = email_content["tos"].([]string)
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

	if err := email.SendEmail(userName, password, tos, params.Subject, params.Content); err != nil {
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

	// 构造响应消息
	st := streamtool.GetStreamTool()
	responseMsg := &models.StreamMessage{
		MessageID:      strconv.Itoa(st.GetMessageID()),
		ReplyID:        msg.MessageID,
		ServiceName:    "response",
		CallbackStream: msg.CallbackStream,
		Playload:       paramData,
	}

	if !st.StreamPush(responseMsg, msg.CallbackStream) {
		log.Printf("Error publishing result: failed to push to stream")
	} else if s.debug {
		log.Printf("Debug - Published email result to stream: %s", msg.CallbackStream)
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
