# WebSupervisor 微服务化架构设计方案

## 一、整体架构

将现有单体应用拆分为5个独立的微服务模块，通过Redis Streams进行异步通信。每个模块可独立启动、部署和替换。

## 二、模块分解结构

### 1. 网络爬取模块 (crawler-service)
**核心功能：** 
- HTTP请求发送（GET/POST）
- 请求头管理
- 响应内容获取

**包含函数：** 
- `CrawlerService.Start()` - 启动爬取服务
- `CrawlerService.ProcessTask(task)` - 处理爬取任务
- `CrawlerService.SendResult(result)` - 发送爬取结果到Redis Streams

**配置文件 (config.json)：** 
```json
{
"redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
},
"consumer_group": "crawler-group",
"input_stream": "crawler-tasks"
}
```

**通信格式设计：**
```json
{
"task_id": "task-123456",
"consumer_group": "parser-group",
"callback_stream": "parser-tasks",
"service_name": "http_request",
"parameter": {
    "url": "https://example.com",
    "method": "GET",
    "headers": {},
    "body": {},
    "str_payload": "string payload for non-json data"
}
}
```

### 2. 文本解析模块 (parser-service)
**核心功能：** 
- JSON格式解析（支持路径提取）
- HTML格式解析（支持左右边界匹配）
- 正则表达式匹配

**包含函数：** 
- `ParserService.Start()` - 启动解析服务
- `ParserService.ParseJSON(content, jsonKeys)` - 解析JSON内容
- `ParserService.ParseHTML(content, htmlKeys)` - 解析HTML内容
- `ParserService.SendResult(result)` - 发送解析结果到Redis Streams

**配置文件 (config.json)：** 
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

**通信格式设计：**
```json
{
"task_id": "task-123456",
"consumer_group": "cache-group",
"callback_stream": "cache-tasks",
"service_name": "parse_html",
"parameter": {
    "content": "html content",
    "keys": [],
    "json_keys": []
}
}
```

### 3. 文本缓存模块 (cache-service)
**核心功能：** 
- Key-Value存储
- 数据对比（检测变化）
- 缓存管理
- 持久化管理
- 提供get、set、delete服务

**包含函数：** 
- `CacheService.Start()` - 启动缓存服务
- `CacheService.CompareData(key, newData)` - 比较新旧数据
- `CacheService.SaveData(key, data)` - 保存数据到缓存
- `CacheService.SendChangeEvent(event)` - 发送变化事件到Redis Streams
- `CacheService.Get(key)` - 获取缓存数据
- `CacheService.Set(key, data)` - 设置缓存数据
- `CacheService.Delete(key)` - 删除缓存数据

**配置文件 (config.json)：** 
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
    "app_name": "path"
    }
}
}
```

**通信格式设计：**
```json
{
"task_id": "task-123456",
"consumer_group": "notifier-group",
"callback_stream": "notifier-tasks",
"service_name": "compare_and_save",
"parameter": {
    "app": "app_name",
    "key": "key-123",
    "data": "data"
}
}
```

### 4. 邮箱通知模块 (notifier-service)
**核心功能：** 
- 邮件发送(多人)
- 消息格式化

**包含函数：** 
- `NotifierService.Start()` - 启动通知服务
- `NotifierService.SendEmail(to, subject, content)` - 发送邮件通知
- `NotifierService.ProcessNotification(notification)` - 处理通知任务

**配置文件 (config.json)：** 
```json
{
"redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
},
"consumer_group": "notifier-group",
"input_stream": "notifier-tasks",
"mail": {
    "userName": "1134815016@qq.com",
    "password": "pwhpkqixofabhaag",
    "tos": ["2667657661@qq.com"]
}
}
```

**通信格式设计：**
```json
{
"task_id": "task-123456",
"consumer_group": "notifier-group",
"callback_stream": "notifier-results",
"service_name": "send_email",
"parameter": {
    "url": "https://example.com",
    "subject": "网页监控匹配成功",
    "content": "匹配到的内容...",
    "userName": "1134815016@qq.com",
    "password": "pwhpkqixofabhaag",
    "tos": ["2667657661@qq.com"]
}
}
```

### 5. 爬虫主控模块 (monitor-service)
**核心功能：** 
- 定时任务调度
- 任务配置管理
- 系统协调

**包含函数：** 
- `MonitorService.Start()` - 启动监控服务
- `MonitorService.LoadJobs(configPath)` - 加载作业配置
- `MonitorService.ScheduleTask(task)` - 调度任务到爬取模块
- `MonitorService.SendTask(task)` - 发送任务到Redis Streams

**配置文件 (config.json)：** 
```json
{
"redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
},
"consumer_group": "monitor-group",
"output_stream": "crawler-tasks",
"crawler-service_input_stream": "crawler-tasks",
"parser-service_input_stream": "parser-tasks",
"cache-service_input_stream": "cache-tasks",
"notifier-service_input_stream": "notifier-tasks",
"jobs_path": "./jobs.json",
"interval_second": 20
}
```

**通信格式设计：**
负责统筹调度其他模块，通信格式具体根据模块的通信格式进行定义。必要格式为：
```json
{
"task_id": "task-123456",
"consumer_group": "notifier-group",
"callback_stream": "notifier-results",
"service_name": "",
"parameter": {}
}
```


## 三、配置管理

### 共享配置结构
每个模块的配置文件都包含：
- Redis连接配置（host, port, password, db）
- 消费者组名（consumer_group）
- 输入输出流配置
- 模块特定配置

### 作业配置文件 (jobs.json)
保持现有格式，由监控模块负责加载和分发。

## 四、部署与运行

每个模块独立运行：
1. 启动Redis服务
2. 分别启动5个微服务模块
3. 监控模块加载作业配置并开始调度

## 五、优势

1. **可独立扩展**：每个模块可根据负载单独扩展
2. **高可用性**：单个模块故障不影响整体系统
3. **灵活替换**：可随时替换或升级单个模块
4. **异步处理**：通过Redis Streams实现可靠的消息传递
5. **便于监控**：每个模块可单独监控和日志记录

这个设计保持了原有功能的完整性，同时实现了微服务化架构，满足您的所有要求。