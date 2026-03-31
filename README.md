# WebSupervisor Microservice

WebSupervisor 是一个基于 Go 语言开发的微服务架构系统，用于监控和管理 Web 服务。系统采用 Redis Stream 作为通信中间件，实现了高性能、可靠的服务间通信。

## 项目架构

项目采用微服务架构，包含多个独立的服务模块，通过 Redis Stream 进行通信。整体架构清晰，职责明确，易于扩展和维护。

### 核心组件

#### 1. 工具包集合 (MyTool)
- `config` - 配置管理，支持JSON配置文件加载和环境变量
- `http` - HTTP客户端，支持GET/POST请求和自定义Headers
- `json` - JSON操作，提供序列化和反序列化功能
- `parser` - 数据解析，支持HTML和JSON内容提取
- `redis` - Redis客户端，封装Redis操作
- `streamtool` - Stream通信工具，提供请求-响应模式
- `streams-manager` - Redis Streams管理工具，用于调试和管理
- `sync` - 同步工具，提供线程安全的操作
- `string` - 字符串操作，提供常用的字符串处理功能
- `email` - 邮件发送，支持SMTP配置
- `file` - 文件操作，支持文件读写和路径处理
- `array` - 数组操作，提供数组处理功能
- `set` - 集合操作，提供集合处理功能

#### 2. 微服务模块 (microservice)
- `cache-service` - 缓存服务，提供数据缓存功能
- `crawler-service` - 爬虫服务，用于爬取网页内容
- `monitor-service` - 监控服务，负责监控其他服务的健康状态和任务调度
- `notifier-service` - 通知服务，发送邮件通知
- `parser-service` - 解析服务，提供HTML和JSON解析功能
- `template-service` - 服务模板，提供微服务开发的标准模板

#### 3. 核心库 (streams-library)
- 基础消息模型定义
- Redis客户端封装
- 消息处理器实现

#### 4. 公共模型 (model)
- 数据模型定义
- 消息结构定义

## 技术栈

- **语言**: Go 1.23+
- **通信**: Redis Stream
- **缓存**: Redis
- **配置**: JSON配置文件
- **依赖管理**: Go Modules

## 快速开始

### 环境要求

- Go 1.23 或更高版本
- Redis 6.0 或更高版本
- Windows/Linux/macOS 操作系统

### 安装依赖

```bash
go mod tidy
```

### 启动服务

#### 启动单个服务

```bash
# 启动缓存服务
cd microservice/cache-service
go run main.go service.go handlers.go utils.go --config ./config.json --debug

# 启动爬虫服务
cd microservice/crawler-service
go run main.go service.go handlers.go utils.go --config ./config.json --debug

# 启动监控服务
cd microservice/monitor-service
go run main.go service.go service_communicator.go task_scheduler.go config_parser.go health_check.go utils.go --config ./config.json --debug

# 启动通知服务
cd microservice/notifier-service
go run main.go service.go handlers.go utils.go --config ./config.json --debug

# 启动解析服务
cd microservice/parser-service
go run main.go service.go handlers.go utils.go --config ./config.json --debug

# 启动模板服务
cd microservice/template-service
go run main.go service.go handlers.go example_usage.go --config ./config.json --debug
```

#### 使用批处理脚本

```bash
# Windows
cd microservice/[service-name]
./run.bat

# Linux/macOS
cd microservice/[service-name]
./run.sh
```

### 构建服务

```bash
# 构建单个服务
go build -o run.exe ./microservice/[service-name]

# 构建所有服务
./build.bat
```

### 运行所有服务

```bash
# 使用统一的启动脚本
./runs.bat
```

## 服务说明

### 1. Cache Service
缓存服务，提供数据缓存功能，支持数据的存储、获取、删除等操作。
- **主要功能**: 数据存储、获取、删除、比较并保存、原子操作
- **配置文件**: `microservice/cache-service/config.json`
- **详细文档**: [cache-service/README.md](microservice/cache-service/README.md)

### 2. Crawler Service
爬虫服务，用于爬取网页内容，支持GET和POST请求。
- **主要功能**: HTTP请求、自定义Headers、请求参数、响应解析
- **配置文件**: `microservice/crawler-service/config.json`
- **详细文档**: [crawler-service/README.md](microservice/crawler-service/README.md)

### 3. Monitor Service
监控服务，负责监控其他服务的健康状态，调度任务执行。
- **主要功能**: 服务监控、任务调度、健康检查、配置解析
- **配置文件**: `microservice/monitor-service/config.json`
- **任务配置**: `microservice/monitor-service/jobs.json`
- **详细文档**: [monitor-service/README.md](microservice/monitor-service/README.md)

### 4. Notifier Service
通知服务，发送邮件通知，支持自定义邮件内容。
- **主要功能**: 邮件发送、自定义内容、多收件人、SMTP配置
- **配置文件**: `microservice/notifier-service/config.json`
- **详细文档**: [notifier-service/README.md](microservice/notifier-service/README.md)

### 5. Parser Service
解析服务，提供HTML和JSON解析功能，支持提取指定内容。
- **主要功能**: HTML解析、JSON解析、多规则解析、灵活配置
- **配置文件**: `microservice/parser-service/config.json`
- **详细文档**: [parser-service/README.md](microservice/parser-service/README.md)

