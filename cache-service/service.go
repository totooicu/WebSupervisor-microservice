package main

import (
	"encoding/json"
	"log"

	"WebSupervisor/new/model"
	"WebSupervisor/new/redis"
	streams_communication "streams-communication"
	streams_model "streams-communication/model"
)

type CacheService struct {
	redisClient *redis.Client
	config      *model.CacheConfig
	debug       bool
}

func NewCacheService(config *model.CacheConfig, debug bool) *CacheService {
	return &CacheService{
		redisClient: redis.NewClient(config.Redis.Host, config.Redis.Port, config.Redis.Password, config.Redis.DB),
		config:      config,
		debug:       debug,
	}
}

// ProcessTask 实现StreamProcessor接口
func (s *CacheService) ProcessTask(task streams_model.Message) {
	log.Printf("Processing task: %s", task.TaskID)
	
	switch task.ServiceName {
	case "compare_and_save":
		s.handleCompareAndSave(task)
	case "get":
		s.handleGet(task)
	case "set":
		s.handleSet(task)
	case "delete":
		s.handleDelete(task)
	case "get_and_set":
		s.handleGetAndSet(task)
	default:
		if s.debug {
			log.Printf("Debug - Unsupported service: %s", task.ServiceName)
		}
	}
}

func (s *CacheService) Start() {
	log.Println("Starting cache service...")
	if s.debug {
		log.Printf("Debug mode enabled, config: %+v", s.config)
	}

	// 创建流消费者
	consumer := streams_communication.NewStreamConsumer(
		s.redisClient,
		s.config,
		s,
		s.debug,
		s.config.InputStream,
		s.config.ConsumerGroup,
		"cache-consumer",
		"compare_and_save",
	)
	
	// 启动流消费者
	consumer.Start()
}



func (s *CacheService) CompareData(oldData, newData interface{}) bool {
	if oldData == nil {
		return true
	}

	oldJSON, _ := json.Marshal(oldData)
	newJSON, _ := json.Marshal(newData)

	return string(oldJSON) != string(newJSON)
}

func (s *CacheService) SaveData(key string, data interface{}) error {
	return s.redisClient.SetKey(key, data, 0)
}

func (s *CacheService) generateKey(app, key string) string {
	return app + ":" + key
}