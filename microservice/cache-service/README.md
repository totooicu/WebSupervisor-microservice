# Cache Service 微服务文档

## 1. 服务概述

Cache Service 是一个高性能的缓存管理微服务，提供统一的缓存操作接口，支持数据的增删改查、数据比较与保存等功能。该服务基于Redis实现，提供了数据持久化和数据变化通知机制。

## 2. 核心功能

- **数据存储**：将数据存储到Redis缓存中
- **数据读取**：从Redis缓存中读取数据
- **数据删除**：从Redis缓存中删除数据
- **数据比较与保存**：比较新旧数据，只有在数据发生变化时才保存并发送通知
- **原子操作**：支持获取并设置的原子操作

## 3. 配置说明

服务配置文件 `config.json` 格式如下：

```json
{
  "redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
  },
  "consumer_group": "cache-group",
  "input_stream": "cache-tasks",
  "persistence": {
    "app_paths": {
      "default": "./data"
    }
  }
}
```

## 4. 接口说明

### 4.1 请求消息格式

所有请求消息都遵循统一格式：

```json
{
  "task_id": "task-1234567890",
  "consumer_group": "notifier-group",
  "callback_stream": "responses-tasks",
  "service_name": "compare_and_save",
  "playload": "{\"app\": \"myapp\", \"key\": \"user:123\", \"data\": {...}}",
  "reply_id": "-1"
}
```

### 4.2 响应消息格式

所有响应消息都遵循统一格式：

```json
{
  "callback_stream": "cache-tasks",
  "consumer_group": "responses-tasks",
  "playload": "{\"changed\": true}",
  "service_name": "response",
  "message_id": "-1",
  "reply_id": "task-1234567890"
}
```

## 5. 参数说明

### CacheParameter 参数结构

```go
type CacheParameter struct {
    App  string      `json:"app"`
    Key  string      `json:"key"`
    Data interface{} `json:"data"`
}
```

| 参数名 | 数据类型 | 是否必填 | 默认值 | 描述 |
|--------|----------|----------|--------|------|
| app | string | 是 | 无 | 应用名称，用于生成唯一的缓存键 |
| key | string | 是 | 无 | 缓存键名称 |
| data | interface{} | 否 | 无 | 要缓存的数据内容 |

### 支持的服务名称

| 服务名称 | 描述 | 参数要求 |
|----------|------|----------|
| compare_and_save | 比较并保存数据，数据变化时发送通知 | app, key, data |
| get | 获取缓存数据 | app, key |
| set | 设置缓存数据 | app, key, data |
| delete | 删除缓存数据 | app, key |
| get_and_set | 获取并设置缓存数据（原子操作） | app, key, data |

## 6. 返回参数说明

### compare_and_save 响应

```json
{
  "changed": true
}
```

| 字段名 | 数据类型 | 描述 |
|--------|----------|------|
| changed | bool | 数据是否发生变化 |

### get 响应

```json
{
  "key": "myapp:user:123",
  "data": {...},
  "error": false
}
```

| 字段名 | 数据类型 | 描述 |
|--------|----------|------|
| key | string | 生成的缓存键 |
| data | interface{} | 获取到的数据 |
| error | bool | 是否发生错误 |

### set 响应

```json
{
  "key": "myapp:user:123",
  "error": null
}
```

| 字段名 | 数据类型 | 描述 |
|--------|----------|------|
| key | string | 生成的缓存键 |
| error | null/error | 错误信息，成功时为null |

### delete 响应

```json
{
  "key": "myapp:user:123",
  "error": null
}
```

| 字段名 | 数据类型 | 描述 |
|--------|----------|------|
| key | string | 生成的缓存键 |
| error | null/error | 错误信息，成功时为null |

### get_and_set 响应

```json
{
  "key": "myapp:user:123",
  "old_data": {...},
  "error": null
}
```

| 字段名 | 数据类型 | 描述 |
|--------|----------|------|
| key | string | 生成的缓存键 |
| old_data | interface{} | 更新前的数据 |
| error | null/error | 错误信息，成功时为null |

## 7. 错误码说明

| 错误码 | 描述 | 可能原因 |
|--------|------|----------|
| JSON解析错误 | Error unmarshalling playload | 请求数据格式不正确 |
| Redis操作错误 | Error saving data | Redis连接问题或权限不足 |
| 数据不存在 | Key not found | 请求的缓存键不存在 |

## 8. 调用示例

### 示例1：比较并保存数据

```bash
# 使用streams-manager工具发送消息
streams-manager add cache-tasks -f compare_and_save.json
```

`compare_and_save.json` 内容：
```json
{
  "callback_stream": "responses-tasks",
  "consumer_group": "notifier-group",
  "message_id": "1",
  "playload": "{\"app\": \"myapp\", \"key\": \"user:123\", \"data\": {\"name\": \"张三\", \"age\": 30}}",
  "reply_id": "-1",
  "service_name": "compare_and_save"
}
```

### 示例2：获取缓存数据

```bash
streams-manager add cache-tasks -f get.json
```

`get.json` 内容：
```json
{
  "callback_stream": "responses-tasks",
  "consumer_group": "notifier-group",
  "message_id": "2",
  "playload": "{\"app\": \"myapp\", \"key\": \"user:123\"}",
  "reply_id": "-1",
  "service_name": "get"
}
```

### 示例3：设置缓存数据

```bash
streams-manager add cache-tasks -f set.json
```

`set.json` 内容：
```json
{
  "callback_stream": "responses-tasks",
  "consumer_group": "notifier-group",
  "message_id": "3",
  "playload": "{\"app\": \"myapp\", \"key\": \"user:123\", \"data\": {\"name\": \"李四\", \"age\": 25}}",
  "reply_id": "-1",
  "service_name": "set"
}
```

### 示例4：删除缓存数据

```bash
streams-manager add cache-tasks -f delete.json
```

`delete.json` 内容：
```json
{
  "callback_stream": "responses-tasks",
  "consumer_group": "notifier-group",
  "message_id": "4",
  "playload": "{\"app\": \"myapp\", \"key\": \"user:123\"}",
  "reply_id": "-1",
  "service_name": "delete"
}
```

## 9. 使用注意事项

1. **数据格式**：所有请求和响应数据都必须是有效的JSON格式
2. **键名规范**：缓存键由 `app:key` 格式生成，确保app和key的唯一性
3. **数据类型**：支持任意JSON可序列化的数据类型
4. **错误处理**：服务会自动处理Redis连接错误和JSON解析错误
5. **性能考虑**：对于大量数据操作，建议分批处理以避免Redis性能问题
6. **持久化**：服务支持数据持久化配置，可在config.json中设置
7. **并发安全**：Redis操作是原子的，但业务逻辑需要考虑并发场景
8. **监控**：建议监控Redis内存使用情况，避免内存溢出
