package streams

import (
	"log"

	"streams-communication/model"
	"streams-communication/redis"
)

// StreamProcessor 定义流处理接口
type StreamProcessor interface {
	ProcessTask(task model.Message)
}

// StreamConsumer 流消费者
type StreamConsumer struct {
	redisClient redis.RedisClient
	config      interface{}
	processor   StreamProcessor
	debug       bool
	streamName  string
	groupName   string
	consumerID  string
	defaultServiceName string
}

// NewStreamConsumer 创建新的流消费者
func NewStreamConsumer(redisClient redis.RedisClient, config interface{}, processor StreamProcessor, debug bool, streamName, groupName, consumerID, defaultServiceName string) *StreamConsumer {
	return &StreamConsumer{
		redisClient: redisClient,
		config:      config,
		processor:   processor,
		debug:       debug,
		streamName:  streamName,
		groupName:   groupName,
		consumerID:  consumerID,
		defaultServiceName: defaultServiceName,
	}
}

// Start 启动流消费者
func (sc *StreamConsumer) Start() {
	log.Printf("Starting stream consumer for stream: %s, group: %s", sc.streamName, sc.groupName)
	
	if err := sc.redisClient.CreateConsumerGroup(sc.streamName, sc.groupName); err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}

	if sc.debug {
		log.Printf("Debug - Created consumer group: %s for stream: %s", sc.groupName, sc.streamName)
	}

	for {
		sc.processMessages()
	}
}

// processMessages 处理消息
func (sc *StreamConsumer) processMessages() {
	messages, err := sc.redisClient.ReadMessages(sc.streamName, sc.groupName, sc.consumerID, 1, 0)
	if err != nil {
		if err.Error() != "redis: nil" && sc.debug {
			log.Printf("Debug - Error reading messages: %v", err)
		}
		return
	}

	if sc.debug {
		log.Printf("Debug - Received %d messages", len(messages))
	}

	for _, msg := range messages {
		task := sc.parseMessage(msg)
		if task.TaskID == "" {
			log.Printf("Skipping message with empty TaskID")
			continue
		}

		log.Printf("Processing task: %s", task.TaskID)
		
		if sc.debug {
			log.Printf("Debug - Task details: %+v", task)
			log.Printf("Debug - ServiceName: '%s'", task.ServiceName)
		}

		// 如果没有设置ServiceName，使用默认值
		if task.ServiceName == "" {
			task.ServiceName = sc.defaultServiceName
			if sc.debug {
				log.Printf("Debug - No ServiceName provided, using default: %s", sc.defaultServiceName)
			}
		}

		// 调用处理器处理任务
		sc.processor.ProcessTask(task)

		// 确认消息（删除）
		if err := sc.redisClient.AcknowledgeMessage(sc.streamName, sc.groupName, msg.ID); err != nil {
			log.Printf("Error acknowledging message: %v", err)
		} else if sc.debug {
			log.Printf("Debug - Acknowledged message: %s", msg.ID)
		}
	}
}

// parseMessage 解析消息
func (sc *StreamConsumer) parseMessage(msg redis.XMessage) model.Message {
	var task model.Message

	if sc.debug {
		log.Printf("Debug - Parsing message with values: %+v", msg.Values)
	}

	// 新格式：直接从消息字段构造Message
	if callbackStream, ok := msg.Values["callback_stream"].(string); ok {
		task.CallbackStream = callbackStream
	}
	if consumerGroup, ok := msg.Values["consumer_group"].(string); ok {
		task.ConsumerGroup = consumerGroup
	}
	// 同时支持旧的parameter字段和新的playload字段
	if playload, ok := msg.Values["playload"].(string); ok {
		task.Playload = playload
	} else if parameter, ok := msg.Values["parameter"].(string); ok {
		task.Playload = parameter
	} else {
		if sc.debug {
			log.Printf("Debug - No playload or parameter field found")
			for k, v := range msg.Values {
				log.Printf("Debug - Field: %s, Type: %T, Value: %v", k, v, v)
			}
		}
	}
	if serviceName, ok := msg.Values["service_name"].(string); ok {
		task.ServiceName = serviceName
	}
	// 同时支持旧的task_id字段和新的message_id字段
	if messageID, ok := msg.Values["message_id"].(string); ok {
		task.TaskID = messageID
	} else if taskID, ok := msg.Values["task_id"].(string); ok {
		task.TaskID = taskID
	} else {
		if sc.debug {
			log.Printf("Debug - No message_id or task_id field found")
		}
	}
	if replyID, ok := msg.Values["reply_id"].(string); ok {
		task.ReplyID = replyID
	}

	return task
}