import logging
from models import StreamMessage, EchoParameter, AddParameter, SubtractParameter, MultiplyParameter, DivideParameter


class Handlers:
    def __init__(self, service):
        self.service = service
        self.logger = logging.getLogger(__name__)
    
    def handle_echo(self, msg: StreamMessage):
        """处理echo服务"""
        if self.service.debug:
            self.logger.info(f"Handling echo message: {msg.message_id}")
        
        try:
            # 解析参数
            params = self.service.parse_parameters(msg, EchoParameter)
            
            # 处理业务逻辑（echo返回相同的消息）
            self.service.send_success_response(msg, {
                "echo": params.message
            })
            
        except Exception as e:
            self.logger.error(f"Error handling echo: {e}")
            self.service.send_error_response(msg, f"Invalid parameters: {str(e)}")
    
    def handle_add(self, msg: StreamMessage):
        """处理加法服务"""
        if self.service.debug:
            self.logger.info(f"Handling add message: {msg.message_id}")
        
        try:
            # 解析参数
            params = self.service.parse_parameters(msg, AddParameter)
            
            # 处理业务逻辑（加法运算）
            result = params.num1 + params.num2
            self.service.send_success_response(msg, {
                "result": result,
                "num1": params.num1,
                "num2": params.num2
            })
            
        except Exception as e:
            self.logger.error(f"Error handling add: {e}")
            self.service.send_error_response(msg, f"Invalid parameters: {str(e)}")
    
    def handle_subtract(self, msg: StreamMessage):
        """处理减法服务"""
        if self.service.debug:
            self.logger.info(f"Handling subtract message: {msg.message_id}")
        
        try:
            # 解析参数
            params = self.service.parse_parameters(msg, SubtractParameter)
            
            # 处理业务逻辑（减法运算）
            result = params.num1 - params.num2
            self.service.send_success_response(msg, {
                "result": result,
                "num1": params.num1,
                "num2": params.num2
            })
            
        except Exception as e:
            self.logger.error(f"Error handling subtract: {e}")
            self.service.send_error_response(msg, f"Invalid parameters: {str(e)}")
    
    def handle_multiply(self, msg: StreamMessage):
        """处理乘法服务"""
        if self.service.debug:
            self.logger.info(f"Handling multiply message: {msg.message_id}")
        
        try:
            # 解析参数
            params = self.service.parse_parameters(msg, MultiplyParameter)
            
            # 处理业务逻辑（乘法运算）
            result = params.num1 * params.num2
            self.service.send_success_response(msg, {
                "result": result,
                "num1": params.num1,
                "num2": params.num2
            })
            
        except Exception as e:
            self.logger.error(f"Error handling multiply: {e}")
            self.service.send_error_response(msg, f"Invalid parameters: {str(e)}")
    
    def handle_divide(self, msg: StreamMessage):
        """处理除法服务"""
        if self.service.debug:
            self.logger.info(f"Handling divide message: {msg.message_id}")
        
        try:
            # 解析参数
            params = self.service.parse_parameters(msg, DivideParameter)
            
            # 检查除数是否为零
            if params.num2 == 0:
                self.service.send_error_response(msg, "Division by zero is not allowed")
                return
            
            # 处理业务逻辑（除法运算）
            result = params.num1 / params.num2
            self.service.send_success_response(msg, {
                "result": result,
                "num1": params.num1,
                "num2": params.num2
            })
            
        except Exception as e:
            self.logger.error(f"Error handling divide: {e}")
            self.service.send_error_response(msg, f"Invalid parameters: {str(e)}")