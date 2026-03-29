package model

type Message struct {
	TaskID         string `json:"task_id"`
	ConsumerGroup  string `json:"consumer_group"`
	CallbackStream string `json:"callback_stream"`
	ServiceName    string `json:"service_name"`
	Playload       string `json:"playload"`
}

type CrawlerParameter struct {
	URL        string                 `json:"url"`
	Method     string                 `json:"method"`
	Headers    map[string]string      `json:"headers"`
	Body       map[string]interface{} `json:"body"`
	StrPayload string                 `json:"str_payload"`
}

type ParserParameter struct {
	Content  string      `json:"content"`
	HTMLKeys []HTMLKey   `json:"HTMLKeys"`
	JSONKeys []JSONKey   `json:"JSONKeys"`
}

type HTMLKey struct {
	Left  string   `json:"left"`
	Right string   `json:"right"`
	Keys  []string `json:"keys"`
}

type JSONKey struct {
	Path []interface{} `json:"path"`
	Keys []string      `json:"keys"`
}

type CacheParameter struct {
	App  string      `json:"app"`
	Key  string      `json:"key"`
	Data interface{} `json:"data"`
}

type NotifierParameter struct {
	URL      string   `json:"url"`
	Subject  string   `json:"subject"`
	Content  string   `json:"content"`
	UserName string   `json:"userName"`
	Password string   `json:"password"`
	Tos      []string `json:"tos"`
}

type MonitorParameter struct {
	JobID         string           `json:"job_id"`
	Interval      int              `json:"interval"`
	CrawlerParams CrawlerParameter `json:"crawler_params"`
}
