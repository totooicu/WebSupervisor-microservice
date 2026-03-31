import redis
import json
import logging
import threading
import time
from typing import Dict, Any, Callable
from models import BaseConfig, StreamMessage


class TemplateService:
    def __init__(self, config: BaseConfig, debug: bool = False):
        self.config = config
        self.debug = debug
        self.redis_client = None
        self.handlers = {}
        self.message_id_counter = 0
        self.message_id_lock = threading.Lock()
        
        # 设置日志
        logging.basicConfig(
            level=logging.DEBUG if debug else logging.INFO,
            format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
        )
        self.logger = logging.getLogger(__name__)
    
    def connect_redis(self):
        """连接Redis"""
        try:
            self.redis_client = redis.Redis(
                host=self.config.redis.host,
                port=self.config.redis.port,
                password=self.config.redis.password,
                db=self.config.redis.db,
                decode_responses=True
            )
            self.redis_client.ping()
            self.logger.info(f"Connected to Redis: {self.config.redis.host}:{self.config.redis.port}")
            return True
        except Exception as e:
            self.logger.error(f"Failed to connect to Redis: {e}")
            return False
    
    def create_consumer_group(self):
        """创建消费者组"""
        try:
            self.redis_client.xgroup_create(
                name=self.config.input_stream,
                groupname=self.config.consumer_group,
                id="0",
                mkstream=True
            )
            self.logger.info(f"Created consumer group: {self.config.consumer_group}")
        except redis.exceptions.ResponseError as e:
            if "BUSYGROUP" in str(e):
                self.logger.info(f"Consumer group already exists: {self.config.consumer_group}")
            else:
                self.logger.error(f"Failed to create consumer group: {e}")
                raise
    
    def get_message_id(self) -> str:
        """获取唯一消息ID"""
        with self.message_id_lock:
            self.message_id_counter += 1
            return f"{int(time.time() * 1000)}-{self.message_id_counter}"
    
    def register_handler(self, service_name: str, handler: Callable[[StreamMessage], None]):
        """注册服务处理器"""
        self.handlers[service_name] = handler
        self.logger.info(f"Registered handler for service: {service_name}")
    
    def parse_parameters(self, msg: StreamMessage, param_class) -> Any:
        """解析消息参数"""
        try:
            return param_class.from_dict(msg.playload)
        except Exception as e:
            self.logger.error(f"Failed to parse parameters: {e}")
            raise
    
    def send_response(self, msg: StreamMessage, success: bool, data: Dict[str, Any]):
        """发送响应消息"""
        if not msg.callback_stream:
            self.logger.warning("No callback stream specified, skipping response")
            return
        
        response_msg = StreamMessage(
            message_id=self.get_message_id(),
            reply_id=msg.message_id,
            service_name="response",
            callback_stream=msg.callback_stream,
            playload={
                "success": success,
                "data": data
            }
        )
        
        try:
            self.redis_client.xadd(
                msg.callback_stream,
                {"message": response_msg.to_json()}
            )
            self.logger.debug(f"Response sent to stream: {msg.callback_stream}")
        except Exception as e:
            self.logger.error(f"Failed to send response: {e}")
    
    def send_success_response(self, msg: StreamMessage, data: Dict[str, Any]):
        """发送成功响应"""
        self.send_response(msg, True, data)
    
    def send_error_response(self, msg: StreamMessage, error_msg: str):
        """发送错误响应"""
        self.send_response(msg, False, {"error": error_msg})
    
    def process_message(self, msg_id: str, msg_data: Dict[str, Any]):
        """处理消息"""
        try:
            # 解析消息
            message_json = msg_data.get('message', '{}')
            msg_dict = json.loads(message_json)
            msg = StreamMessage.from_dict(msg_dict)
            msg.message_id = msg_id
            
            self.logger.debug(f"Received message: {msg.service_name} (ID: {msg_id})")
            
            # 查找处理器
            if msg.service_name in self.handlers:
                handler = self.handlers[msg.service_name]
                handler(msg)
            else:
                self.logger.warning(f"No handler found for service: {msg.service_name}")
                self.send_error_response(msg, f"Service not found: {msg.service_name}")
            
            # 确认消息处理完成
            self.redis_client.xack(
                self.config.input_stream,
                self.config.consumer_group,
                msg_id
            )
            
        except Exception as e:
            self.logger.error(f"Error processing message {msg_id}: {e}")
    
    def start_consuming(self):
        """开始消费消息"""
        self.logger.info(f"Starting consumer for stream: {self.config.input_stream}")
        
        while True:
            try:
                # 从Stream获取消息
                messages = self.redis_client.xreadgroup(
                    groupname=self.config.consumer_group,
                    consumername="consumer-1",
                    streams={self.config.input_stream: ">"},
                    count=1,
                    block=1000  # 1秒超时
                )
                
                if messages:
                    for stream, stream_messages in messages:
                        for msg_id, msg_data in stream_messages:
                            self.process_message(msg_id, msg_data)
                
            except Exception as e:
                self.logger.error(f"Error consuming messages: {e}")
                time.sleep(1)  # 出错后等待1秒再重试
    
    def start(self):
        """启动服务"""
        # 连接Redis
        if not self.connect_redis():
            self.logger.error("Failed to connect to Redis, exiting")
            return
        
        # 创建消费者组
        self.create_consumer_group()
        
        # 开始消费消息
        self.start_consuming()