import redis
import json

redis_client = redis.Redis(host='localhost', port=6379, decode_responses=True)

# 清空测试流
redis_client.delete("test-response-stream")

# 发送一个测试消息到Go服务
print("发送测试消息...")
msg = {
    "message_id": "",
    "reply_id": "", 
    "service_name": "add",
    "callback_stream": "test-response-stream",
    "playload": {"num1": 5, "num2": 3}
}
redis_client.xadd("template-service-input-stream", {"message": json.dumps(msg)})

print("等待响应...")
import time
time.sleep(1)

# 查看流中的所有消息
print("\n查看test-response-stream中的所有消息:")
messages = redis_client.xread(streams={"test-response-stream": "0"}, count=10, block=1000)
print(f"找到 {len(messages)} 条消息")

if messages:
    for stream, stream_messages in messages:
        print(f"\nStream: {stream}")
        for msg_id, msg_data in stream_messages:
            print(f"  Message ID: {msg_id}")
            print(f"  Message data: {msg_data}")
            for key, value in msg_data.items():
                print(f"    {key}: {value}")
                
                # 尝试解析message字段
                if key == 'message':
                    try:
                        msg_dict = json.loads(value)
                        print(f"    Parsed message: {json.dumps(msg_dict, indent=4)}")
                    except json.JSONDecodeError:
                        print(f"    Failed to parse JSON: {value}")
