package main

// RedisConfig Redis配置结构体
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// ParserConfig 解析服务配置结构体
type ParserConfig struct {
	Redis         RedisConfig `json:"redis"`
	ConsumerGroup string      `json:"consumer_group"`
	InputStream   string      `json:"input_stream"`
}

// HTMLKey HTML解析键结构体
type HTMLKey struct {
	Left  string   `json:"left"`
	Right string   `json:"right"`
	Keys  []string `json:"key"`
}

// JSONKey JSON解析键结构体
type JSONKey struct {
	Path []interface{} `json:"path"`
	Keys []string      `json:"key"`
}

// ParserParameter 解析服务参数结构体
type ParserParameter struct {
	Content  string    `json:"content"`
	HTMLKeys []HTMLKey `json:"htmlKeys"`
	JSONKeys []JSONKey `json:"jsonKeys"`
}