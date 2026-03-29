package model

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type MailConfig struct {
	UserName string   `json:"userName"`
	Password string   `json:"password"`
	Tos      []string `json:"tos"`
}

type PersistenceConfig struct {
	AppPaths map[string]string `json:"app_paths"`
}

type BaseConfig struct {
	Redis         RedisConfig `json:"redis"`
	ConsumerGroup string      `json:"consumer_group"`
	InputStream   string      `json:"input_stream"`
}

type CrawlerConfig struct {
	BaseConfig
}

type ParserConfig struct {
	BaseConfig
}

type CacheConfig struct {
	BaseConfig
	Persistence PersistenceConfig `json:"persistence"`
}

type NotifierConfig struct {
	BaseConfig
	Mail MailConfig `json:"mail"`
}

type MonitorConfig struct {
	BaseConfig
	OutputStream               string `json:"output_stream"`
	CrawlerServiceInputStream  string `json:"crawler-service_input_stream"`
	ParserServiceInputStream   string `json:"parser-service_input_stream"`
	CacheServiceInputStream    string `json:"cache-service_input_stream"`
	NotifierServiceInputStream string `json:"notifier-service_input_stream"`
	JobsPath                   string `json:"jobs_path"`
	IntervalSecond             int    `json:"interval_second"`
	HealthCheckPort            string `json:"health_check_port"`
}
