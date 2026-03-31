package main

// RedisConfig Redis配置结构体
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// MailConfig 邮件配置结构体
type MailConfig struct {
	UserName string   `json:"userName"`
	Password string   `json:"password"`
	Tos      []string `json:"tos"`
}

// NotifierConfig 通知服务配置结构体
type NotifierConfig struct {
	Redis         RedisConfig `json:"redis"`
	ConsumerGroup string      `json:"consumer_group"`
	InputStream   string      `json:"input_stream"`
	Mail          MailConfig  `json:"mail"`
}

// NotifierParameter 通知服务参数结构体
type NotifierParameter struct {
	URL      string   `json:"url"`
	Subject  string   `json:"subject"`
	Content  string   `json:"content"`
	UserName string   `json:"userName"`
	Password string   `json:"password"`
	Tos      []string `json:"tos"`
}