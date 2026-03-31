# Cache Service

缓存服务，提供数据缓存功能，支持数据的存储、获取、删除等操作。

## 功能特性

- **数据存储**: 支持将数据存储到Redis中
- **数据获取**: 支持从Redis中获取数据
- **数据删除**: 支持删除Redis中的数据
- **数据比较**: 支持比较新旧数据并保存
- **原子操作**: 支持GET和SET的原子操作

## API接口

### 服务名称

- `compare_and_save`: 比较并保存数据
- `get`: 获取数据
- `set`: 设置数据
- `delete`: 删除数据
- `get_and_set`: 获取并设置数据（原子操作）

### 请求参数

#### compare_and_save
```json
{
  "app": "app_name",
  "key": "data_key",
  "data": "data_value"
}
```

#### get
```json
{
  "app": "app_name",
  "key": "data_key"
}
```

#### set
```json
{
  "app": "app_name",
  "key": "data_key",
  "data": "data_value"
}
```

#### delete
```json
{
  "app": "app_name",
  "key": "data_key"
}
```

#### get_and_set
```json
{
  "app": "app_name",
  "key": "data_key",
  "data": "data_value"
}
```

### 响应格式

```json
{
  "success": true,
  "data": {
    "result": "operation_result",
    "data": "returned_data"
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
  "consumer_group": "cache-service-group",
  "input_stream": "cache-service-input-stream",
  "persistence": {
    "app_paths": {
      "app1": "/path/to/app1",
      "app2": "/path/to/app2"
    }
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

## 数据存储格式

数据存储在Redis中的格式为：
```
[app]:[key] -> data
```

## 开发说明

### 扩展功能

1. 修改 `handlers.go` 文件，添加新的处理函数
2. 在 `service.go` 中注册新的服务
3. 更新配置文件添加必要的配置项

### 最佳实践

- 使用统一的应用名称前缀
- 合理设置数据过期时间
- 实现数据备份机制
- 添加数据验证逻辑

## 故障排除

### 常见问题

1. **Redis连接失败**
   - 检查Redis服务是否运行
   - 验证连接配置是否正确

2. **数据存储失败**
   - 检查Redis内存是否充足
   - 验证数据格式是否正确

3. **服务注册失败**
   - 检查服务名称是否重复
   - 验证输入流配置

## 性能优化

- 使用Redis集群提高性能
- 实现数据压缩减少存储空间
- 添加缓存预热机制
- 优化数据序列化方式