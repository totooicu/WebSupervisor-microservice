# Notifier Service 详细说明文档

## 服务概述

Notifier Service 是一个通知服务微服务，主要负责发送邮件通知。该服务通过 Redis Streams 接收邮件发送任务，支持从配置文件和请求参数中获取邮件配置，并提供邮件格式验证功能。服务采用异步处理模式，确保邮件发送任务的可靠执行。

## 核心功能说明

### 1. 邮件发送功能
- **配置优先级**：支持从请求参数覆盖配置文件中的邮件设置
- **邮箱格式验证**：对发送者和接收者邮箱进行格式验证
- **异步处理**：通过 Redis Streams 实现任务的异步处理和结果返回

### 2. 配置管理
- **默认配置**：从配置文件中读取默认的邮件发送配置
- **动态覆盖**：支持通过请求参数动态覆盖默认配置
- **安全验证**：确保邮箱格式正确，防止发送失败

## 参数说明

### 请求参数

| 参数名 | 数据类型 | 是否必填 | 默认值 | 详细描述 |
|--------|----------|----------|--------|----------|
| url | string | 否 | 空字符串 | 邮件发送的相关URL（保留字段） |
| subject | string | 是 | 无 | 邮件主题 |
| content | string | 是 | 无 | 邮件内容 |
| userName | string | 否 | 配置文件中的值 | 发送者邮箱地址 |
| password | string | 否 | 配置文件中的值 | 发送者邮箱密码 |
| tos | array[string] | 否 | 配置文件中的值 | 接收者邮箱地址数组 |

### 数据结构定义

```go
type NotifierParameter struct {
    URL      string   `json:"url"`
    Subject  string   `json:"subject"`
    Content  string   `json:"content"`
    UserName string   `json:"userName"`
    Password string   `json:"password"`
    Tos      []string `json:"tos"`
}
```

## 返回参数说明

### 成功响应

| 字段名 | 数据类型 | 说明 |
|--------|----------|------|
| success | boolean | 邮件发送是否成功 |
| tos | array[string] | 实际接收邮件的邮箱地址列表 |

### 返回数据结构

```json
{
"success": true,
"tos": ["recipient@example.com", "another@example.com"]
}
```

## 错误码说明

| 错误码 | 错误描述 | 可能原因 |
|--------|----------|----------|
| 400 | 参数解析失败 | 请求参数格式错误或缺失必填字段 |
| 400 | 无效的发送者邮箱格式 | userName 格式不符合邮箱规范 |
| 400 | 无效的接收者邮箱格式 | tos 中包含格式不正确的邮箱地址 |
| 500 | 邮件发送失败 | SMTP 服务器连接问题、认证失败或网络错误 |

## 调用示例

### 示例 1: 使用配置文件中的默认邮箱配置

```json
{
"task_id": "email-123",
"consumer_group": "monitor-group",
"callback_stream": "monitor-responses",
"service_name": "send_email",
"playload": "{\"subject\":\"系统告警通知\",\"content\":\"检测到网站异常，请及时处理\"}"
}
```

### 示例 2: 覆盖配置文件中的邮箱配置

```json
{
"task_id": "email-456",
"consumer_group": "monitor-group",
"callback_stream": "monitor-responses",
"service_name": "send_email",
"playload": "{\"subject\":\"自定义通知\",\"content\":\"这是一条自定义通知\",\"userName\":\"custom@example.com\",\"password\":\"custom-pass\",\"tos\":[\"user1@example.com\",\"user2@example.com\"]}"
}
```

### 示例 3: 仅指定接收者邮箱

```json
{
"task_id": "email-789",
"consumer_group": "monitor-group",
"callback_stream": "monitor-responses",
"service_name": "send_email",
"playload": "{\"subject\":\"定向通知\",\"content\":\"仅发送给指定用户\",\"tos\":[\"specific@example.com\"]}"
}
```

## 使用注意事项

### 1. 参数优先级规则
- 请求参数中的 `userName`、`password` 和 `tos` 会覆盖配置文件中的默认值
- 如果请求中未提供这些参数，则使用配置文件中的默认值

### 2. 邮箱格式验证
- 发送者邮箱必须包含 @ 符号和有效的域名部分
- 接收者邮箱列表中的每个邮箱都必须通过格式验证
- 邮箱验证失败会导致整个邮件发送任务失败

### 3. 配置要求
- Redis 连接配置必须正确，包括主机、端口、密码和数据库索引
- 邮件配置必须包含有效的发送者邮箱和密码
- 消费者组和输入流名称必须与配置文件一致

### 4. 安全考虑
- 避免在日志中记录邮箱密码等敏感信息
- 建议使用环境变量或加密配置文件存储邮箱密码
- 考虑实现邮件发送频率限制，防止滥用

### 5. 性能建议
- 对于大量邮件发送任务，考虑实现批量处理机制
- 建议配置合理的超时时间，避免长时间阻塞
- 实现重试机制，处理临时的网络问题

## 配置文件说明

```json
{
"redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
},
"consumer_group": "notifier-group",
"input_stream": "notifier-tasks",
"mail": {
    "userName": "default@example.com",
    "password": "default-password",
    "tos": ["admin@example.com", "support@example.com"]
}
}
```

## 运行方式

```bash
# 开发模式（带调试日志）
./notifier-service.exe -debug -config ./notifier-service/config.json

# 生产模式
./notifier-service.exe -config ./notifier-service/config.json
```