import pefile
import json
import io
import os
import zipfile
import requests
import re
import shutil
import subprocess
import sys
sys.path.append(r"../")

from loguru import logger
from pathlib import Path
from typing import Iterable

from information.information import __fullname__, __electron_repo_root__, __electron_source__, __electron_repo__, __platform__, PEError, DownloadError
from database.database import updateDB, checkDB

http_proxy = os.environ.get("HTTP_PROXY")
https_proxy = os.environ.get("HTTPS_PROXY")

def pathFormatChange(drive: str) -> Path:
    return Path(drive).joinpath('/')

def backUpfile(target:Path, backUpFilePath:Path, sourceFilePath:Path, backup_history:list):
    try:
        backUpFilePath=backUpFilePath.joinpath(Path(sourceFilePath.relative_to(target)))
        if not backUpFilePath.parent.exists():
            os.makedirs(backUpFilePath.parent)
        shutil.copy(sourceFilePath, backUpFilePath)
        backup_info = (str(sourceFilePath), str(backUpFilePath))
        backup_history.append(backup_info)
        logger.debug(f"将 {sourceFilePath} 备份到 {backUpFilePath}")
    except Exception as e:
        logger.error(f"备份 {sourceFilePath} 到 {backUpFilePath} 时发生错误: {str(e)}")
    return backup_history

def restore_backup(backup_history:list):
    for source, backup in backup_history:
        sourcePath=Path(source)
        backupPath=Path(backup)
        try:
            # 还原文件
            shutil.copy(backupPath, sourcePath)
            logger.debug(f"还原 {backup} 到 {source}")
        except Exception as e:
            logger.error(f"还原 {backup} 到 {source} 时发生错误: {str(e)}")

def create_link(repo: Path, repo_name: Path | str, target: Path, target_name: Path):
    logger.debug(f"将 {repo_name} 链接到 {target_name}")
    repo_file = repo.joinpath(repo_name)
    target_file = target.joinpath(target_name)
    target_file.unlink()
    target_file.hardlink_to(repo_file)

def getAllFiles(root: str):
    files_and_directories = list(map(lambda filename: root + '\\' + filename, os.listdir(root)))
    files = []

    for path in files_and_directories:
        if (os.path.isfile(path)):
            files.append(path)
        elif (os.path.isdir(path)):
            files_and_directories.extend(list(map(lambda filename: path + '\\' + filename, os.listdir(path))))

    return files

