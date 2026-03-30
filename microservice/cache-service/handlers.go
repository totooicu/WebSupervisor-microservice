package main

import (
	"encoding/json"
	"log"
	"strconv"

	"WebSupervisor/MyTool/streamtool"
	"WebSupervisor/MyTool/streamtool/models"
	"WebSupervisor/model"
)

func (s *CacheService) handleCompareAndSave(msg *models.StreamMessage) {
	// 解析参数
	var params model.CacheParameter
	playloadData, err := json.Marshal(msg.Playload)
	if err != nil {
		log.Printf("Error marshalling playload: %v", err)
		return
	}
	
	if err := json.Unmarshal(playloadData, &params); err != nil {
		log.Printf("Error unmarshalling playload: %v", err)
		return
	}

	key := s.generateKey(params.App, params.Key)

	if s.debug {
		log.Printf("Debug - Comparing data for key: %s", key)
	}

	var oldData interface{}
	err = s.redisClient.GetKey(key, &oldData)

	if err != nil {
		if s.debug {
			log.Printf("Debug - Key not found, saving new data: %v", err)
		}
	} else if s.debug {
		log.Printf("Debug - Found existing data for key: %s", key)
	}
	
	paramData := map[string]interface{}{"changed": false}
	if s.CompareData(oldData, params.Data) {
		log.Printf("Data changed, saving and sending notification")

		if err := s.SaveData(key, params.Data); err != nil {
			log.Printf("Error saving data: %v", err)
			return
		}

		if s.debug {
			log.Printf("Debug - Data saved successfully for key: %s", key)
		}

		paramData["changed"] = true
	} else if s.debug {
		log.Printf("Debug - Data unchanged for key: %s", key)
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
		log.Printf("Debug - Published notification to stream: %s", msg.CallbackStream)
	}
}

func (s *CacheService) handleGet(msg *models.StreamMessage) {
	// 解析参数
	var params model.CacheParameter
	playloadData, err := json.Marshal(msg.Playload)
	if err != nil {
		log.Printf("Error marshalling playload: %v", err)
		return
	}
	
	if err := json.Unmarshal(playloadData, &params); err != nil {
		log.Printf("Error unmarshalling playload: %v", err)
		return
	}

	key := s.generateKey(params.App, params.Key)

	if s.debug {
		log.Printf("Debug - Getting data for key: %s", key)
	}

	var data interface{}
	err = s.redisClient.GetKey(key, &data)

	if err != nil {
		if s.debug {
			log.Printf("Debug - Error getting data: %v", err)
		}
	} else if s.debug {
		log.Printf("Debug - Got data for key: %s", key)
	}

	paramData := map[string]interface{}{
		"key":   key,
		"data":  data,
		"error": err != nil,
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
		log.Printf("Debug - Published get result to stream: %s", msg.CallbackStream)
	}
}

func (s *CacheService) handleSet(msg *models.StreamMessage) {
	// 解析参数
	var params model.CacheParameter
	playloadData, err := json.Marshal(msg.Playload)
	if err != nil {
		log.Printf("Error marshalling playload: %v", err)
		return
	}
	
	if err := json.Unmarshal(playloadData, &params); err != nil {
		log.Printf("Error unmarshalling playload: %v", err)
		return
	}

	key := s.generateKey(params.App, params.Key)

	if s.debug {
		log.Printf("Debug - Setting data for key: %s", key)
	}

	if err := s.SaveData(key, params.Data); err != nil {
		log.Printf("Error saving data: %v", err)
		return
	}

	if s.debug {
		log.Printf("Debug - Data saved successfully for key: %s", key)
	}

	paramData := map[string]interface{}{
		"key":   key,
		"error": nil,
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
		log.Printf("Debug - Published set result to stream: %s", msg.CallbackStream)
	}
}

func (s *CacheService) handleDelete(msg *models.StreamMessage) {
	// 解析参数
	var params model.CacheParameter
	playloadData, err := json.Marshal(msg.Playload)
	if err != nil {
		log.Printf("Error marshalling playload: %v", err)
		return
	}
	
	if err := json.Unmarshal(playloadData, &params); err != nil {
		log.Printf("Error unmarshalling playload: %v", err)
		return
	}

	key := s.generateKey(params.App, params.Key)

	if s.debug {
		log.Printf("Debug - Deleting key: %s", key)
	}

	if err := s.redisClient.DeleteKey(key); err != nil {
		log.Printf("Error deleting key: %v", err)
		return
	}

	if s.debug {
		log.Printf("Debug - Key deleted successfully: %s", key)
	}

	paramData := map[string]interface{}{
		"key":   key,
		"error": nil,
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
		log.Printf("Debug - Published delete result to stream: %s", msg.CallbackStream)
	}
}

func (s *CacheService) handleGetAndSet(msg *models.StreamMessage) {
	// 解析参数
	var params model.CacheParameter
	playloadData, err := json.Marshal(msg.Playload)
	if err != nil {
		log.Printf("Error marshalling playload: %v", err)
		return
	}
	
	if err := json.Unmarshal(playloadData, &params); err != nil {
		log.Printf("Error unmarshalling playload: %v", err)
		return
	}

	key := s.generateKey(params.App, params.Key)

	if s.debug {
		log.Printf("Debug - Getting and setting data for key: %s", key)
	}

	var oldData interface{}
	err = s.redisClient.GetKey(key, &oldData)

	if err != nil {
		if s.debug {
			log.Printf("Debug - Key not found, setting new data: %v", err)
		}
	} else if s.debug {
		log.Printf("Debug - Found existing data for key: %s", key)
	}

	if err := s.SaveData(key, params.Data); err != nil {
		log.Printf("Error saving data: %v", err)
		return
	}

	if s.debug {
		log.Printf("Debug - Data saved successfully for key: %s", key)
	}

	paramData := map[string]interface{}{
		"key":      key,
		"old_data": oldData,
		"error":    nil,
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
		log.Printf("Debug - Published get_and_set result to stream: %s", msg.CallbackStream)
	}
}
