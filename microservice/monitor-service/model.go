package main

// RedisConfig Redis配置结构体
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// MonitorConfig 监控服务配置结构体
type MonitorConfig struct {
	Redis                     RedisConfig `json:"redis"`
	ConsumerGroup             string      `json:"consumer_group"`
	InputStream               string      `json:"input_stream"`
	OutputStream              string      `json:"output_stream"`
	CrawlerServiceInputStream string      `json:"crawler-service_input_stream"`
	ParserServiceInputStream  string      `json:"parser-service_input_stream"`
	CacheServiceInputStream   string      `json:"cache-service_input_stream"`
	NotifierServiceInputStream string     `json:"notifier-service_input_stream"`
	JobsPath                  string      `json:"jobs_path"`
	IntervalSecond            int         `json:"interval_second"`
	HealthCheckPort           string      `json:"health_check_port"`
}

// Message 通用消息结构体
type Message struct {
	TaskID         string `json:"task_id"`
	ConsumerGroup  string `json:"consumer_group"`
	CallbackStream string `json:"callback_stream"`
	ServiceName    string `json:"service_name"`
	Playload       string `json:"playload"`
}

// CrawlerParameter 爬虫服务参数结构体
type CrawlerParameter struct {
	URL        string                 `json:"url"`
	Method     string                 `json:"method"`
	Headers    map[string]string      `json:"headers"`
	Body       map[string]interface{} `json:"body"`
	StrPayload string                 `json:"str_payload"`
}

// ParserParameter 解析服务参数结构体
type ParserParameter struct {
	Content  string                 `json:"content"`
	HTMLKeys []map[string]interface{} `json:"htmlKeys"`
	JSONKeys []map[string]interface{} `json:"jsonKeys"`
}

// CacheParameter 缓存服务参数结构体
type CacheParameter struct {
	App  string      `json:"app"`
	Key  string      `json:"key"`
	Data interface{} `json:"data"`
}

// MonitorParameter 监控服务参数结构体
type MonitorParameter struct {
	JobID         string           `json:"job_id"`
	Interval      int              `json:"interval"`
	CrawlerParams CrawlerParameter `json:"crawler_params"`
}