def link(app_entry_str:str):
    app_entry=Path(app_entry_str)
    logger.debug(app_entry)

    with pefile.PE(app_entry, fast_load=True) as pe:
        match pefile.MACHINE_TYPE[pe.FILE_HEADER.Machine]:
            case 'IMAGE_FILE_MACHINE_I386':
                electron_arch = 'ia32'
            case 'IMAGE_FILE_MACHINE_AMD64':
                electron_arch = 'x64'
            case 'IMAGE_FILE_MACHINE_ARM64':
                electron_arch = 'arm64'
            case _:
                msg = f'{app_entry.name} 的 CPU 架构未知:{pefile.MACHINE_TYPE[pe.FILE_HEADER.Machine]}'
                logger.warning(msg)
                raise PEError(msg)
        logger.debug(electron_arch)

        section_rdata = list(filter(lambda section: section.Name.strip(b'\x00') == b'.rdata', pe.sections))
        if len(section_rdata) == 0:
            msg = f'{app_entry.name} 无 .rdata 段'
            logger.warning(msg)
            raise PEError(msg)
        elif len(section_rdata) > 1:
            logger.warning(f'{app_entry.name} 的 .rdata 段不唯一，默认使用第一个')
        section_rdata = section_rdata[0]
        rdata = pe.get_data(section_rdata.VirtualAddress, section_rdata.SizeOfRawData)

    versions = [i.decode() for i in (set(re.findall(rb'Chrome/(?:[0-9.]+?|%s) Electron/(\S+?)\x00', rdata)))]
    logger.debug(versions)
    if len(versions) == 0:
        msg = f'{app_entry.name} 的 .rdata 段中找不到版本信息，如确认为应用入口'
        logger.warning(msg)
        raise PEError(msg)
    elif len(versions) > 1:
        logger.warning(f'{app_entry.name} 的 .rdata 段中版本信息不唯一:{versions}，默认使用第一个')
    electron_version = versions[0]

    drive = pathFormatChange(app_entry.drive)

    repo = drive.joinpath(__electron_repo_root__).joinpath(__electron_repo__.format(version=electron_version, arch=electron_arch))

    if repo.exists() and not repo.is_dir():
        logger.warning(f"{repo} 存在但不是文件夹, 自动删除")
        repo.unlink()
    if not repo.exists():
        logger.warning(f"{repo} 不存在, 自动创建")
        try:
            repo.mkdir(parents=True)
        except FileExistsError as e:
            msg = "无法创建链接仓库，可能是链接仓库根目录不是目录，为防止误删，请自行检查并删除"
            logger.warning(msg)
            raise DownloadError(msg)

    logger.debug("正在下载 Electron 预编译程序")
    logger.debug(f"目标链接仓库: {repo}")

    source = __electron_source__[1]
    electron_url = source.format(version=electron_version, platform=__platform__, arch=electron_arch)
    logger.debug(f"下载地址: {electron_url}")

    signed_file = repo.joinpath('apricrot.signed.json')
    if signed_file.exists():
        logger.warning(f"{signed_file} 存在, 自动删除")
        signed_file.unlink()  # 删除校验文件，等价于将此目录标记为未下载完成/下载失败

    try:
        electron_resp = requests.get(electron_url, proxies={
            http_proxy,
            https_proxy
        } if (http_proxy or https_proxy) else None)
    except requests.exceptions.InvalidProxyURL as e:
        msg = "代理格式存在问题，请检查工具中或系统中的代理设置"
        logger.warning(msg)
        raise DownloadError(msg)
    except requests.exceptions.ProxyError as e:
        msg = "无法连接到代理，请检查工具中或系统中的代理设置"
        logger.warning(msg)
        raise DownloadError(msg)
    except requests.exceptions.ConnectionError as e:
        msg = "网络存在问题"
        logger.warning(msg)
        raise DownloadError(msg)

    if electron_resp.status_code != 200:
        msg = f"下载失败，HTTP {electron_resp.status_code}"
        logger.warning(msg)
        raise DownloadError(msg)

    signed = {}
    with io.BytesIO(electron_resp.content) as f:
        with zipfile.ZipFile(f, 'r') as zipf:
            for file in zipf.filelist:
                filename = repo.joinpath(file.filename)
                if filename.parent != repo:
                    if not filename.parent.exists():
                        filename.parent.mkdir(parents=True)
                with open(filename, 'wb') as rf:
                    # 通过覆写文件的方式来更新已被硬链接的应用，如果直接替换文件会导致已有的硬链接失效
                    rf.write(zipf.read(file.filename))
                signed[file.filename] = hex(file.CRC)
    with open(signed_file, 'w') as f:  # 写入校验文件，等价于将此目录标记为下载完成
        f.write(json.dumps(signed))
    logger.debug("下载 Electron 预编译程序完成")

    target = app_entry.parent
    files = getAllFiles(str(target))
    Ofiles = getAllFiles(str(repo))
    Jfiles = []

    for absPath in files:
        relaPath = Path(absPath).relative_to(target)
        if Ofiles.count(str(repo.joinpath(relaPath))) != 0:
            Jfiles.append(absPath)

    backUpfiles=repo.joinpath(app_entry.relative_to(target))
    os.makedirs(backUpfiles)

    backup_history = []

    for file in Jfiles:
        fileP = Path(file)
        relative_name = fileP.relative_to(target)
        backup_history = backUpfile(target, backUpfiles, fileP, backup_history)
        create_link(repo,  relative_name, target, relative_name)

    try:
        test=subprocess.Popen(app_entry, shell=False)
        test.terminate()
        shutil.rmtree(backUpfiles)
        logger.debug("删除备份")
        updateDB(app_entry_str, electron_version, electron_arch)
        logger.debug("成功")
    except Exception as e:
        logger.error(f"文件无法正常打开，已还原文件: {str(e)}")
        restore_backup(backup_history)
        shutil.rmtree(backUpfiles)
        logger.debug("删除备份")
        #这里现在会让electron留下，将来再考虑怎么处理

