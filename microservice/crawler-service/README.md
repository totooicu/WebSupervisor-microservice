# Crawler Service

爬虫服务，用于爬取网页内容，支持GET和POST请求。

## 功能特性

- **HTTP请求**: 支持GET和POST请求
- **自定义Headers**: 支持自定义请求头
- **请求参数**: 支持URL查询参数和请求体
- **响应解析**: 返回响应内容和状态码
- **错误处理**: 完善的错误处理机制

## API接口

### 服务名称

- `http_request`: 发送HTTP请求

### 请求参数

```json
{
  "url": "https://example.com",
  "method": "GET",
  "headers": {
    "Content-Type": "application/json",
    "User-Agent": "WebSupervisor-Crawler"
  },
  "body": {
    "key1": "value1",
    "key2": "value2"
  },
  "str_payload": "{\"key\": \"value\"}"
}
```

### 参数说明

- `url`: 请求URL（必填）
- `method`: 请求方法，支持GET、POST（必填）
- `headers`: 请求头，键值对形式
- `body`: POST请求的JSON数据
- `str_payload`: 字符串形式的请求体，优先级高于body

### 响应格式

```json
{
  "success": true,
  "data": {
    "content": "response_content",
    "status": 200
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
  "consumer_group": "crawler-service-group",
  "input_stream": "crawler-service-input-stream"
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

### GET请求示例

```json
{
  "service": "http_request",
  "playload": {
    "url": "https://api.example.com/data",
    "method": "GET",
    "headers": {
      "Authorization": "Bearer token123"
    }
  }
}
```

### POST请求示例

```json
{
  "service": "http_request",
  "playload": {
    "url": "https://api.example.com/data",
    "method": "POST",
    "headers": {
      "Content-Type": "application/json"
    },
    "body": {
      "username": "admin",
      "password": "password123"
    }
  }
}
```

## 开发说明

### 扩展功能

1. 修改 `handlers.go` 文件，添加新的HTTP方法支持
2. 在 `service.go` 中注册新的服务
3. 更新配置文件添加必要的配置项

### 最佳实践

- 设置合理的请求超时时间
- 实现请求重试机制
- 添加请求频率限制
- 使用代理IP避免被封禁
- 设置合适的User-Agent

## 故障排除

### 常见问题

1. **请求失败**
   - 检查URL是否正确
   - 验证网络连接
   - 查看目标网站是否可访问

2. **响应为空**
   - 检查状态码是否为200
   - 验证请求头是否正确
   - 查看目标网站是否有反爬虫机制

3. **服务无法启动**
   - 检查Redis连接
   - 验证配置文件格式

## 性能优化

- 使用连接池管理HTTP连接
- 实现请求并发控制
- 添加缓存机制避免重复请求
- 使用异步处理提高吞吐量

## 安全注意事项

- 避免爬取敏感网站
- 遵守网站的robots.txt规则
- 设置合理的请求间隔
- 不要爬取需要登录的内容
- 尊重网站的使用条款