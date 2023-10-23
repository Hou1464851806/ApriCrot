import sys

from loguru import logger
from link.link import link

def main():
    logger.remove()
    log_format = "<level>{level: ^8}</level> | <level>{message}</level>"

    logger.add(sys.stdout, format=log_format, filter=lambda log_instance: log_instance['level'].name == "DEBUG")
    logger.add(sys.stderr, format=log_format, filter=lambda log_instance: log_instance['level'].name == "WARNING")
    logger.add(sys.stdout, format=log_format, filter=lambda log_instance: log_instance['level'].name == "INFO")
    logger.add(sys.stderr, format=log_format, filter=lambda log_instance: log_instance['level'].name == "ERROR")

    if sys.platform != 'win32':
        logger.error('当前仅支持 Windows 平台，其他平台敬请期待')
        exit(0)
        
    
    if len(sys.argv) != 2:
        print("Usage: my_python_script.py <input_string>")
        sys.exit(1)

    input_string = sys.argv[1]
    result = f"Python received: {input_string}"
    print(result)

    logger.debug('加载成功')

    link(str(sys.argv[1]))

if __name__ == '__main__':
    main()