# Monitor Service

监控服务，负责监控其他服务的健康状态，调度任务执行。

## 功能特性

- **服务监控**: 监控其他微服务的健康状态
- **任务调度**: 定时执行爬取、解析等任务
- **健康检查**: 提供服务健康状态检查接口
- **配置解析**: 解析和验证任务配置
- **服务通信**: 与其他服务进行通信

## API接口

### 服务名称

- `health_check`: 健康检查
- `task_schedule`: 任务调度
- `service_monitor`: 服务监控

### 请求参数

#### health_check
```json
{
  "service_name": "service_name"
}
```

#### task_schedule
```json
{
  "job_id": "job_123",
  "interval": 60,
  "crawler_params": {
    "url": "https://example.com",
    "method": "GET"
  }
}
```

### 响应格式

```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "services": [
      {
        "name": "crawler-service",
        "status": "running"
      }
    ]
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
  "consumer_group": "monitor-service-group",
  "input_stream": "monitor-service-input-stream",
  "output_stream": "monitor-service-output-stream",
  "crawler-service_input_stream": "crawler-service-input-stream",
  "parser-service_input_stream": "parser-service-input-stream",
  "cache-service_input_stream": "cache-service-input-stream",
  "notifier-service_input_stream": "notifier-service-input-stream",
  "jobs_path": "./jobs.json",
  "interval_second": 60,
  "health_check_port": "8080"
}
```

### jobs.json

```json
{
  "jobs": [
    {
      "id": "job_001",
      "name": "Example Job",
      "interval": 300,
      "crawler_params": {
        "url": "https://api.example.com/data",
        "method": "GET",
        "headers": {
          "Content-Type": "application/json"
        }
      },
      "parser_params": {
        "html_keys": [
          {
            "left": "<title>",
            "right": "</title>"
          }
        ]
      },
      "cache_params": {
        "app": "example",
        "key": "data"
      },
      "notifier_params": {
        "subject": "Job Notification",
        "content": "Job completed successfully"
      }
    }
  ]
}
```

## 运行服务

```bash
# 启动服务
go run main.go service.go service_communicator.go task_scheduler.go config_parser.go health_check.go utils.go --config ./config.json --debug

# 或使用批处理脚本
./run.bat
```

## 健康检查

监控服务提供健康检查接口，可以通过HTTP访问：

```
http://localhost:8080/health
```

## 开发说明

### 扩展功能

1. 修改 `task_scheduler.go` 文件，添加新的任务类型
2. 在 `service_communicator.go` 中添加新的服务通信逻辑
3. 更新配置文件添加必要的配置项

### 最佳实践

- 合理设置任务执行间隔
- 实现任务失败重试机制
- 添加任务执行日志
- 监控任务执行状态
- 实现任务优先级管理

## 故障排除

### 常见问题

1. **任务执行失败**
   - 检查任务配置是否正确
   - 验证目标服务是否正常运行
   - 查看任务执行日志

2. **服务通信失败**
   - 检查Redis连接
   - 验证流名称配置
   - 查看服务注册状态

3. **健康检查失败**
   - 检查服务端口配置
   - 验证HTTP服务是否启动
   - 查看网络连接

## 性能优化

- 使用并发执行提高任务处理效率
- 实现任务队列管理
- 添加任务执行统计
- 优化配置解析性能
- 使用缓存减少重复操作

## 安全注意事项

- 限制健康检查接口的访问权限
- 验证任务配置的合法性
- 防止恶意任务注入
- 监控异常任务执行
- 实现任务执行超时控制