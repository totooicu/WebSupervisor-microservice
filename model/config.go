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

type TerminalConfig struct {
	BaseConfig
	DefaultCommand []string         `json:"defaultCommand"`
	Services       []ServiceConfig  `json:"services"`
	Gateway        GatewayConfig    `json:"gateway"`
	Security       SecurityConfig   `json:"security"`
	Monitoring     MonitoringConfig `json:"monitoring"`
}

type ServiceConfig struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Input       InputConfig  `json:"input"`
	Output      OutputConfig `json:"output"`
}

type InputConfig struct {
	Cmd           []string          `json:"cmd"`
	Params        map[string]string `json:"params"`
	AutoStart     bool              `json:"autoStart"`
	RestartPolicy string            `json:"restartPolicy"`
}

type OutputConfig struct {
	Adapters []AdapterConfig `json:"adapters"`
}

type AdapterConfig struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
	Params []ParamConfig          `json:"params,omitempty"`
}

type ParamConfig struct {
	Name string            `json:"name,omitempty"`
	Keys map[string]string `json:"keys,omitempty"`
}

type GatewayConfig struct {
	Port           int    `json:"port"`
	Host           string `json:"host"`
	Timeout        int    `json:"timeout"`
	MaxConnections int    `json:"maxConnections"`
}

type SecurityConfig struct {
	AllowedCommands []string `json:"allowedCommands"`
	Timeout         int      `json:"timeout"`
	MaxOutputSize   int      `json:"maxOutputSize"`
}

type MonitoringConfig struct {
	LogLevel    string            `json:"logLevel"`
	Metrics     bool              `json:"metrics"`
	HealthCheck HealthCheckConfig `json:"healthCheck"`
}

type HealthCheckConfig struct {
	Interval int `json:"interval"`
	Timeout  int `json:"timeout"`
}

type TemplateConfig struct {
	BaseConfig
}
