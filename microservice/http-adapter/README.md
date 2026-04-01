# HTTP Adapter

HTTP Adapter 是一个双向适配器服务，用于连接 HTTP 协议和 Redis Stream 消息队列系统。它提供了两种工作模式：

- **服务器模式**：接收 HTTP 请求并转发到 Redis Stream，等待响应后返回 HTTP 响应
- **客户端模式**：从 Redis Stream 读取消息并转发到配置的 HTTP 端点

## 🚀 功能特性

### 核心功能
- ✅ **双向通信**：HTTP ↔ Redis Stream 的无缝转换
- ✅ **灵活配置**：支持多种配置选项和自定义映射
- ✅ **超时管理**：请求超时处理和响应等待机制
- ✅ **错误处理**：完善的错误处理和重试机制
- ✅ **安全特性**：CORS 支持和速率限制
- ✅ **监控日志**：详细的日志记录和运行状态监控

### 服务器模式特性
- 支持自定义 HTTP 请求头
- 动态 Stream 选择
- 响应超时控制
- 并发请求处理
- 统一的错误响应格式

### 客户端模式特性
- 多 Stream 监听
- 支持 GET/POST 方法
- 自定义 HTTP 头
- 指数退避重试
- 死信队列支持

## 🏗️ 架构设计

### 系统架构
```
┌─────────────────┐         ┌─────────────────┐         ┌─────────────────┐
│                 │         │                 │         │                 │
│   HTTP Client   │────────▶│   HTTP Adapter  │────────▶│  Redis Stream   │
│                 │         │                 │         │                 │
└─────────────────┘         └─────────────────┘         └─────────────────┘
         ▲                           │                           ▲
         │                           │                           │
         │                           │                           │
         │                           ▼                           │
┌─────────────────┐         ┌─────────────────┐         ┌─────────────────┐
│                 │         │                 │         │                 │
│   HTTP Server   │◀────────│   HTTP Adapter  │◀────────│  Microservices  │
│                 │         │                 │         │                 │
└─────────────────┘         └─────────────────┘         └─────────────────┘
```

### 核心组件
1. **配置管理器**：加载和验证配置文件
2. **HTTP 服务器**：处理 HTTP 请求和响应
3. **Redis 客户端**：管理 Redis 连接和 Stream 操作
4. **消息转换器**：在 HTTP 和 Stream 格式间转换
5. **响应等待器**：管理请求-响应匹配和超时
6. **流监控器**：客户端模式下监听多个 Stream

## ⚙️ 配置说明

### 配置文件结构

```json
{
  "redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
  },
  "server": {
    "enabled": true,
    "port": 8080,
    "host": "0.0.0.0",
    "default_input_stream": "http-adapter-input-stream",
    "default_response_stream": "http-adapter-response-stream",
    "timeout": 30,
    "debug": true,
    "max_concurrent_requests": 100,
    "request_id_header": "X-Request-ID",
    "service_name_header": "X-Service-Name",
    "stream_name_header": "X-Stream-Name",
    "callback_stream_header": "X-Callback-Stream",
    "timeout_header": "X-Timeout"
  },
  "client": {
    "enabled": true,
    "stream_mappings": [
      {
        "stream": "notification-service-output-stream",
        "http_endpoint": "http://localhost:3000/api/notifications",
        "method": "POST",
        "timeout": 5,
        "headers": {
          "Content-Type": "application/json",
          "Authorization": "Bearer your-api-token"
        },
        "retry": {
          "enabled": true,
          "max_attempts": 3,
          "backoff_ms": 1000
        },
        "dead_letter_stream": "http-adapter-dead-letter-stream"
      }
    ],
    "consumer_group": "http-adapter-client-group",
    "poll_interval_ms": 1000
  },
  "logging": {
    "level": "INFO",
    "format": "json",
    "file": "",
    "max_size_mb": 100,
    "max_backups": 5
  },
  "security": {
    "enable_cors": true,
    "allowed_origins": ["*"],
    "allowed_methods": ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
    "allowed_headers": ["*"],
    "enable_rate_limiting": false,
    "rate_limit_requests": 100,
    "rate_limit_window_seconds": 60
  }
}
```

