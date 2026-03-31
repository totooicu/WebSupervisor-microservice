import redis
import json
import time

redis_client = redis.Redis(host='localhost', port=6379, decode_responses=True)

print("=== 单独测试乘法功能 ===")

# 发送乘法请求到Python服务
print("发送乘法请求到Python服务...")
redis_client.xadd(
    "template-service-python-input-stream",
    {
        "message_id": "",
        "reply_id": "",
        "service_name": "multiply",
        "callback_stream": "test-response-stream",
        "playload": json.dumps({"num1": 8, "num2": 7})
    }
)

time.sleep(1)

# 读取响应
messages = redis_client.xread(
    streams={"test-response-stream": "0"},
    count=1,
    block=5000
)

if messages:
    for stream, stream_messages in messages:
        for msg_id, msg_data in stream_messages:
            playload = json.loads(msg_data.get('playload', '{}'))
            print(f"乘法结果: {json.dumps(playload, indent=2, ensure_ascii=False)}")
            redis_client.xdel("test-response-stream", msg_id)
else:
    print("未收到乘法响应")
    # 检查Python服务是否正常运行
    print("\n检查Python服务日志...")
