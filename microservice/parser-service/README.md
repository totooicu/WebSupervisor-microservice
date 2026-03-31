# Parser Service

解析服务，提供HTML和JSON解析功能，支持提取指定内容。

## 功能特性

- **HTML解析**: 支持从HTML内容中提取指定内容
- **JSON解析**: 支持从JSON内容中提取指定字段
- **多规则解析**: 支持同时应用多个解析规则
- **灵活配置**: 支持配置解析规则
- **错误处理**: 完善的错误处理机制

## API接口

### 服务名称

- `parse_html`: 解析HTML内容
- `parse_json`: 解析JSON内容

### 请求参数

#### parse_html
```json
{
  "content": "<html><body><title>Test</title></body></html>",
  "htmlKeys": [
    {
      "left": "<title>",
      "right": "</title>",
      "keys": ["title"]
    }
  ]
}
```

#### parse_json
```json
{
  "content": "{\"name\": \"John\", \"age\": 30}",
  "jsonKeys": [
    {
      "path": ["name"],
      "keys": ["name"]
    }
  ]
}
```

### 参数说明

#### HTML解析
- `content`: HTML内容（必填）
- `htmlKeys`: HTML解析规则数组
  - `left`: 左边界字符串
  - `right`: 右边界字符串
  - `keys`: 提取的字段名称

#### JSON解析
- `content`: JSON内容（必填）
- `jsonKeys`: JSON解析规则数组
  - `path`: JSON路径数组
  - `keys`: 提取的字段名称

### 响应格式

```json
{
  "success": true,
  "data": {
    "results": [
      {
        "title": "Test"
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
  "consumer_group": "parser-service-group",
  "input_stream": "parser-service-input-stream"
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

### HTML解析示例

```json
{
  "service": "parse_html",
  "playload": {
    "content": "<div><h1>Hello World</h1><p>Welcome to WebSupervisor</p></div>",
    "htmlKeys": [
      {
        "left": "<h1>",
        "right": "</h1>",
        "keys": ["title"]
      },
      {
        "left": "<p>",
        "right": "</p>",
        "keys": ["content"]
      }
    ]
  }
}
```

### JSON解析示例

```json
{
  "service": "parse_json",
  "playload": {
    "content": "{\"user\": {\"name\": \"John\", \"age\": 30, \"address\": {\"city\": \"New York\"}}}",
    "jsonKeys": [
      {
        "path": ["user", "name"],
        "keys": ["username"]
      },
      {
        "path": ["user", "age"],
        "keys": ["userage"]
      },
      {
        "path": ["user", "address", "city"],
        "keys": ["city"]
      }
    ]
  }
}
```

## 开发说明

### 扩展功能

1. 修改 `handlers.go` 文件，添加新的解析器
2. 在 `service.go` 中注册新的服务
3. 更新配置文件添加必要的配置项

### 最佳实践

- 使用精确的解析规则
- 处理解析失败的情况
- 添加解析结果验证
- 优化解析性能
- 支持复杂的解析规则

## 故障排除

### 常见问题

1. **解析结果为空**
   - 检查解析规则是否正确
   - 验证内容格式是否符合预期
   - 查看解析日志

2. **解析错误**
   - 检查内容格式是否正确
   - 验证解析规则的语法
   - 查看错误日志

3. **服务无法启动**
   - 检查Redis连接
   - 验证配置文件格式

## 性能优化

- 使用高效的解析算法
- 实现解析结果缓存
- 优化正则表达式
- 使用并发解析提高效率
- 添加解析超时控制

## 安全注意事项

- 验证输入内容的安全性
- 防止注入攻击
- 限制解析内容的大小
- 实现解析频率限制
- 监控异常解析请求