import json
from typing import Dict, Any, Optional


class RedisConfig:
    def __init__(self, host: str, port: int, password: str, db: int):
        self.host = host
        self.port = port
        self.password = password
        self.db = db
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'RedisConfig':
        return cls(
            host=data.get('host', 'localhost'),
            port=data.get('port', 6379),
            password=data.get('password', ''),
            db=data.get('db', 0)
        )


class BaseConfig:
    def __init__(self, redis: RedisConfig, consumer_group: str, input_stream: str):
        self.redis = redis
        self.consumer_group = consumer_group
        self.input_stream = input_stream
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'BaseConfig':
        redis_config = RedisConfig.from_dict(data.get('redis', {}))
        return cls(
            redis=redis_config,
            consumer_group=data.get('consumer_group', 'template-service-python-group'),
            input_stream=data.get('input_stream', 'template-service-python-input-stream')
        )


class StreamMessage:
    def __init__(self, message_id: str = "", reply_id: str = "", service_name: str = "", 
                 callback_stream: str = "", playload: Dict[str, Any] = None):
        self.message_id = message_id
        self.reply_id = reply_id
        self.service_name = service_name
        self.callback_stream = callback_stream
        self.playload = playload or {}
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "message_id": self.message_id,
            "reply_id": self.reply_id,
            "service_name": self.service_name,
            "callback_stream": self.callback_stream,
            "playload": self.playload
        }
    
    def to_json(self) -> str:
        return json.dumps(self.to_dict(), ensure_ascii=False, indent=2)
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'StreamMessage':
        return cls(
            message_id=data.get('message_id', ''),
            reply_id=data.get('reply_id', ''),
            service_name=data.get('service_name', ''),
            callback_stream=data.get('callback_stream', ''),
            playload=data.get('playload', {})
        )


class EchoParameter:
    def __init__(self, message: str):
        self.message = message
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'EchoParameter':
        return cls(message=data.get('message', ''))


class AddParameter:
    def __init__(self, num1: float, num2: float):
        self.num1 = num1
        self.num2 = num2
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'AddParameter':
        return cls(
            num1=data.get('num1', 0.0),
            num2=data.get('num2', 0.0)
        )


class SubtractParameter:
    def __init__(self, num1: float, num2: float):
        self.num1 = num1
        self.num2 = num2
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'SubtractParameter':
        return cls(
            num1=data.get('num1', 0.0),
            num2=data.get('num2', 0.0)
        )


class MultiplyParameter:
    def __init__(self, num1: float, num2: float):
        self.num1 = num1
        self.num2 = num2
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'MultiplyParameter':
        return cls(
            num1=data.get('num1', 0.0),
            num2=data.get('num2', 0.0)
        )


class DivideParameter:
    def __init__(self, num1: float, num2: float):
        self.num1 = num1
        self.num2 = num2
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'DivideParameter':
        return cls(
            num1=data.get('num1', 0.0),
            num2=data.get('num2', 0.0)
        )