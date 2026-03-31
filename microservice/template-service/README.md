# Template Service 开发手册

## 概述

Template Service 是一个基于 Stream 架构的服务模板，用于演示如何正确接入 Stream 系统。本手册将指导您如何使用这个模板开发自己的微服务。

## 核心特性

- **手动服务注册**：用户可以手动注册服务处理器
- **带响应消息**：支持发送带响应的消息，使用 `response.get()` 阻塞获取
- **无响应消息**：支持发送不带响应的消息
- **统一接口**：提供清晰的API接口设计

## 目录结构

```
microservice/template-service/
├── main.go              # 应用入口文件
├── service.go           # 服务接口和核心方法
├── handlers.go          # 服务处理逻辑示例
├── example_usage.go     # API使用示例
├── config.json          # 配置文件
└── README.md            # 开发手册
```

## 核心API

### 1. 服务注册

```go
// RegisterHandler 手动注册服务处理器
func (s *TemplateService) RegisterHandler(serviceName string, handler func(msg *models.StreamMessage))
```

**参数说明**：
- `serviceName`: 服务名称
- `handler`: 消息处理函数

**使用示例**：
```go
service.RegisterHandler("echo", service.HandleEcho)
service.RegisterHandler("add", service.HandleAdd)
```

### 2. 发送带响应的消息

```go
// SendMessageWithResponse 发送带响应的消息
func (s *TemplateService) SendMessageWithResponse(msg *models.StreamMessage, stream string) *streamtool.Response
```

**参数说明**：
- `msg`: 消息对象
- `stream`: 目标流名称

**返回值**：
- `*streamtool.Response`: 响应对象，可使用 `.Get()` 阻塞获取响应

**使用示例**：
```go
response := service.SendMessageWithResponse(msg, "target-stream")
result := response.Get() // 阻塞等待响应
```

### 3. 发送不带响应的消息

```go
// SendMessageWithoutResponse 发送不带响应的消息
func (s *TemplateService) SendMessageWithoutResponse(msg *models.StreamMessage, stream string) bool
```

**参数说明**：
- `msg`: 消息对象
- `stream`: 目标流名称

**返回值**：
- `bool`: 发送是否成功

**使用示例**：
```go
success := service.SendMessageWithoutResponse(msg, "notification-stream")
```

## 开发步骤

### 1. 创建服务结构体

继承 `TemplateService` 或创建自己的服务结构体：

```go
type YourService struct {
    *TemplateService
}

func NewYourService(config *model.TemplateConfig, debug bool) *YourService {
    return &YourService{
        TemplateService: NewTemplateService(config, debug),
    }
}
```

### 2. 实现消息处理器

创建消息处理方法：

```go
func (s *YourService) HandleYourService(msg *models.StreamMessage) {
    // 解析参数
    var params model.YourParameter
    if err := s.parseParameters(msg, &params); err != nil {
        s.sendErrorResponse(msg, "Invalid parameters: "+err.Error())
        return
    }
    
    // 业务逻辑处理
    // ...
    
    // 发送响应
    s.sendSuccessResponse(msg, map[string]interface{}{
        "result": "success",
    })
}
```

### 3. 注册服务

在 `main.go` 中注册服务：

```go
service := NewYourService(&config, *debug)

// 手动注册服务处理器
service.RegisterHandler("your-service", service.HandleYourService)

// 启动服务
service.Start()
```

### 4. 使用API发送消息

#### 发送带响应的消息

```go
msg := &models.StreamMessage{
    ServiceName:    "target-service",
    CallbackStream: "response-stream",
    Playload: map[string]interface{}{
        "param1": "value1",
        "param2": "value2",
    },
}

response := service.SendMessageWithResponse(msg, "target-stream")
if response != nil {
    result := response.Get() // 阻塞获取响应
    fmt.Printf("Response: %v\n", result)
}
```

#### 发送不带响应的消息

```go
msg := &models.StreamMessage{
    ServiceName: "notification",
    Playload: map[string]interface{}{
        "type":    "info",
        "message": "System notification",
    },
}

success := service.SendMessageWithoutResponse(msg, "notification-stream")
if success {
    fmt.Println("Message sent successfully")
}
```

## 配置说明

### config.json 配置项

```json
{
  "redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
  },
  "consumer_group": "template-service-group",
  "input_stream": "template-service-input-stream"
}
```

- `redis`: Redis连接配置
- `consumer_group`: 消费者组名称
- `input_stream`: 输入流名称

## 运行服务

```bash
go run main.go service.go handlers.go example_usage.go --config ./config.json --debug
```

或使用 run.bat 脚本：

```bash
./run.bat
```

## 完整示例

### 服务端示例

```go
// 服务端：注册并处理消息
service.RegisterHandler("echo", func(msg *models.StreamMessage) {
    var params model.EchoParameter
    s.parseParameters(msg, &params)
    
    responseMsg := &models.StreamMessage{
        MessageID:      s.stream.GetMessageID(),
        ReplyID:        msg.MessageID,
        ServiceName:    "response",
        CallbackStream: msg.CallbackStream,
        Playload: map[string]interface{}{
            "success": true,
            "data": map[string]interface{}{
                "echo": params.Message,
            },
        },
    }
    
    s.stream.StreamPush(responseMsg, msg.CallbackStream)
})
```

### 客户端示例

```go
// 客户端：发送带响应的消息
msg := &models.StreamMessage{
    ServiceName:    "echo",
    CallbackStream: "response-stream",
    Playload: map[string]interface{}{
        "message": "Hello World",
    },
}

response := service.SendMessageWithResponse(msg, "template-service-input-stream")
result := response.Get() // 阻塞等待响应
fmt.Printf("Received response: %v\n", result)
```

## 最佳实践

### 1. 服务命名规范

- 使用小写字母和连字符
- 服务名应该清晰表达服务功能
- 示例：`user-authentication`, `data-processor`

### 2. 错误处理

- 使用统一的错误响应格式
- 提供详细的错误信息
- 记录错误日志以便调试

### 3. 消息设计

- 消息结构应该简洁明了
- 避免在消息中传递大量数据
- 使用合适的数据类型

### 4. 性能考虑

- 避免长时间阻塞操作
- 合理设置超时时间
- 使用异步处理模式

## 故障排除

### 常见问题

1. **服务无法启动**
   - 检查 Redis 连接配置
   - 验证配置文件格式
   - 查看端口占用情况

2. **消息发送失败**
   - 检查目标流是否存在
   - 验证消息格式
   - 查看 Redis 连接状态

3. **响应获取超时**
   - 检查服务是否正常运行
   - 验证回调流配置
   - 查看网络连接

## 总结

Template Service 提供了一个完整的微服务开发模板，遵循 Stream 架构设计。通过这个模板，您可以快速开发自己的微服务，并与其他服务进行通信。

核心要点：
- 使用 `RegisterHandler()` 手动注册服务
- 使用 `SendMessageWithResponse()` 发送带响应的消息
- 使用 `SendMessageWithoutResponse()` 发送不带响应的消息
- 使用 `response.Get()` 阻塞获取响应

按照本手册的指导，您可以轻松开发和部署基于 Stream 架构的微服务。