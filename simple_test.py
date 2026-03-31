import redis
import json
import time

# 创建Redis客户端
redis_client = redis.Redis(host='localhost', port=6379, decode_responses=True)

def test_add_go():
    print("=== 测试Go版本加法 ===")
    
    print("发送加法请求到Go服务...")
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
    
    time.sleep(1)
    
    messages = redis_client.xread(
        streams={"test-response-stream": "0"},
        count=1,
        block=3000
    )
    
    if messages:
        for stream, stream_messages in messages:
            for msg_id, msg_data in stream_messages:
                playload = json.loads(msg_data.get('playload', '{}'))
                print(f"Go服务响应: {json.dumps(playload, indent=2, ensure_ascii=False)}")
                redis_client.xdel("test-response-stream", msg_id)
    else:
        print("未收到Go服务响应")
    
    print()

def test_add_python():
    print("=== 测试Python版本加法 ===")
    
    print("发送加法请求到Python服务...")
    msg_dict = {
        "message_id": "",
        "reply_id": "",
        "service_name": "add",
        "callback_stream": "test-response-stream",
        "playload": {"num1": 15, "num2": 25}
    }
    redis_client.xadd(
        "template-service-python-input-stream",
        {"message": json.dumps(msg_dict)}
    )
    
    time.sleep(1)
    
    messages = redis_client.xread(
        streams={"test-response-stream": "0"},
        count=1,
        block=3000
    )
    
    if messages:
        for stream, stream_messages in messages:
            for msg_id, msg_data in stream_messages:
                message_json = msg_data.get('message', '{}')
                msg_dict = json.loads(message_json)
                print(f"Python服务响应: {json.dumps(msg_dict['playload'], indent=2, ensure_ascii=False)}")
                redis_client.xdel("test-response-stream", msg_id)
    else:
        print("未收到Python服务响应")
    
    print()

def test_subtract():
    print("=== 测试减法功能 ===")
    
    print("发送减法请求...")
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
    
    time.sleep(1)
    
    messages = redis_client.xread(
        streams={"test-response-stream": "0"},
        count=1,
        block=3000
    )
    
    if messages:
        for stream, stream_messages in messages:
            for msg_id, msg_data in stream_messages:
                playload = json.loads(msg_data.get('playload', '{}'))
                print(f"减法结果: {json.dumps(playload, indent=2, ensure_ascii=False)}")
                redis_client.xdel("test-response-stream", msg_id)
    else:
        print("未收到减法响应")
    
    print()

def test_multiply():
    print("=== 测试乘法功能 ===")
    
    print("发送乘法请求...")
    msg_dict = {
        "message_id": "",
        "reply_id": "",
        "service_name": "multiply",
        "callback_stream": "test-response-stream",
        "playload": {"num1": 8, "num2": 7}
    }
    redis_client.xadd(
        "template-service-python-input-stream",
        {"message": json.dumps(msg_dict)}
    )
    
    time.sleep(1)
    
    messages = redis_client.xread(
        streams={"test-response-stream": "0"},
        count=1,
        block=3000
    )
    
    if messages:
        for stream, stream_messages in messages:
            for msg_id, msg_data in stream_messages:
                message_json = msg_data.get('message', '{}')
                msg_dict = json.loads(message_json)
                print(f"乘法结果: {json.dumps(msg_dict['playload'], indent=2, ensure_ascii=False)}")
                redis_client.xdel("test-response-stream", msg_id)
    else:
        print("未收到乘法响应")
    
    print()

def test_divide():
    print("=== 测试除法功能 ===")
    
    print("发送除法请求...")
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
    
    time.sleep(1)
    
    messages = redis_client.xread(
        streams={"test-response-stream": "0"},
        count=1,
        block=3000
    )
    
    if messages:
        for stream, stream_messages in messages:
            for msg_id, msg_data in stream_messages:
                playload = json.loads(msg_data.get('playload', '{}'))
                print(f"除法结果: {json.dumps(playload, indent=2, ensure_ascii=False)}")
                redis_client.xdel("test-response-stream", msg_id)
    else:
        print("未收到除法响应")
    
    print()

if __name__ == "__main__":
    print("开始测试Template Service...")
    print("=" * 50)
    
    try:
        test_add_go()
        test_add_python()
        test_subtract()
        test_multiply()
        test_divide()
        
        print("测试完成！")
        
    except Exception as e:
        print(f"测试过程中出现错误: {e}")
    finally:
        redis_client.delete("test-response-stream")
        print("测试流已清理")