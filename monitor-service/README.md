# Monitor Service 文档

## 概述

Monitor Service 是一个爬虫主控服务，负责任务调度和管理。该服务从配置文件加载监控任务，定期调度任务到各个微服务进行处理。

## 提供的功能

- **任务加载**: 从配置文件加载监控任务
- **任务调度**: 定期调度任务到爬虫服务
- **任务管理**: 支持多项目、多URL的监控配置

## 配置文件说明

### 主配置文件 (config.json)

```json
{
  "redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
  },
  "jobs_path": "./jobs.json",
  "interval_second": 60,
  "crawler_service_input_stream": "crawler-tasks"
}
```

### 任务配置文件 (jobs.json)

```json
{
  "projectName": "项目名称",
  "header": {
    "User-Agent": "Mozilla/5.0"
  },
  "urls": [
    {
      "url": "https://example.com",
      "output": "output.json",
      "test": false,
      "header": {
        "Custom-Header": "value"
      },
      "method": "GET",
      "body": {
        "key": "value"
      },
      "stringPlayLoad": "",
      "type": "html",
      "jsonKeys": [],
      "htmlKeys": [
        {
          "left": "<h1>",
          "right": "</h1>",
          "keys": ["title"]
        }
      ],
      "intervalSecond": 300
    }
  ],
  "intervalSecond": 60
}
```

## 运行方式

```bash
./monitor-service.exe -debug -config ./monitor-service/config.json
```

## 任务调度流程

1. 从配置文件加载监控任务
2. 按照配置的时间间隔定期调度任务
3. 将任务发送到 crawler-service 的输入流
4. 等待下一次调度周期

## URL 配置参数说明

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| url | string | 是 | 监控的URL地址 |
| output | string | 否 | 输出文件名 |
| test | boolean | 否 | 是否为测试模式 |
| header | object | 否 | 自定义HTTP头 |
| method | string | 否 | HTTP方法（GET/POST） |
| body | object | 否 | JSON格式的请求体 |
| stringPlayLoad | string | 否 | 字符串格式的请求体 |
| type | string | 否 | 内容类型（html/json） |
| jsonKeys | array | 否 | JSON提取规则 |
| htmlKeys | array | 否 | HTML提取规则 |
| intervalSecond | int | 否 | 单个URL的调度间隔 |

## 调度机制

- **全局间隔**: `interval_second` 参数控制整体调度频率
- **URL级别间隔**: 每个URL可以设置独立的 `intervalSecond` 参数
- **任务生成**: 每次调度为每个URL生成一个唯一的任务ID

## 任务消息格式

```json
{
  "task_id": "task-20260328123456",
  "consumer_group": "parser-group",
  "callback_stream": "parser-tasks",
  "service_name": "http_request",
  "parameter": {
    "url": "https://example.com",
    "method": "GET",
    "headers": {
      "User-Agent": "Mozilla/5.0"
    },
    "body": {},
    "str_payload": ""
  }
}
```

## 监控指标

- 当前加载的任务数量
- 调度频率
- 任务发送状态

## 错误处理

- 配置文件加载失败时会记录错误日志并退出
- 任务发送失败时会记录错误日志但继续运行