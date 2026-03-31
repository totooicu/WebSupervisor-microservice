package main

import (
	"encoding/json"
	"log"
	"strconv"

	"WebSupervisor/MyTool/streamtool/models"
)

// Example handlers for demonstration
func handleStreamMessage(){

	
}
// HandleEcho 处理echo服务
func (s *TemplateService) HandleEcho(msg *models.StreamMessage) {
	if s.debug {
		log.Printf("Handling echo message: %s", msg.MessageID)
	}

	// 解析参数
	var params EchoParameter
	if err := s.parseParameters(msg, &params); err != nil {
		s.sendErrorResponse(msg, "Invalid parameters: "+err.Error())
		return
	}

	// 处理业务逻辑（echo返回相同的消息）
	s.sendSuccessResponse(msg, map[string]interface{}{
		"echo": params.Message,
	})
}

// HandleAdd 处理加法服务
func (s *TemplateService) HandleAdd(msg *models.StreamMessage) {
	if s.debug {
		log.Printf("Handling add message: %s", msg.MessageID)
	}

	// 解析参数
	var params AddParameter
	if err := s.parseParameters(msg, &params); err != nil {
		s.sendErrorResponse(msg, "Invalid parameters: "+err.Error())
		return
	}

	// 处理业务逻辑（加法运算）
	result := params.Num1 + params.Num2
	s.sendSuccessResponse(msg, map[string]interface{}{
		"result": result,
		"num1":   params.Num1,
		"num2":   params.Num2,
	})
}

// HandleSubtract 处理减法服务
func (s *TemplateService) HandleSubtract(msg *models.StreamMessage) {
	if s.debug {
		log.Printf("Handling subtract message: %s", msg.MessageID)
	}

	// 解析参数
	var params SubtractParameter
	if err := s.parseParameters(msg, &params); err != nil {
		s.sendErrorResponse(msg, "Invalid parameters: "+err.Error())
		return
	}

	// 处理业务逻辑（减法运算）
	result := params.Num1 - params.Num2
	s.sendSuccessResponse(msg, map[string]interface{}{
		"result": result,
		"num1":   params.Num1,
		"num2":   params.Num2,
	})
}

// HandleMultiply 处理乘法服务
func (s *TemplateService) HandleMultiply(msg *models.StreamMessage) {
	if s.debug {
		log.Printf("Handling multiply message: %s", msg.MessageID)
	}

	// 解析参数
	var params MultiplyParameter
	if err := s.parseParameters(msg, &params); err != nil {
		s.sendErrorResponse(msg, "Invalid parameters: "+err.Error())
		return
	}

	// 处理业务逻辑（乘法运算）
	result := params.Num1 * params.Num2
	s.sendSuccessResponse(msg, map[string]interface{}{
		"result": result,
		"num1":   params.Num1,
		"num2":   params.Num2,
	})
}

// HandleDivide 处理除法服务
func (s *TemplateService) HandleDivide(msg *models.StreamMessage) {
	if s.debug {
		log.Printf("Handling divide message: %s", msg.MessageID)
	}

	// 解析参数
	var params DivideParameter
	if err := s.parseParameters(msg, &params); err != nil {
		s.sendErrorResponse(msg, "Invalid parameters: "+err.Error())
		return
	}

	// 检查除数是否为零
	if params.Num2 == 0 {
		s.sendErrorResponse(msg, "Division by zero is not allowed")
		return
	}

	// 处理业务逻辑（除法运算）
	result := params.Num1 / params.Num2
	s.sendSuccessResponse(msg, map[string]interface{}{
		"result": result,
		"num1":   params.Num1,
		"num2":   params.Num2,
	})
}

// parseParameters 解析消息参数
func (s *TemplateService) parseParameters(msg *models.StreamMessage, params interface{}) error {
	jsonData, err := json.Marshal(msg.Playload)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, params)
}

// sendResponse 发送响应消息
func (s *TemplateService) sendResponse(msg *models.StreamMessage, success bool, data map[string]interface{}) {
	responseMsg := &models.StreamMessage{
		MessageID:      strconv.Itoa(s.stream.GetMessageID()),
		ReplyID:        msg.MessageID,
		ServiceName:    "response",
		CallbackStream: msg.CallbackStream,
		Playload: map[string]interface{}{
			"success": success,
			"data":    data,
		},
	}

	s.stream.StreamPush(responseMsg, msg.CallbackStream)
}

// sendSuccessResponse 发送成功响应
func (s *TemplateService) sendSuccessResponse(msg *models.StreamMessage, data map[string]interface{}) {
	s.sendResponse(msg, true, data)
}

// sendErrorResponse 发送错误响应
func (s *TemplateService) sendErrorResponse(msg *models.StreamMessage, errorMsg string) {
	s.sendResponse(msg, false, map[string]interface{}{
		"error": errorMsg,
	})
}