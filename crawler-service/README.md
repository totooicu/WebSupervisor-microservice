# Crawler Service 详细说明文档

## 服务概述

Crawler Service 是一个网络爬取微服务，通过 Redis Streams 接收任务并执行 HTTP 请求。该服务支持 GET 和 POST 请求方法，能够处理不同格式的请求体，并将响应结果发送到指定的回调流中。服务采用异步处理模式，确保高并发和可靠的任务处理能力。

## 核心功能说明

### 1. HTTP 请求处理
- **支持的请求方法**：GET、POST
- **请求体格式**：支持 JSON 对象和字符串格式
- **自定义请求头**：允许设置自定义 HTTP 请求头
- **异步处理**：通过 Redis Streams 实现任务的异步处理和结果返回

### 2. 错误处理机制
- 参数解析错误处理
- HTTP 请求失败处理
- 不支持的 HTTP 方法处理
- 调试模式支持，提供详细的日志输出

## 参数说明

### 请求参数

| 参数名 | 数据类型 | 是否必填 | 默认值 | 详细描述 |
|--------|----------|----------|--------|----------|
| url | string | 是 | 无 | 请求的目标 URL 地址，必须包含完整的协议（http:// 或 https://） |
| method | string | 是 | 无 | HTTP 请求方法，支持的值：GET、POST |
| headers | map[string]string | 否 | 空对象 | HTTP 请求头，用于设置自定义请求头信息，如 Content-Type、Authorization 等 |
| body | map[string]interface{} | 否 | null | JSON 格式的请求体，仅在 method 为 POST 时有效 |
| str_payload | string | 否 | 空字符串 | 字符串格式的请求体，仅在 method 为 POST 时有效，优先级高于 body |

### 数据结构定义

```go
type CrawlerParameter struct {
    URL        string                 `json:"url"`
    Method     string                 `json:"method"`
    Headers    map[string]string      `json:"headers"`
    Body       map[string]interface{} `json:"body"`
    StrPayload string                 `json:"str_payload"`
}
```

## 返回参数说明

### 成功响应

| 字段名 | 数据类型 | 说明 |
|--------|----------|------|
| content | string | HTTP 响应的内容，包含服务器返回的完整响应体 |

### 返回数据结构

```json
{
"content": "HTTP响应内容"
}
```

## 错误码说明

| 错误码 | 错误描述 | 可能原因 |
|--------|----------|----------|
| 400 | 参数解析失败 | 请求参数格式错误或缺失必填字段 |
| 405 | 不支持的HTTP方法 | 请求方法不是 GET 或 POST |
| 500 | HTTP请求失败 | 网络连接问题、目标服务器错误或超时 |

## 调用示例

### 示例 1: GET 请求

```json
{
"task_id": "task-123",
"consumer_group": "monitor-group",
"callback_stream": "monitor-responses",
"service_name": "http_request",
"playload": "{\"url\":\"https://api.example.com/data\",\"method\":\"GET\",\"headers\":{\"Accept\":\"application/json\"}}"
}
```

### 示例 2: POST 请求（JSON 格式）

```json
{
"task_id": "task-456",
"consumer_group": "monitor-group",
"callback_stream": "monitor-responses",
"service_name": "http_request",
"playload": "{\"url\":\"https://api.example.com/submit\",\"method\":\"POST\",\"headers\":{\"Content-Type\":\"application/json\"},\"body\":{\"username\":\"test\",\"password\":\"123456\"}}"
}
```

### 示例 3: POST 请求（字符串格式）

```json
{
"task_id": "task-789",
"consumer_group": "monitor-group",
"callback_stream": "monitor-responses",
"service_name": "http_request",
"playload": "{\"url\":\"https://api.example.com/submit\",\"method\":\"POST\",\"headers\":{\"Content-Type\":\"application/x-www-form-urlencoded\"},\"str_payload\":\"username=test&password=123456\"}"
}
```

## 使用注意事项

### 1. 参数优先级
- 当同时提供 `body` 和 `str_payload` 时，`str_payload` 优先级更高
- `method` 参数必须严格匹配 "GET" 或 "POST"（大小写敏感）

### 2. 配置要求
- Redis 连接配置必须正确，包括主机、端口、密码和数据库索引
- 消费者组和输入流名称必须与配置文件一致

### 3. 性能建议
- 避免在短时间内发送大量请求，可能导致目标服务器限流
- 对于大型响应内容，注意内存使用情况
- 建议在生产环境中启用调试日志进行问题排查

### 4. 安全考虑
- 不要在请求头或请求体中包含敏感信息
- 确保目标 URL 是可信的，避免安全风险
- 考虑实现请求超时机制，防止长时间阻塞

## 配置文件说明

```json
{
"redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
},
"consumer_group": "crawler-group",
"input_stream": "crawler-tasks"
}
```

## 运行方式

```bash
# 开发模式（带调试日志）
./crawler-service.exe -debug -config ./crawler-service/config.json

# 生产模式
./crawler-service.exe -config ./crawler-service/config.json
```