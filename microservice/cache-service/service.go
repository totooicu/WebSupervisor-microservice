package main

import (
	"encoding/json"
	"log"

	"WebSupervisor/MyTool/redis"
	"WebSupervisor/MyTool/streamtool"
	"WebSupervisor/MyTool/streamtool/models"
)

type CacheService struct {
	redisClient *redis.Client
	config      *CacheConfig
	debug       bool
}

func NewCacheService(config *CacheConfig, debug bool) *CacheService {
	return &CacheService{
		redisClient: redis.NewClient(config.Redis.Host, config.Redis.Port, config.Redis.Password, config.Redis.DB),
		config:      config,
		debug:       debug,
	}
}

// Start 启动缓存服务
func (s *CacheService) Start() {
	log.Println("Starting cache service...")
	if s.debug {
		log.Printf("Debug mode enabled, config: %+v", s.config)
	}

	// 初始化StreamTool
	streamToolConfig := &models.StreamToolConfig{
		Redis: models.RedisConfig{
			Host:     s.config.Redis.Host,
			Port:     s.config.Redis.Port,
			Password: s.config.Redis.Password,
			DB:       s.config.Redis.DB,
		},
		Services: []models.ServiceConfig{
			{
				Name:          "compare_and_save",
				StreamName:    s.config.InputStream,
				ConsumerGroup: s.config.ConsumerGroup,
				ConsumerID:    "cache-consumer",
			},
			{
				Name:          "get",
				StreamName:    s.config.InputStream,
				ConsumerGroup: s.config.ConsumerGroup,
				ConsumerID:    "cache-consumer",
			},
			{
				Name:          "set",
				StreamName:    s.config.InputStream,
				ConsumerGroup: s.config.ConsumerGroup,
				ConsumerID:    "cache-consumer",
			},
			{
				Name:          "delete",
				StreamName:    s.config.InputStream,
				ConsumerGroup: s.config.ConsumerGroup,
				ConsumerID:    "cache-consumer",
			},
			{
				Name:          "get_and_set",
				StreamName:    s.config.InputStream,
				ConsumerGroup: s.config.ConsumerGroup,
				ConsumerID:    "cache-consumer",
			},
		},
		Debug: s.debug,
	}
	
	streamtool.InitStreamTool(streamToolConfig)
	
	// 启动消息网关
	st := streamtool.GetStreamTool()
	st.StartGateway(s.config.InputStream, s.config.ConsumerGroup, "cache-gateway")
	
	// 启动服务
	st.StartService("compare_and_save", s.handleCompareAndSave)
	st.StartService("get", s.handleGet)
	st.StartService("set", s.handleSet)
	st.StartService("delete", s.handleDelete)
	st.StartService("get_and_set", s.handleGetAndSet)
	
	log.Println("Cache service started successfully")
	
	// 保持服务运行
	select {}
}

// handleStreamMessage 处理流消息
func (s *CacheService) handleStreamMessage(msg *models.StreamMessage) {
	log.Printf("Processing task: %s", msg.MessageID)
	
	switch msg.ServiceName {
	case "compare_and_save":
		s.handleCompareAndSave(msg)
	case "get":
		s.handleGet(msg)
	case "set":
		s.handleSet(msg)
	case "delete":
		s.handleDelete(msg)
	case "get_and_set":
		s.handleGetAndSet(msg)
	default:
		if s.debug {
			log.Printf("Debug - Unsupported service: %s", msg.ServiceName)
		}
	}
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