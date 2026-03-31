package main

import (
	"encoding/json"
	"log"
	"strconv"

	http "WebSupervisor/MyTool/http"
	"WebSupervisor/MyTool/streamtool"
	"WebSupervisor/MyTool/streamtool/models"
)

func (s *CrawlerService) handleHttpRequest(msg *models.StreamMessage) {
	log.Printf(">>>rocessing msg.Playload: %v", msg.Playload)
	// 检查url字段是否存在
	if _, ok := msg.Playload["url"]; !ok {
		log.Printf("Error: url field not found in playload")
		return
	}

	// 解析参数
	var params CrawlerParameter
	playloadData, err := json.Marshal(msg.Playload)
	if err != nil {
		log.Printf("Error marshalling playload: %v", err)
		return
	}

	if err := json.Unmarshal(playloadData, &params); err != nil {
		log.Printf("Error unmarshalling playload: %v", err)
		return
	}

	if s.debug {
		log.Printf("Debug - HTTP request params: %+v", params)
	}

	var response string
	var statusCode int

	httpClient := http.NewHttpHeader(params.URL, params.Headers)

	switch params.Method {
	case "GET":
		if s.debug {
			log.Printf("Debug - Sending GET request to: %s", params.URL)
		}
		httpClient.Get("").Read()
		response = httpClient.GetBodyString()
		statusCode = httpClient.GetStatusCode()
	case "POST":
		if s.debug {
			log.Printf("Debug - Sending POST request to: %s", params.URL)
		}
		httpClient.Post(params.Body, params.StrPayload).Read()
		response = httpClient.GetBodyString()
		statusCode = httpClient.GetStatusCode()
	default:
		log.Printf("Unsupported method: %s", params.Method)
		
		return
	}

	// 检查状态码是否有效
	if statusCode == 0 {
		log.Printf("Warning: Status code is 0, request may have failed")
	} else if s.debug {
		log.Printf("Debug - HTTP response status code: %d", statusCode)
	}

	if s.debug {
		log.Printf("Debug - HTTP response received, length: %d", len(response))
	}
	log.Printf(">>>handleHttpRequest response: %s", response)
	// 构造响应消息
	paramData := map[string]interface{}{
		"content": response,
		"status": statusCode,
	}

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
		log.Printf("Debug - Published result to stream: %s", msg.CallbackStream)
	}
}
