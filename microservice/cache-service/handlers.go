package main

import (
	"encoding/json"
	"log"
	streams_model "streams-communication/model"

	"WebSupervisor/model"
)

func (s *CacheService) handleCompareAndSave(task streams_model.Message) {
	var params model.CacheParameter
	if err := json.Unmarshal([]byte(task.Playload), &params); err != nil {
		log.Printf("Error unmarshalling playload: %v", err)
		return
	}

	key := s.generateKey(params.App, params.Key)

	if s.debug {
		log.Printf("Debug - Comparing data for key: %s", key)
	}

	var oldData interface{}
	err := s.redisClient.GetKey(key, &oldData)

	if err != nil {
		if s.debug {
			log.Printf("Debug - Key not found, saving new data: %v", err)
		}
	} else if s.debug {
		log.Printf("Debug - Found existing data for key: %s", key)
	}
	paramData := map[string]interface{}{"changed": false,}
	if s.CompareData(oldData, params.Data) {
		log.Printf("Data changed, saving and sending notification")

		if err := s.SaveData(key, params.Data); err != nil {
			log.Printf("Error saving data: %v", err)
			return
		}

		if s.debug {
			log.Printf("Debug - Data saved successfully for key: %s", key)
		}

		paramData["changed"]= true

	} else if s.debug {
		log.Printf("Debug - Data unchanged for key: %s", key)
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
			log.Printf("Debug - Published notification to stream: %s", task.CallbackStream)
		}
}

func (s *CacheService) handleGet(task streams_model.Message) {
	var params model.CacheParameter
	if err := json.Unmarshal([]byte(task.Playload), &params); err != nil {
		log.Printf("Error unmarshalling playload: %v", err)
		return
	}

	key := s.generateKey(params.App, params.Key)

	if s.debug {
		log.Printf("Debug - Getting data for key: %s", key)
	}

	var data interface{}
	err := s.redisClient.GetKey(key, &data)

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
		log.Printf("Debug - Published get result to stream: %s", task.CallbackStream)
	}
}

func (s *CacheService) handleSet(task streams_model.Message) {
	var params model.CacheParameter
	if err := json.Unmarshal([]byte(task.Playload), &params); err != nil {
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
		log.Printf("Debug - Published set result to stream: %s", task.CallbackStream)
	}
}

func (s *CacheService) handleDelete(task streams_model.Message) {
	var params model.CacheParameter
	if err := json.Unmarshal([]byte(task.Playload), &params); err != nil {
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
		log.Printf("Debug - Published delete result to stream: %s", task.CallbackStream)
	}
}

func (s *CacheService) handleGetAndSet(task streams_model.Message) {
	var params model.CacheParameter
	if err := json.Unmarshal([]byte(task.Playload), &params); err != nil {
		log.Printf("Error unmarshalling playload: %v", err)
		return
	}

	key := s.generateKey(params.App, params.Key)

	if s.debug {
		log.Printf("Debug - Getting and setting data for key: %s", key)
	}

	var oldData interface{}
	err := s.redisClient.GetKey(key, &oldData)

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
		log.Printf("Debug - Published get_and_set result to stream: %s", task.CallbackStream)
	}
}
