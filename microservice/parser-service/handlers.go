package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"strconv"

	"WebSupervisor/MyTool"
	"WebSupervisor/MyTool/streamtool"
	"WebSupervisor/MyTool/streamtool/models"
	"WebSupervisor/model"
)

func (s *ParserService) handleParseHTML(msg *models.StreamMessage) {
	// 解析参数
	var params model.ParserParameter
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
		log.Printf("Debug - HTML parse params: HTMLKeys=%v, content length=%d", params.HTMLKeys, len(params.Content))
	}

	results := MyTool.GetMid(params.Content, params.HTMLKeys[0].Left, params.HTMLKeys[0].Right, 0)

	if s.debug {
		log.Printf("Debug - HTML parse completed, found %d results", len(results))
	}

	paramData := map[string]interface{}{
		"parsed_data": results,
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
		log.Printf("Debug - Published result to stream: %s", msg.CallbackStream)
	}
}

func (s *ParserService) handleParseJSON(msg *models.StreamMessage) {
	// 解析参数
	var params model.ParserParameter
	playloadData, err := json.Marshal(msg.Playload)
	if err != nil {
		log.Printf("Error marshalling playload: %v", err)
		return
	}
	
	if err := json.Unmarshal(playloadData, &params); err != nil {
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
		log.Printf("Debug - Published result to stream: %s", msg.CallbackStream)
	}
}