### 配置项详细说明

#### Redis 配置
- `host`: Redis 服务器地址（默认：localhost）
- `port`: Redis 服务器端口（默认：6379）
- `password`: Redis 认证密码（可选）
- `db`: Redis 数据库索引（默认：0）

#### 服务器配置
- `enabled`: 是否启用服务器模式
- `port`: HTTP 服务器端口（默认：8080）
- `host`: HTTP 服务器绑定地址（默认：0.0.0.0）
- `default_input_stream`: 默认输入 Stream 名称
- `default_response_stream`: 默认响应 Stream 名称
- `timeout`: 默认请求超时时间（秒）
- `debug`: 是否启用调试模式
- `max_concurrent_requests`: 最大并发请求数
- `request_id_header`: 请求 ID 头名称
- `service_name_header`: 服务名称头名称
- `stream_name_header`: Stream 名称头名称
- `callback_stream_header`: 回调 Stream 头名称
- `timeout_header`: 超时设置头名称

#### 客户端配置
- `enabled`: 是否启用客户端模式
- `stream_mappings`: Stream 到 HTTP 端点的映射列表
- `consumer_group`: Redis Stream 消费者组名称
- `poll_interval_ms`: 轮询间隔（毫秒）

#### Stream 映射配置
- `stream`: Redis Stream 名称
- `http_endpoint`: 目标 HTTP 端点 URL
- `method`: HTTP 方法（GET/POST）
- `timeout`: HTTP 请求超时时间（秒）
- `headers`: 自定义 HTTP 头
- `retry`: 重试配置
  - `enabled`: 是否启用重试
  - `max_attempts`: 最大重试次数
  - `backoff_ms`: 退避时间（毫秒）
- `dead_letter_stream`: 死信队列 Stream 名称

#### 日志配置
- `level`: 日志级别（DEBUG/INFO/WARN/ERROR）
- `format`: 日志格式（json/text）
- `file`: 日志文件路径（空表示控制台输出）
- `max_size_mb`: 日志文件最大大小（MB）
- `max_backups`: 保留的日志文件数量

#### 安全配置
- `enable_cors`: 是否启用 CORS
- `allowed_origins`: 允许的源列表
- `allowed_methods`: 允许的 HTTP 方法
- `allowed_headers`: 允许的 HTTP 头
- `enable_rate_limiting`: 是否启用速率限制
- `rate_limit_requests`: 速率限制请求数
- `rate_limit_window_seconds`: 速率限制窗口时间（秒）

## 📡 HTTP API 使用指南

### 服务器模式 - 发送请求

#### 请求格式

**HTTP 请求示例**：
```http
POST /api/stream
Host: localhost:8080
Content-Type: application/json
X-Service-Name: add
X-Stream-Name: template-service-input-stream
X-Callback-Stream: http-adapter-response-stream
X-Timeout: 10

{
  "num1": 10,
  "num2": 20
}
```

#### 请求头说明

| 头名称 | 必填 | 说明 | 默认值 |
|--------|------|------|--------|
| `X-Service-Name` | ✅ | 要调用的服务名称 | - |
| `X-Stream-Name` | ❌ | 目标 Stream 名称 | 配置的 default_input_stream |
| `X-Callback-Stream` | ❌ | 响应回调 Stream | 配置的 default_response_stream |
| `X-Timeout` | ❌ | 请求超时时间（秒） | 配置的 timeout |
| `X-Request-ID` | ❌ | 请求唯一标识 | 自动生成 |

#### 响应格式

**成功响应**：
```json
{
  "success": true,
  "data": {
    "result": 30,
    "num1": 10,
    "num2": 20
  },
  "request_id": "req-123456",
  "processing_time": 150
}
```

**错误响应**：
```json
{
  "success": false,
  "error": "Service not found",
  "error_code": "SERVICE_NOT_FOUND",
  "request_id": "req-123456",
  "processing_time": 50
}
```

### 客户端模式 - Stream 监听

