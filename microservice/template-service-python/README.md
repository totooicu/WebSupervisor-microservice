# Template Service Python Version

Python版本的微服务模板，提供加减乘除功能，基于Redis Stream进行通信。

## 功能特性

- **加法运算**: 支持两个数字的加法运算
- **减法运算**: 支持两个数字的减法运算  
- **乘法运算**: 支持两个数字的乘法运算
- **除法运算**: 支持两个数字的除法运算，包含除零检查
- **Echo功能**: 返回相同的消息
- **基于Redis Stream**: 使用Redis Stream进行消息通信
- **手动服务注册**: 支持手动注册服务处理器
- **请求-响应模式**: 支持带响应的消息发送

## 目录结构

```
template-service-python/
├── main.py              # 应用入口文件
├── service.py           # 核心服务类
├── handlers.py          # 服务处理函数
├── models.py            # 数据模型定义
├── config.json          # 配置文件
├── requirements.txt     # Python依赖
├── run.bat              # 启动脚本
└── README.md            # 说明文档
```

## 技术栈

- **语言**: Python 3.8+
- **Redis客户端**: redis>=4.5.4
- **通信**: Redis Stream

## 快速开始

### 环境要求

- Python 3.8 或更高版本
- Redis 6.0 或更高版本

### 安装依赖

```bash
pip install -r requirements.txt
```

### 启动服务

```bash
# 使用启动脚本
./run.bat

# 或直接运行
python main.py --config ./config.json --debug
```

## API接口

### 服务名称

- `echo`: 返回相同的消息
- `add`: 加法运算
- `subtract`: 减法运算
- `multiply`: 乘法运算
- `divide`: 除法运算

### 请求参数

#### echo

```json
{
  "message": "Hello World"
}
```

#### add

```json
{
  "num1": 10,
  "num2": 20
}
```

#### subtract

```json
{
  "num1": 30,
  "num2": 15
}
```

#### multiply

```json
{
  "num1": 5,
  "num2": 6
}
```

#### divide

```json
{
  "num1": 20,
  "num2": 4
}
```

### 响应格式

```json
{
  "success": true,
  "data": {
    "result": 30,
    "num1": 10,
    "num2": 20
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
  "consumer_group": "template-service-python-group",
  "input_stream": "template-service-python-input-stream"
}
```

## 使用示例

### 发送加法请求

```python
import redis
import json
from models import StreamMessage

# 创建Redis客户端
redis_client = redis.Redis(host='localhost', port=6379, decode_responses=True)

# 创建消息
msg = StreamMessage(
    service_name="add",
    callback_stream="response-stream",
    playload={
        "num1": 10,
        "num2": 20
    }
)

# 发送消息
redis_client.xadd(
    "template-service-python-input-stream",
    {"message": msg.to_json()}
)

print("Message sent successfully")
```

### 接收响应

```python
# 从响应流获取消息
messages = redis_client.xread(
    streams={"response-stream": ">"},
    count=1,
    block=5000
)

if messages:
    for stream, stream_messages in messages:
        for msg_id, msg_data in stream_messages:
            message_json = msg_data.get('message', '{}')
            msg_dict = json.loads(message_json)
            print(f"Received response: {msg_dict}")
```

## 开发指南

### 添加新服务

1. **在models.py中定义参数结构**:

```python
class NewParameter:
    def __init__(self, param1: str, param2: int):
        self.param1 = param1
        self.param2 = param2
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'NewParameter':
        return cls(
            param1=data.get('param1', ''),
            param2=data.get('param2', 0)
        )
```

2. **在handlers.py中添加处理函数**:

```python
def handle_new_service(self, msg: StreamMessage):
    if self.service.debug:
        self.logger.info(f"Handling new_service message: {msg.message_id}")
    
    try:
        params = self.service.parse_parameters(msg, NewParameter)
        
        # 处理业务逻辑
        result = f"Processed: {params.param1}, {params.param2}"
        
        self.service.send_success_response(msg, {
            "result": result,
            "param1": params.param1,
            "param2": params.param2
        })
        
    except Exception as e:
        self.logger.error(f"Error handling new_service: {e}")
        self.service.send_error_response(msg, f"Invalid parameters: {str(e)}")
```

3. **在main.py中注册服务**:

```python
service.register_handler('new_service', handlers.handle_new_service)
```

## 故障排除

### 常见问题

1. **Redis连接失败**
   - 检查Redis服务是否运行
   - 验证连接配置是否正确
   - 查看网络连接

2. **消息发送失败**
   - 检查Redis连接状态
   - 验证消息格式是否正确
   - 查看目标流是否存在

3. **响应获取超时**
   - 检查服务是否正常运行
   - 验证回调流配置
   - 查看消息处理逻辑

4. **除零错误**
   - 检查除法运算的除数是否为零
   - 实现适当的错误处理

## 性能优化

- 使用连接池管理Redis连接
- 实现异步消息处理
- 优化消息序列化和反序列化
- 使用批量操作减少网络往返

## 安全注意事项

- 验证输入参数的有效性
- 处理异常情况，避免服务崩溃
- 实现适当的日志记录
- 限制消息大小，防止内存溢出

## 许可证

MIT License