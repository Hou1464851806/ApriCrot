package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	execFile := "./electron/electron.exe"
	err := createDaemon("daemon.log", true)
	if err != nil {
		log.Println("守护进程启动失败", err)
	}
	log.Println("现在守护进程已经启动；", "我的进程号：", os.Getpid(), "；我的父进程号：", os.Getppid())
	log.Println("正在运行守护进程")
	//run("link.log", "./test.py")
	daemonize("service.log", execFile)
}

//func run(logFile string, execFile string) {
//	serviceLog, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
//	if err != nil {
//		log.Println("日志文件打开失败", err)
//	}
//	cmd := exec.Command("python", execFile)
//	cmd.Stderr = serviceLog
//	cmd.Stdout = serviceLog
//	err = cmd.Start()
//	if err != nil {
//		log.Println("处理进程启动失败")
//	}
//}

func daemonize(logFile string, execFile string) {
	serviceLog, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("日志文件打开失败", err)
	}

	for {
		cmd := exec.Command(execFile)
		cmd.Stderr = serviceLog
		cmd.Stdout = serviceLog
		err := cmd.Start()
		if err != nil {
			log.Println("服务启动失败", err)
		}
		cmd.Wait()
		log.Println("服务再次启动")
	}
}

func createDaemon(logFile string, isBackground bool) error {
	envName := "PROC_TYPE"
	envValue := "DAEMON"
	env := os.Getenv(envName)
	if env == envValue {
		return nil
	}

	cmd := &exec.Cmd{Path: os.Args[0], Args: os.Args, Env: os.Environ()}
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", envName, envValue))

	if logFile != "" {
		daemonLog, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Println("日志文件打开错误", err)
		}
		cmd.Stderr = daemonLog
		cmd.Stdout = daemonLog
	}

	err := cmd.Start()
	if err != nil {
		return err
	}
	if isBackground {
		os.Exit(0)
	}
	return nil
}
