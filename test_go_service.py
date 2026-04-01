import redis
import json
import time

# 创建Redis客户端
redis_client = redis.Redis(host='localhost', port=6379, decode_responses=True)

def test_go_service():
    print("=== 测试Go版本template-service ===")
    print("=" * 50)
    
    # 清空测试流
    redis_client.delete("test-response-stream")
    print("测试流已清空")
    print()
    
    # 测试加法
    print("1. 测试加法功能 (10 + 20)")
    redis_client.xadd(
        "template-service-input-stream",
        {
            "message_id": "",
            "reply_id": "",
            "service_name": "add",
            "callback_stream": "test-response-stream",
            "playload": json.dumps({"num1": 10, "num2": 20})
        }
    )
    
    time.sleep(0.5)
    
    messages = redis_client.xread(
        streams={"test-response-stream": "0"},
        count=1,
        block=3000
    )
    
    if messages:
        for stream, stream_messages in messages:
            for msg_id, msg_data in stream_messages:
                playload = json.loads(msg_data.get('playload', '{}'))
                print(f"   加法结果: {json.dumps(playload, indent=4)}")
                redis_client.xdel("test-response-stream", msg_id)
    else:
        print("   ❌ 未收到加法响应")
    
    print()
    
    # 测试减法
    print("2. 测试减法功能 (50 - 20)")
    redis_client.xadd(
        "template-service-input-stream",
        {
            "message_id": "",
            "reply_id": "",
            "service_name": "subtract",
            "callback_stream": "test-response-stream",
            "playload": json.dumps({"num1": 50, "num2": 20})
        }
    )
    
    time.sleep(0.5)
    
    messages = redis_client.xread(
        streams={"test-response-stream": "0"},
        count=1,
        block=3000
    )
    
    if messages:
        for stream, stream_messages in messages:
            for msg_id, msg_data in stream_messages:
                playload = json.loads(msg_data.get('playload', '{}'))
                print(f"   减法结果: {json.dumps(playload, indent=4)}")
                redis_client.xdel("test-response-stream", msg_id)
    else:
        print("   ❌ 未收到减法响应")
    
    print()
    
    # 测试乘法
    print("3. 测试乘法功能 (8 × 7)")
    redis_client.xadd(
        "template-service-input-stream",
        {
            "message_id": "",
            "reply_id": "",
            "service_name": "multiply",
            "callback_stream": "test-response-stream",
            "playload": json.dumps({"num1": 8, "num2": 7})
        }
    )
    
    time.sleep(0.5)
    
    messages = redis_client.xread(
        streams={"test-response-stream": "0"},
        count=1,
        block=3000
    )
    
    if messages:
        for stream, stream_messages in messages:
            for msg_id, msg_data in stream_messages:
                playload = json.loads(msg_data.get('playload', '{}'))
                print(f"   乘法结果: {json.dumps(playload, indent=4)}")
                redis_client.xdel("test-response-stream", msg_id)
    else:
        print("   ❌ 未收到乘法响应")
    
    print()
    
    # 测试除法
    print("4. 测试除法功能 (100 ÷ 20)")
    redis_client.xadd(
        "template-service-input-stream",
        {
            "message_id": "",
            "reply_id": "",
            "service_name": "divide",
            "callback_stream": "test-response-stream",
            "playload": json.dumps({"num1": 100, "num2": 20})
        }
    )
    
    time.sleep(0.5)
    
    messages = redis_client.xread(
        streams={"test-response-stream": "0"},
        count=1,
        block=3000
    )
    
    if messages:
        for stream, stream_messages in messages:
            for msg_id, msg_data in stream_messages:
                playload = json.loads(msg_data.get('playload', '{}'))
                print(f"   除法结果: {json.dumps(playload, indent=4)}")
                redis_client.xdel("test-response-stream", msg_id)
    else:
        print("   ❌ 未收到除法响应")
    
    print()
    
    # 测试echo功能
    print("5. 测试echo功能")
    redis_client.xadd(
        "template-service-input-stream",
        {
            "message_id": "",
            "reply_id": "",
            "service_name": "echo",
            "callback_stream": "test-response-stream",
            "playload": json.dumps({"message": "Hello, Go Service!"})
        }
    )
    
    time.sleep(0.5)
    
    messages = redis_client.xread(
        streams={"test-response-stream": "0"},
        count=1,
        block=3000
    )
    
    if messages:
        for stream, stream_messages in messages:
            for msg_id, msg_data in stream_messages:
                playload = json.loads(msg_data.get('playload', '{}'))
                print(f"   Echo响应: {json.dumps(playload, indent=4)}")
                redis_client.xdel("test-response-stream", msg_id)
    else:
        print("   ❌ 未收到echo响应")
    
    print()
    print("测试完成！")

if __name__ == "__main__":
    test_go_service()
