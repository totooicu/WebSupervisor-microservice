import redis
import json
import time
import sys
import os

# 添加项目根目录到Python路径
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

# 从正确的路径导入StreamMessage
from microservice.template_service_python.models import StreamMessage

# 创建Redis客户端
redis_client = redis.Redis(host='localhost', port=6379, decode_responses=True)

def test_go_service():
    print("=== 测试Go版本服务 ===")
    
    # 创建测试消息
    msg = StreamMessage(
        service_name="add",
        callback_stream="test-response-stream",
        playload={
            "num1": 10,
            "num2": 20
        }
    )
    
    # 发送消息到Go服务
    print("发送加法请求到Go服务...")
    redis_client.xadd(
        "template-service-input-stream",
        {"message": msg.to_json()}
    )
    
    # 等待响应
    print("等待Go服务响应...")
    time.sleep(1)
    
    # 从响应流获取消息
    messages = redis_client.xread(
        streams={"test-response-stream": ">"},
        count=1,
        block=3000
    )
    
    if messages:
        for stream, stream_messages in messages:
            for msg_id, msg_data in stream_messages:
                message_json = msg_data.get('message', '{}')
                msg_dict = json.loads(message_json)
                print(f"Go服务响应: {json.dumps(msg_dict, indent=2, ensure_ascii=False)}")
    else:
        print("未收到Go服务响应")
    
    print()

def test_python_service():
    print("=== 测试Python版本服务 ===")
    
    # 创建测试消息
    msg = StreamMessage(
        service_name="add",
        callback_stream="test-response-stream",
        playload={
            "num1": 15,
            "num2": 25
        }
    )
    
    # 发送消息到Python服务
    print("发送加法请求到Python服务...")
    redis_client.xadd(
        "template-service-python-input-stream",
        {"message": msg.to_json()}
    )
    
    # 等待响应
    print("等待Python服务响应...")
    time.sleep(1)
    
    # 从响应流获取消息
    messages = redis_client.xread(
        streams={"test-response-stream": ">"},
        count=1,
        block=3000
    )
    
    if messages:
        for stream, stream_messages in messages:
            for msg_id, msg_data in stream_messages:
                message_json = msg_data.get('message', '{}')
                msg_dict = json.loads(message_json)
                print(f"Python服务响应: {json.dumps(msg_dict, indent=2, ensure_ascii=False)}")
    else:
        print("未收到Python服务响应")
    
    print()

def test_subtract_service():
    print("=== 测试减法功能 ===")
    
    # 测试Go服务减法
    msg = StreamMessage(
        service_name="subtract",
        callback_stream="test-response-stream",
        playload={
            "num1": 50,
            "num2": 20
        }
    )
    
    print("发送减法请求到Go服务...")
    redis_client.xadd(
        "template-service-input-stream",
        {"message": msg.to_json()}
    )
    
    time.sleep(1)
    
    messages = redis_client.xread(
        streams={"test-response-stream": ">"},
        count=1,
        block=3000
    )
    
    if messages:
        for stream, stream_messages in messages:
            for msg_id, msg_data in stream_messages:
                message_json = msg_data.get('message', '{}')
                msg_dict = json.loads(message_json)
                print(f"减法结果: {json.dumps(msg_dict, indent=2, ensure_ascii=False)}")
    else:
        print("未收到减法响应")
    
    print()

def test_multiply_service():
    print("=== 测试乘法功能 ===")
    
    # 测试Python服务乘法
    msg = StreamMessage(
        service_name="multiply",
        callback_stream="test-response-stream",
        playload={
            "num1": 8,
            "num2": 7
        }
    )
    
    print("发送乘法请求到Python服务...")
    redis_client.xadd(
        "template-service-python-input-stream",
        {"message": msg.to_json()}
    )
    
    time.sleep(1)
    
    messages = redis_client.xread(
        streams={"test-response-stream": ">"},
        count=1,
        block=3000
    )
    
    if messages:
        for stream, stream_messages in messages:
            for msg_id, msg_data in stream_messages:
                message_json = msg_data.get('message', '{}')
                msg_dict = json.loads(message_json)
                print(f"乘法结果: {json.dumps(msg_dict, indent=2, ensure_ascii=False)}")
    else:
        print("未收到乘法响应")
    
    print()

def test_divide_service():
    print("=== 测试除法功能 ===")
    
    # 测试Go服务除法
    msg = StreamMessage(
        service_name="divide",
        callback_stream="test-response-stream",
        playload={
            "num1": 100,
            "num2": 20
        }
    )
    
    print("发送除法请求到Go服务...")
    redis_client.xadd(
        "template-service-input-stream",
        {"message": msg.to_json()}
    )
    
    time.sleep(1)
    
    messages = redis_client.xread(
        streams={"test-response-stream": ">"},
        count=1,
        block=3000
    )
    
    if messages:
        for stream, stream_messages in messages:
            for msg_id, msg_data in stream_messages:
                message_json = msg_data.get('message', '{}')
                msg_dict = json.loads(message_json)
                print(f"除法结果: {json.dumps(msg_dict, indent=2, ensure_ascii=False)}")
    else:
        print("未收到除法响应")
    
    print()

def test_echo_service():
    print("=== 测试Echo功能 ===")
    
    # 测试Python服务echo
    msg = StreamMessage(
        service_name="echo",
        callback_stream="test-response-stream",
        playload={
            "message": "Hello from Python test script!"
        }
    )
    
    print("发送Echo请求到Python服务...")
    redis_client.xadd(
        "template-service-python-input-stream",
        {"message": msg.to_json()}
    )
    
    time.sleep(1)
    
    messages = redis_client.xread(
        streams={"test-response-stream": ">"},
        count=1,
        block=3000
    )
    
    if messages:
        for stream, stream_messages in messages:
            for msg_id, msg_data in stream_messages:
                message_json = msg_data.get('message', '{}')
                msg_dict = json.loads(message_json)
                print(f"Echo结果: {json.dumps(msg_dict, indent=2, ensure_ascii=False)}")
    else:
        print("未收到Echo响应")
    
    print()

if __name__ == "__main__":
    print("开始测试Template Service...")
    print("=" * 50)
    
    try:
        # 测试各个功能
        test_go_service()
        test_python_service()
        test_subtract_service()
        test_multiply_service()
        test_divide_service()
        test_echo_service()
        
        print("测试完成！")
        
    except Exception as e:
        print(f"测试过程中出现错误: {e}")
    finally:
        # 清理测试流
        redis_client.delete("test-response-stream")
        print("测试流已清理")