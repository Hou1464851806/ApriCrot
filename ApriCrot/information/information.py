__version__ = "0.1.0"
#__description__ = "一个通过将磁盘上所有 Electron 应用中相同文件硬链接到统一位置来减少磁盘占用的解决方案，就像 pnpm 一样。"
__fullname__ = f'ApriCrot v{__version__}'
__electron_source__ = [
    'https://github.com/electron/electron/releases/download/v{version}/electron-v{version}-{platform}-{arch}.zip',
    'https://repo.huaweicloud.com/electron/{version}/electron-v{version}-{platform}-{arch}.zip',
    'https://registry.npmmirror.com/-/binary/electron/{version}/electron-v{version}-{platform}-{arch}.zip',
]
__electron_repo_root__ = ".apricrot"
__electron_repo__ = "{version}-{arch}"
__electron_repo_re__ = r"^(\S+)-(\S+?)$"
__platform__ = "win32"

class RepoError(Exception):
    """
    链接仓库错误
    """
    pass


class ScanError(Exception):
    """
    扫描时错误
    """
    pass


class TargetError(Exception):
    """
    目标 Electron 应用错误
    """
    pass


class DownloadError(Exception):
    """
    下载时错误
    """
    pass


class PEError(Exception):
    """
    解析 PE 文件时错误
    """
    pass
