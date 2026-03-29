package main

import (
	"encoding/json"
	"log"

	"WebSupervisor/model"
	"WebSupervisor/MyTool"
	streams_model "streams-communication/model"
)

func (s *CrawlerService) handleHttpRequest(task streams_model.Message) {
	var params model.CrawlerParameter
	if err := json.Unmarshal([]byte(task.Playload), &params); err != nil {
		log.Printf("Error unmarshalling playload: %v", err)
		return
	}

	if s.debug {
		log.Printf("Debug - HTTP request params: %+v", params)
	}

	var response string
	var err error

	httpClient := MyTool.NewHttpHeader(params.URL, params.Headers)

	switch params.Method {
	case "GET":
		if s.debug {
			log.Printf("Debug - Sending GET request to: %s", params.URL)
		}
		response = httpClient.Get("").GetBodyString()
	case "POST":
		if s.debug {
			log.Printf("Debug - Sending POST request to: %s", params.URL)
		}
		response = httpClient.Post(params.Body, params.StrPayload).GetBodyString()
	default:
		log.Printf("Unsupported method: %s", params.Method)
		return
	}

	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		return
	}

	if s.debug {
		log.Printf("Debug - HTTP response received, length: %d", len(response))
	}

	paramData := map[string]interface{}{
		"content": response,
	}
	paramBytes, _ := json.Marshal(paramData)

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
		log.Printf("Debug - Published result to stream: %s", task.CallbackStream)
	}
}
