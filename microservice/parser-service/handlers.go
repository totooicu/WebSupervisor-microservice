package main

import (
	"encoding/json"
	"WebSupervisor/MyTool"
	"fmt"
	"log"
	"strings"
	"WebSupervisor/model"

	streams_model "streams-communication/model"
)

func (s *ParserService) handleParseHTML(task streams_model.Message) {
	var params model.ParserParameter
	log.Printf(">>>handleParseHTML task.Playload: %v\n", task.Playload)

		
	if err := json.Unmarshal([]byte(task.Playload), &params); err != nil {
		log.Printf("Error unmarshalling playload: %v", err)
		return
	}
	log.Printf(">>>handleParseHTML params: %v\n", params)
	log.Printf(">>>handleParseHTML params.Content: %v\n", params.Content)
	log.Printf(">>>handleParseHTML params.HTMLKeys: %v\n", params.HTMLKeys)
	log.Printf(">>>handleParseHTML params.HTMLKeys[0].Keys: %v\n", params.HTMLKeys[0].Keys)

	results := MyTool.GetMid(params.Content, params.HTMLKeys[0].Left, params.HTMLKeys[0].Right,0)
	log.Printf(">>>handleParseHTML results: %V\n", results)
	if s.debug {
		log.Printf("Debug - HTML parse params: HTMLKeys=%v, content length=%d", params.HTMLKeys, len(params.Content))
	}

	

	if s.debug {
		log.Printf("Debug - HTML parse completed, found %d results", len(results))
	}

	paramData := map[string]interface{}{
		"parsed_data": results,
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

func (s *ParserService) handleParseJSON(task streams_model.Message) {
	var params model.ParserParameter
	if err := json.Unmarshal([]byte(task.Playload), &params); err != nil {
		log.Printf("Error unmarshalling playload: %v", err)
		return
	}

	// 将JSONKeys转换为utils.ParseJSON期望的格式 ["path1.path2", ...]
	var jsonKeys []string
	for _, jsonKey := range params.JSONKeys {
		var pathParts []string
		for _, pathPart := range jsonKey.Path {
			pathParts = append(pathParts, fmt.Sprintf("%v", pathPart))
		}
		pathStr := strings.Join(pathParts, ".")
		jsonKeys = append(jsonKeys, pathStr)
	}

	if s.debug {
		log.Printf("Debug - JSON parse params: JSONKeys=%v, content length=%d", params.JSONKeys, len(params.Content))
	}

	results := MyTool.ParseJSON(params.Content, jsonKeys)

	if s.debug {
		log.Printf("Debug - JSON parse completed, found %d results", len(results))
	}

	paramData := map[string]interface{}{
		"parsed_data": results,
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
