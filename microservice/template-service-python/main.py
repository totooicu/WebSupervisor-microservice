import argparse
import json
import logging
from models import BaseConfig
from service import TemplateService
from handlers import Handlers


def load_config(config_path: str) -> BaseConfig:
    """加载配置文件"""
    try:
        with open(config_path, 'r', encoding='utf-8') as f:
            config_data = json.load(f)
        return BaseConfig.from_dict(config_data)
    except Exception as e:
        logging.error(f"Failed to load config: {e}")
        raise


def main():
    """主函数"""
    # 解析命令行参数
    parser = argparse.ArgumentParser(description='Template Service Python Version')
    parser.add_argument('--config', type=str, default='./config.json', help='Path to config file')
    parser.add_argument('--debug', action='store_true', help='Enable debug logging')
    args = parser.parse_args()
    
    # 设置日志
    logging.basicConfig(
        level=logging.DEBUG if args.debug else logging.INFO,
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )
    logger = logging.getLogger(__name__)
    
    if args.debug:
        logger.info("Debug mode enabled")
    
    # 加载配置
    try:
        config = load_config(args.config)
        logger.info(f"Config loaded from: {args.config}")
    except Exception as e:
        logger.error(f"Failed to load config: {e}")
        return
    
    # 创建服务
    service = TemplateService(config, args.debug)
    
    # 创建处理器
    handlers = Handlers(service)
    
    # 注册服务处理器
    service.register_handler('echo', handlers.handle_echo)
    service.register_handler('add', handlers.handle_add)
    service.register_handler('subtract', handlers.handle_subtract)
    service.register_handler('multiply', handlers.handle_multiply)
    service.register_handler('divide', handlers.handle_divide)
    
    logger.info("All handlers registered successfully")
    
    # 启动服务
    try:
        service.start()
    except KeyboardInterrupt:
        logger.info("Service stopped by user")
    except Exception as e:
        logger.error(f"Service error: {e}")


if __name__ == "__main__":
    main()