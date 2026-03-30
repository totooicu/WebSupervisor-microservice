package models

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// StreamMessage 流消息结构
type StreamMessage struct {
	MessageID   string                 `json:"message_id"`
	ReplyID     string                 `json:"reply_id"`
	ServiceName string                 `json:"service_name"`
	CallbackStream string              `json:"callback_stream"`
	Playload    map[string]interface{} `json:"playload"`
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	Name         string `json:"name"`
	StreamName   string `json:"stream_name"`
	ConsumerGroup string `json:"consumer_group"`
	ConsumerID   string `json:"consumer_id"`
}

// StreamToolConfig Stream工具配置
type StreamToolConfig struct {
	Redis     RedisConfig     `json:"redis"`
	Services  []ServiceConfig `json:"services"`
	Debug     bool            `json:"debug"`
}