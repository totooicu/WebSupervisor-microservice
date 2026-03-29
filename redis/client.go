package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// XMessage 导出Redis XMessage类型
type XMessage = redis.XMessage

type Client struct {
	client *redis.Client
	ctx    context.Context
}

func NewClient(host string, port int, password string, db int) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	return &Client{
		client: client,
		ctx:    ctx,
	}
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) PublishMessage(stream string, message interface{}) error {
	// 将消息转换为map[string]interface{}
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	
	var msgMap map[string]interface{}
	if err := json.Unmarshal(data, &msgMap); err != nil {
		return err
	}

	_, err = c.client.XAdd(c.ctx, &redis.XAddArgs{
		Stream: stream,
		Values: msgMap,
	}).Result()

	return err
}

func (c *Client) CreateConsumerGroup(stream, group string) error {
	_, err := c.client.XGroupCreateMkStream(c.ctx, stream, group, "$").Result()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return err
	}
	return nil
}

func (c *Client) ReadMessages(stream, group, consumer string, count, block int64) ([]redis.XMessage, error) {
	xstreams, err := c.client.XReadGroup(c.ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{stream, ">"},
		Count:    count,
		Block:    time.Duration(block),
	}).Result()

	if err != nil {
		return nil, err
	}

	var messages []redis.XMessage
	for _, xstream := range xstreams {
		messages = append(messages, xstream.Messages...)
	}

	return messages, nil
}

func (c *Client) AcknowledgeMessage(stream, group, id string) error {
	return c.client.XAck(c.ctx, stream, group, id).Err()
}

func (c *Client) SetKey(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(c.ctx, key, string(data), expiration).Err()
}

func (c *Client) GetKey(key string, dest interface{}) error {
	data, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

func (c *Client) DeleteKey(key string) error {
	return c.client.Del(c.ctx, key).Err()
}
