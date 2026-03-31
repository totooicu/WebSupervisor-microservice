# Notifier Service

通知服务，发送邮件通知，支持自定义邮件内容。

## 功能特性

- **邮件发送**: 支持发送邮件通知
- **自定义内容**: 支持自定义邮件主题和内容
- **多收件人**: 支持发送给多个收件人
- **SMTP配置**: 支持配置SMTP服务器
- **错误处理**: 完善的错误处理机制

## API接口

### 服务名称

- `send_email`: 发送邮件通知

### 请求参数

```json
{
  "subject": "Email Subject",
  "content": "Email Content",
  "userName": "sender@example.com",
  "password": "email_password",
  "tos": ["recipient1@example.com", "recipient2@example.com"]
}
```

### 参数说明

- `subject`: 邮件主题（必填）
- `content`: 邮件内容（必填）
- `userName`: 发件人邮箱（可选，使用配置中的默认值）
- `password`: 发件人密码（可选，使用配置中的默认值）
- `tos`: 收件人列表（可选，使用配置中的默认值）

### 响应格式

```json
{
  "success": true,
  "data": {
    "status": "sent",
    "tos": ["recipient@example.com"]
  }
}
```

## 配置说明

### config.json

```json
{
  "redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
  },
  "consumer_group": "notifier-service-group",
  "input_stream": "notifier-service-input-stream",
  "mail": {
    "userName": "sender@example.com",
    "password": "email_password",
    "tos": ["admin@example.com"]
  }
}
```

## 运行服务

```bash
# 启动服务
go run main.go service.go handlers.go utils.go --config ./config.json --debug

# 或使用批处理脚本
./run.bat
```

## 使用示例

### 基本邮件发送

```json
{
  "service": "send_email",
  "playload": {
    "subject": "Test Email",
    "content": "This is a test email from WebSupervisor"
  }
}
```

### 自定义收件人

```json
{
  "service": "send_email",
  "playload": {
    "subject": "Custom Email",
    "content": "Email with custom recipients",
    "tos": ["user1@example.com", "user2@example.com"]
  }
}
```

### 自定义发件人

```json
{
  "service": "send_email",
  "playload": {
    "subject": "Custom Sender",
    "content": "Email with custom sender",
    "userName": "custom@example.com",
    "password": "custom_password"
  }
}
```

## 开发说明

### 扩展功能

1. 修改 `handlers.go` 文件，添加新的通知方式
2. 在 `service.go` 中注册新的服务
3. 更新配置文件添加必要的配置项

### 最佳实践

- 使用安全的SMTP服务器
- 设置合理的邮件发送频率
- 实现邮件发送失败重试机制
- 添加邮件发送日志
- 验证邮箱格式

## 故障排除

### 常见问题

1. **邮件发送失败**
   - 检查SMTP服务器配置
   - 验证邮箱密码是否正确
   - 查看网络连接

2. **收件人格式错误**
   - 检查邮箱地址格式
   - 验证收件人列表是否为空

3. **服务无法启动**
   - 检查Redis连接
   - 验证配置文件格式

## 性能优化

- 使用邮件队列管理发送任务
- 实现批量邮件发送
- 添加邮件发送统计
- 使用异步处理提高吞吐量

## 安全注意事项

- 不要在代码中硬编码邮箱密码
- 使用加密方式存储邮箱密码
- 限制邮件发送频率
- 验证邮件内容的安全性
- 实现邮件发送权限控制