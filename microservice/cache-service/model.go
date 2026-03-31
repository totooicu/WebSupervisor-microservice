package main

// RedisConfig Redis配置结构体
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// PersistenceConfig 持久化配置结构体
type PersistenceConfig struct {
	AppPaths map[string]string `json:"app_paths"`
}

// CacheConfig 缓存服务配置结构体
type CacheConfig struct {
	Redis         RedisConfig         `json:"redis"`
	ConsumerGroup string              `json:"consumer_group"`
	InputStream   string              `json:"input_stream"`
	Persistence   PersistenceConfig   `json:"persistence"`
}

// CacheParameter 缓存服务参数结构体
type CacheParameter struct {
	App  string      `json:"app"`
	Key  string      `json:"key"`
	Data interface{} `json:"data"`
}