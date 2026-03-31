# Streams Library

一个基于Redis Streams的微服务通信库，提供核心的消息处理和通信功能。

## 功能特性

- **消息模型**：定义标准的Stream消息结构
- **Redis客户端**：封装Redis操作，提供流式数据处理
- **消息处理器**：提供消息处理和分发机制
- **类型安全**：使用Go结构体定义消息格式
- **可扩展性**：支持自定义消息处理逻辑

## 目录结构

```
streams-library/
├── model/
│   └── message.go      # 消息模型定义
├── redis/
│   └── client.go       # Redis客户端封装
├── processor.go        # 消息处理器
├── go.mod              # Go模块定义
└── README.md           # 说明文档
```

## 核心组件

### 1. 消息模型 (model/message.go)

定义标准的Stream消息结构：

```go
type StreamMessage struct {
    MessageID      string                 `json:"message_id"`
    ReplyID        string                 `json:"reply_id"`
    ServiceName    string                 `json:"service_name"`
    CallbackStream string                 `json:"callback_stream"`
    Playload       map[string]interface{} `json:"playload"`
}
```

### 2. Redis客户端 (redis/client.go)

提供Redis操作的封装：

- 连接管理
- Stream消息推送
- Stream消息读取
- 消费者组管理

### 3. 消息处理器 (processor.go)

提供消息处理和分发功能：

- 消息解析
- 服务路由
- 响应处理
- 错误处理

## 使用示例

### 基本用法

```go
package main

import (
    "log"
    
    "WebSupervisor/streams-library/model"
    "WebSupervisor/streams-library/redis"
)

func main() {
    // 创建Redis客户端
    client, err := redis.NewRedisClient(&redis.Config{
        Host:     "localhost",
        Port:     6379,
        Password: "",
        DB:       0,
    })
    if err != nil {
        log.Fatalf("Failed to create Redis client: %v", err)
    }
    defer client.Close()

    // 创建消息
    msg := &model.StreamMessage{
        ServiceName:    "test-service",
        CallbackStream: "response-stream",
        Playload: map[string]interface{}{
            "key": "value",
        },
    }

    // 发送消息
    err = client.XAdd("input-stream", msg)
    if err != nil {
        log.Fatalf("Failed to send message: %v", err)
    }

    log.Println("Message sent successfully")
}
```

### 消息处理

```go
func processMessage(msg *model.StreamMessage) {
    log.Printf("Processing message: %s", msg.ServiceName)
    
    // 处理业务逻辑
    result := map[string]interface{}{
        "status": "success",
        "data":   "processed",
    }
    
    // 创建响应消息
    responseMsg := &model.StreamMessage{
        MessageID:      generateMessageID(),
        ReplyID:        msg.MessageID,
        ServiceName:    "response",
        CallbackStream: msg.CallbackStream,
        Playload:       result,
    }
    
    // 发送响应
    client.XAdd(msg.CallbackStream, responseMsg)
}
```

## 依赖

- github.com/go-redis/redis/v8

## 配置说明

### Redis配置

```go
type Config struct {
    Host     string
    Port     int
    Password string
    DB       int
}
```

## 最佳实践

1. **连接管理**：使用defer关闭Redis连接
2. **错误处理**：妥善处理Redis操作错误
3. **消息格式**：遵循标准的消息结构
4. **并发安全**：注意并发访问的安全性

## 性能优化

- 使用连接池管理Redis连接
- 实现批量消息处理
- 合理设置超时时间
- 监控消息处理延迟

## 许可证

MIT License