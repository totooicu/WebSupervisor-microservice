# Monitor Service API 文档

## 服务概述

Monitor Service 是一个定时任务调度服务，负责根据配置文件中的时间参数，周期性地将任务发送到各个微服务进行处理。

## 核心功能

1. **配置文件解析**：准确读取并验证 jobs 配置信息与邮箱配置参数
2. **服务通信**：将解析后的 jobs 信息分别发送至 crawler-service、parser-service 和 cache-service
3. **事件监听**：监听服务数据更新，自动触发 notifier-service 执行通知流程
4. **定时任务调度**：根据 jobs 配置中的时间参数，周期性执行服务通信与事件监听流程
5. **健康检查**：提供 HTTP 接口用于服务健康状态检查

## 配置文件结构

### jobs.json 配置示例

```json
{
  "projectName": "test",
  "header": {},
  "intervalSecond": 60,
  "urls": [
    {
      "url": "https://example.com/api/data",
      "output": "./output/data.json",
      "test": false,
      "header": {
        "User-Agent": "Monitor-Service/1.0"
      },
      "method": "GET",
      "body": {},
      "type": "json",
      "jsonKeys": [
        {
          "path": ["data", "items"],
          "key": ["id", "name", "value"]
        }
      ],
      "htmlKeys": []
    }
  ]
}
```

### config.json 配置示例

```json
{
  "redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
  },
  "consumer_group": "monitor-group",
  "input_stream": "monitor-tasks",
  "output_stream": "monitor-output",
  "crawler-service_input_stream": "crawler-tasks",
  "parser-service_input_stream": "parser-tasks",
  "cache-service_input_stream": "cache-tasks",
  "notifier-service_input_stream": "notifier-tasks",
  "jobs_path": "./jobs.json",
  "interval_second": 60,
  "health_check_port": "8080"
}
```

## API 接口

### 健康检查接口

**端点**: `/health`  
**方法**: GET  
**描述**: 检查服务健康状态  

**响应**:
- 成功 (200 OK):
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

- 失败 (503 Service Unavailable):
```json
{
  "status": "unhealthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "issues": "Redis connection failed: connection refused; Config files check failed: jobs.json not found"
}
```

## 服务通信流程

1. **任务调度**: Monitor Service 根据配置的时间间隔，周期性地调度任务
2. **发送到爬虫服务**: 将 URL 请求任务发送到 crawler-service
3. **爬虫处理**: crawler-service 执行 HTTP 请求，获取原始数据
4. **发送到解析服务**: 将原始数据发送到 parser-service 进行解析
5. **解析处理**: parser-service 根据配置提取所需数据
6. **发送到缓存服务**: 将解析后的数据发送到 cache-service 进行缓存
7. **缓存处理**: cache-service 将数据存储到 Redis
8. **发送通知**: 通过 notifier-service 发送操作结果通知

## 事件监听机制

Monitor Service 通过 Redis Stream 监听各个服务的回调事件：

1. **crawler_callback**: 爬虫服务完成后的回调
2. **parser_callback**: 解析服务完成后的回调  
3. **cache_callback**: 缓存服务完成后的回调

当收到回调事件时，服务会根据事件类型执行相应的处理逻辑，包括错误处理和通知发送。

## 错误处理

服务包含完整的错误处理机制：

1. **配置验证**: 在启动时验证所有配置参数的有效性
2. **连接错误**: 检测 Redis 连接状态，自动重连
3. **任务执行错误**: 捕获并记录任务执行过程中的所有错误
4. **通知机制**: 发生错误时通过 notifier-service 发送通知

## 日志记录

服务使用 Go 标准库的 log 包进行日志记录：

- **INFO**: 记录服务启动、任务执行等正常操作
- **DEBUG**: 在 debug 模式下记录详细的调试信息
- **ERROR**: 记录所有错误信息和异常情况

日志格式包含时间戳、日志级别和详细信息。

## 部署指南

### 前置条件

- Go 1.20+
- Redis 6.0+
- 其他微服务 (crawler-service, parser-service, cache-service, notifier-service)

### 部署步骤

1. **配置文件准备**:
   - 创建 `config.json` 配置文件
   - 创建 `jobs.json` 任务配置文件

2. **环境变量**:
   配置文件中支持环境变量替换，格式为 `${ENV_VAR_NAME}`

3. **启动服务**:
   ```bash
   go run main.go --config=config.json --debug=true
   ```

4. **健康检查**:
   ```bash
   curl http://localhost:8080/health
   ```

### Docker 部署

```dockerfile
FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go build -o monitor-service .

EXPOSE 8080

CMD ["./monitor-service", "--config=config.json"]
```

## 故障排除

### 常见问题

1. **Redis 连接失败**:
   - 检查 Redis 服务器是否运行
   - 验证配置中的主机名、端口和密码

2. **配置文件错误**:
   - 检查 JSON 格式是否正确
   - 验证必填字段是否存在

3. **服务通信失败**:
   - 检查其他微服务是否正常运行
   - 验证输入流名称是否正确

4. **健康检查失败**:
   - 检查 Redis 连接
   - 验证配置文件路径

## 版本控制策略

- **主分支**: 稳定版本，经过充分测试
- **开发分支**: 新功能开发和测试
- **标签**: 使用语义化版本号 (v1.0.0, v1.0.1)

## 贡献指南

1. Fork 仓库
2. 创建功能分支
3. 提交代码
4. 创建 Pull Request
5. 代码审查
6. 合并到主分支

## 许可证

MIT License

## 联系方式

- 项目地址: https://github.com/your-repo/WebSupervisor
- 问题反馈: https://github.com/your-repo/WebSupervisor/issues