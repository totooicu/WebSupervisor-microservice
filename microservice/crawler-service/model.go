package main

// RedisConfig Redis配置结构体
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// CrawlerConfig 爬虫服务配置结构体
type CrawlerConfig struct {
	Redis         RedisConfig `json:"redis"`
	ConsumerGroup string      `json:"consumer_group"`
	InputStream   string      `json:"input_stream"`
}

// CrawlerParameter 爬虫服务参数结构体
type CrawlerParameter struct {
	URL        string                 `json:"url"`
	Method     string                 `json:"method"`
	Headers    map[string]string      `json:"headers"`
	Body       map[string]interface{} `json:"body"`
	StrPayload string                 `json:"str_payload"`
}