### 6. Template Service
服务模板，提供微服务开发的标准模板，包含完整的服务注册和消息处理机制。
- **主要功能**: 手动服务注册、带响应消息、无响应消息、统一接口
- **配置文件**: `microservice/template-service/config.json`
- **详细文档**: [template-service/README.md](microservice/template-service/README.md)

## 开发指南

### 创建新服务

参考 `template-service` 目录结构，创建新的微服务：

1. **创建服务目录**: `mkdir microservice/your-service`
2. **复制模板文件**: 复制template-service的文件结构
3. **实现服务逻辑**: 
   - 修改 `service.go` 实现服务结构体
   - 修改 `handlers.go` 实现消息处理逻辑
   - 更新 `config.json` 配置文件
4. **注册服务**: 在 `main.go` 中注册服务处理器
5. **测试服务**: 启动服务并验证功能

### 消息通信

使用 `streamtool` 进行服务间通信：

```go
// 发送带响应的消息
response := service.SendMessageWithResponse(msg, "target-stream")
result := response.Get() // 阻塞获取响应

// 发送不带响应的消息
success := service.SendMessageWithoutResponse(msg, "notification-stream")
```

### 工具包使用

```go
// 使用HTTP客户端
httpClient := MyTool.NewHttpClient()
response, err := httpClient.Get("https://example.com")

// 使用JSON操作
data := MyTool.ParseJSON(jsonString)

// 使用Redis客户端
redisClient := MyTool.NewRedisClient(config)
err := redisClient.Set("key", "value")
```

## 配置说明

每个服务都有自己的 `config.json` 配置文件，主要包含：

### 基础配置

```json
{
  "redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
  },
  "consumer_group": "service-group",
  "input_stream": "service-input-stream",
  "debug": false
}
```

### 服务特定配置

- **cache-service**: 持久化配置
- **monitor-service**: 任务配置路径、健康检查端口
- **notifier-service**: 邮件配置
- **其他服务**: 根据功能添加特定配置

## 部署指南

### 开发环境

1. 安装 Go 1.23+
2. 安装 Redis 6.0+
3. 克隆代码库
4. 运行 `go mod tidy` 安装依赖
5. 配置 Redis 连接
6. 启动各个服务

### 生产环境

1. 使用 `build.bat` 构建所有服务
2. 配置生产环境的 `config.json`
3. 使用 systemd 或其他进程管理器管理服务
4. 设置日志记录和监控
5. 配置防火墙和安全设置

## 监控和维护

### 健康检查

监控服务提供健康检查接口：

```
http://localhost:8080/health
```

### 日志管理

每个服务都会输出日志，包含：
- 服务启动信息
- 消息处理日志
- 错误信息
- 性能指标

### 故障排除

#### 常见问题

1. **服务无法启动**
   - 检查Redis连接配置
   - 验证配置文件格式
   - 查看端口占用情况
   - 检查Go版本兼容性

2. **消息发送失败**
   - 检查目标流是否存在
   - 验证消息格式是否正确
   - 查看Redis连接状态
   - 检查消费者组配置

3. **响应获取超时**
   - 检查目标服务是否正常运行
   - 验证回调流配置
   - 查看网络连接
   - 检查消息处理逻辑

4. **Redis连接问题**
   - 检查Redis服务是否运行
   - 验证连接参数
   - 查看防火墙设置
   - 检查Redis密码配置

5. **内存占用过高**
   - 检查消息队列积压
   - 优化消息处理逻辑
   - 实现消息过期机制
   - 增加资源监控

## 性能优化

### 服务优化

1. **使用连接池**: 优化Redis和HTTP连接管理
2. **异步处理**: 实现异步消息处理
3. **批量操作**: 减少网络往返
4. **缓存策略**: 实现多级缓存
5. **负载均衡**: 支持服务水平扩展

### 架构优化

1. **服务拆分**: 合理拆分服务职责
2. **消息设计**: 优化消息结构
3. **错误处理**: 完善错误处理机制
4. **监控告警**: 实现全面监控
5. **容灾备份**: 配置数据备份策略

## 安全考虑

1. **密码管理**: 使用环境变量或加密存储密码
2. **访问控制**: 实现服务间认证
3. **输入验证**: 验证所有输入数据
4. **日志脱敏**: 敏感信息脱敏处理
5. **网络安全**: 配置防火墙和TLS

## 贡献指南

欢迎提交Issue和Pull Request来改进项目。

### 开发流程

1. Fork代码库
2. 创建功能分支
3. 实现功能或修复bug
4. 编写测试用例
5. 提交Pull Request
6. 代码审查
7. 合并代码

### 代码规范

- 遵循Go语言标准规范
- 使用Go Modules管理依赖
- 编写清晰的注释
- 实现完整的错误处理
- 编写单元测试

## 许可证

MIT License

## 联系方式

- GitHub: [WebSupervisor](https://github.com/yourusername/WebSupervisor)
- 邮箱: your@email.com

## 更新日志

### v1.0.0
- 初始版本发布
- 实现核心微服务架构
- 支持Redis Stream通信
- 提供完整的开发文档