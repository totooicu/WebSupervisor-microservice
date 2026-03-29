# Parser Service 详细说明文档

## 服务概述

Parser Service 是一个数据解析微服务，提供 HTML 和 JSON 格式数据的解析功能。该服务通过 Redis Streams 接收解析任务，支持使用左右边界匹配解析 HTML 内容，以及使用路径表达式解析 JSON 数据。服务采用异步处理模式，确保高效可靠的数据解析能力。

## 核心功能说明

### 1. HTML 解析功能
- **左右边界匹配**：通过指定左右边界字符串提取 HTML 内容
- **多结果支持**：支持提取多个匹配结果
- **灵活配置**：支持配置多个解析规则

### 2. JSON 解析功能
- **路径表达式**：支持使用点号分隔的路径表达式解析 JSON 数据
- **多字段提取**：支持同时提取多个 JSON 字段
- **类型转换**：自动处理不同数据类型的提取

### 3. 异步处理机制
- 通过 Redis Streams 实现任务的异步处理
- 支持批量解析请求
- 提供详细的调试日志

## 参数说明

### HTML 解析参数

| 参数名 | 数据类型 | 是否必填 | 默认值 | 详细描述 |
|--------|----------|----------|--------|----------|
| content | string | 是 | 无 | 需要解析的 HTML 内容 |
| HTMLKeys | array[HTMLKey] | 是 | 无 | HTML 解析规则数组 |

#### HTMLKey 结构

| 字段名 | 数据类型 | 是否必填 | 默认值 | 详细描述 |
|--------|----------|----------|--------|----------|
| left | string | 是 | 无 | 左边界字符串，用于定位开始位置 |
| right | string | 是 | 无 | 右边界字符串，用于定位结束位置 |
| keys | array[string] | 否 | 空数组 | 提取结果的键名列表 |

### JSON 解析参数

| 参数名 | 数据类型 | 是否必填 | 默认值 | 详细描述 |
|--------|----------|----------|--------|----------|
| content | string | 是 | 无 | 需要解析的 JSON 字符串 |
| JSONKeys | array[JSONKey] | 是 | 无 | JSON 解析规则数组 |

#### JSONKey 结构

| 字段名 | 数据类型 | 是否必填 | 默认值 | 详细描述 |
|--------|----------|----------|--------|----------|
| path | array[interface{}] | 是 | 无 | JSON 路径数组，支持字符串和数字索引 |
| keys | array[string] | 否 | 空数组 | 提取结果的键名列表 |

### 数据结构定义

```go
type ParserParameter struct {
    Content  string      `json:"content"`
    HTMLKeys []HTMLKey   `json:"HTMLKeys"`
    JSONKeys []JSONKey   `json:"JSONKeys"`
}

type HTMLKey struct {
    Left  string   `json:"left"`
    Right string   `json:"right"`
    Keys  []string `json:"keys"`
}

type JSONKey struct {
    Path []interface{} `json:"path"`
    Keys []string      `json:"keys"`
}
```

## 返回参数说明

### HTML 解析成功响应

| 字段名 | 数据类型 | 说明 |
|--------|----------|------|
| parsed_data | array[string] | 解析提取的字符串数组，包含所有匹配的结果 |

### JSON 解析成功响应

| 字段名 | 数据类型 | 说明 |
|--------|----------|------|
| parsed_data | map[string]interface{} | 解析提取的数据，键为路径，值为提取的数据 |

### 返回数据结构

```json
{
"parsed_data": ["提取结果1", "提取结果2"]
}
```

或者

```json
{
"parsed_data": {
    "user.name": "张三",
    "user.age": 25
}
}
```

## 错误码说明

| 错误码 | 错误描述 | 可能原因 |
|--------|----------|----------|
| 400 | 参数解析失败 | 请求参数格式错误或缺失必填字段 |
| 400 | HTMLKeys 为空 | HTML 解析时未提供解析规则 |
| 400 | JSONKeys 为空 | JSON 解析时未提供解析规则 |
| 500 | 解析失败 | 内容格式错误或解析规则不正确 |

## 调用示例

### 示例 1: HTML 解析

```json
{
"task_id": "html-123",
"consumer_group": "monitor-group",
"callback_stream": "monitor-responses",
"service_name": "parse_html",
"playload": "{\"content\":\"<div class=\\\"title\\\">标题1</div><div class=\\\"title\\\">标题2</div>\",\"HTMLKeys\":[{\"left\":\"<div class=\\\"title\\\">\",\"right\":\"</div>\",\"keys\":[\"titles\"]}]}"
}
```

### 示例 2: JSON 解析

```json
{
"task_id": "json-456",
"consumer_group": "monitor-group",
"callback_stream": "monitor-responses",
"service_name": "parse_json",
"playload": "{\"content\":\"{\\\"user\\\":{\\\"name\\\":\\\"张三\\\",\\\"age\\\":25,\\\"address\\\":{\\\"city\\\":\\\"北京\\\"}}}\",\"JSONKeys\":[{\"path\":[\"user\",\"name\"],\"keys\":[\"username\"]},{\"path\":[\"user\",\"age\"],\"keys\":[\"userage\"]}]}"
}
```

### 示例 3: 复杂 JSON 解析（包含数组索引）

```json
{
"task_id": "json-789",
"consumer_group": "monitor-group",
"callback_stream": "monitor-responses",
"service_name": "parse_json",
"playload": "{\"content\":\"{\\\"users\\\":[{\\\"name\\\":\\\"张三\\\"},{\\\"name\\\":\\\"李四\\\"}]}\",\"JSONKeys\":[{\"path\":[\"users\",0,\"name\"],\"keys\":[\"first_user\"]}]}"
}
```

## 使用注意事项

### 1. HTML 解析注意事项
- `left` 和 `right` 边界字符串必须精确匹配目标内容
- 支持提取多个匹配结果，按出现顺序返回
- 边界字符串中的特殊字符需要正确转义

### 2. JSON 解析注意事项
- `path` 数组支持字符串键和数字索引
- 数字索引用于访问数组元素
- 路径不存在时返回空值，不会报错

### 3. 配置要求
- Redis 连接配置必须正确，包括主机、端口、密码和数据库索引
- 消费者组和输入流名称必须与配置文件一致

### 4. 性能建议
- 对于大型内容，考虑分块处理以减少内存占用
- 避免使用过于宽泛的解析规则，可能导致匹配过多结果
- 对于频繁解析的内容，考虑使用缓存机制

### 5. 错误处理
- 解析失败时会记录详细的错误日志
- 建议在生产环境中启用调试模式以便排查问题
- 对输入内容进行验证，确保格式正确

## 配置文件说明

```json
{
"redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
},
"consumer_group": "parser-group",
"input_stream": "parser-tasks"
}
```

## 运行方式

```bash
# 开发模式（带调试日志）
./parser-service.exe -debug -config ./parser-service/config.json

# 生产模式
./parser-service.exe -config ./parser-service/config.json
```