客户端模式会自动监听配置的 Stream，并将消息转发到对应的 HTTP 端点。

#### GET 请求处理
- 将 Stream 消息的 payload 转换为查询参数
- 发送 GET 请求到配置的 HTTP 端点

#### POST 请求处理
- 将 Stream 消息的 payload 作为请求体
- 发送 POST 请求到配置的 HTTP 端点

## 🚀 部署说明

### 环境要求
- Redis 5.0+
- Go 1.16+ 或 Python 3.8+
- 足够的内存和 CPU 资源

### 安装步骤

1. **克隆代码库**
```bash
git clone <repository-url>
cd http-adapter
```

2. **配置文件**
```bash
cp config.json.example config.json
# 编辑配置文件
```

3. **启动服务**
```bash
# Go 版本
go run main.go --config config.json --debug

# Python 版本
python main.py --config config.json --debug
```

4. **验证服务**
```bash
curl -X POST http://localhost:8080/api/stream \
  -H "Content-Type: application/json" \
  -H "X-Service-Name: echo" \
  -d '{"message": "test"}'
```

### Docker 部署

```bash
docker build -t http-adapter .
docker run -p 8080:8080 -v $(pwd)/config.json:/app/config.json http-adapter
```

## 📊 使用示例

### 示例 1：调用模板服务的加法功能

```bash
curl -X POST http://localhost:8080/api/stream \
  -H "Content-Type: application/json" \
  -H "X-Service-Name: add" \
  -H "X-Stream-Name: template-service-input-stream" \
  -d '{
    "num1": 10,
    "num2": 20
  }'
```

### 示例 2：自定义回调 Stream

```bash
curl -X POST http://localhost:8080/api/stream \
  -H "Content-Type: application/json" \
  -H "X-Service-Name: multiply" \
  -H "X-Callback-Stream: my-custom-response-stream" \
  -H "X-Timeout: 5" \
  -d '{
    "num1": 8,
    "num2": 7
  }'
```

### 示例 3：客户端模式配置

配置文件示例：
```json
{
  "client": {
    "enabled": true,
    "stream_mappings": [
      {
        "stream": "notification-stream",
        "http_endpoint": "http://webhook.example.com/notifications",
        "method": "POST",
        "headers": {
          "Authorization": "Bearer secret-token"
        }
      }
    ]
  }
}
```

## 🔧 故障排除

### 常见问题

1. **连接超时**
   - 检查 Redis 连接配置
   - 验证网络连通性
   - 检查防火墙设置

2. **服务未找到**
   - 确认 `X-Service-Name` 头设置正确
   - 检查目标微服务是否正常运行
   - 验证 Stream 名称配置

3. **响应超时**
   - 增加超时时间配置
   - 检查微服务处理时间
   - 验证 Redis Stream 消费者状态

4. **客户端模式不工作**
   - 检查 Stream 映射配置
   - 验证 HTTP 端点可达性
   - 查看日志中的错误信息

### 日志说明

日志级别：
- `DEBUG`: 详细的调试信息
- `INFO`: 正常运行信息
- `WARN`: 警告信息
- `ERROR`: 错误信息

关键日志指标：
- 请求处理时间
- 成功/失败统计
- 重试次数
- 超时统计

## 🤝 贡献指南

### 开发环境设置

1. 克隆代码库
2. 安装依赖
3. 运行测试

### 代码规范

- 遵循项目的代码风格
- 添加单元测试
- 更新文档

### 提交流程

1. 创建功能分支
2. 提交代码
3. 创建 Pull Request
4. 等待代码审查

## 📝 版本历史

- **v1.0.0**: 初始版本，支持基本的 HTTP ↔ Stream 转换
- **v1.1.0**: 添加客户端模式和重试机制
- **v1.2.0**: 增强安全特性和监控功能

## 📄 许可证

MIT License

## 🆘 支持

如有问题或建议，请创建 Issue 或联系维护者。

---

**注意**：本项目正在开发中，API 和配置格式可能会发生变化。请定期更新配置文件以保持兼容性。
