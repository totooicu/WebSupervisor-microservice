package model

// Message 定义消息结构
type Message struct {
	TaskID         string `json:"task_id"`
	ConsumerGroup  string `json:"consumer_group"`
	CallbackStream string `json:"callback_stream"`
	ServiceName    string `json:"service_name"`
	Playload       string `json:"playload"`
	ReplyID        string `json:"reply_id"`
}