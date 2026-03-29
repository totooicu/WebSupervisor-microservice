package redis

import "github.com/go-redis/redis/v8"

// XMessage 定义Redis XMessage类型
type XMessage = redis.XMessage

// RedisClient 定义Redis客户端接口
type RedisClient interface {
	CreateConsumerGroup(stream, group string) error
	ReadMessages(stream, group, consumer string, count, block int64) ([]XMessage, error)
	AcknowledgeMessage(stream, group, messageID string) error
}