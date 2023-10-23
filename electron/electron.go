package main

import (
	"Daemon-test/electron/everything"
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

var numOfEXE int

func main() {
	linkFile := "./ApriCrot/__main__.py"
	app, err := os.OpenFile("app.txt", os.O_CREATE, 0666)
	if err != nil {
		log.Println("文件打开失败")
	}
	defer app.Close()
	//appFile := "app.txt"
	writer := bufio.NewWriter(app)
	reader := bufio.NewReader(app)

	firstLine, err := reader.ReadString('\n')
	if err == nil {
		firstLine = firstLine[0 : len(firstLine)-1]
		numOfEXE, _ = strconv.Atoi(firstLine)
	} else {
	}

	for {
		everything.SetSearch("(atom|CefDetectorX).exe")
		everything.SetRegex(true)
		everything.SetRequestFlags(everything.EVERYTHING_REQUEST_FILE_NAME | everything.EVERYTHING_REQUEST_PATH | everything.EVERYTHING_REQUEST_SIZE)
		everything.SetSort(everything.EVERYTHING_SORT_NAME_ASCENDING)
		everything.Query(true)
		log.Println("everything查询结果数量", everything.GetNumResults())
		if numOfEXE != everything.GetNumResults() {
			numOfEXE = everything.GetNumResults()
			writer.WriteString(strconv.Itoa(numOfEXE) + "\n")
			for i := 0; i < everything.GetNumResults(); i++ {
				fileName := everything.GetResultFileName(i)
				filePath := everything.GetResultPathW(i)
				fullFile := fmt.Sprintf("%s\\%s", filePath, fileName)
				if isElectron, version := filterElectron(fullFile); isElectron {
					link("link.log", linkFile, fullFile)
					fullFile = fullFile + " " + version + "\n"
					writer.WriteString(fullFile)
					log.Printf(fullFile)
				}
			}
		}
		writer.Flush()
		log.Println("ElectronApp全部找出")
		//link("link.log", linkFile)
		time.Sleep(time.Hour * 1)
	}
}

func filterElectron(filePath string) (bool, string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		//log.Println("读取文件失败", err)
		return false, ""
	}
	content := string(data)
	re := regexp.MustCompile("Electron/([0-9]+).([0-9]+).([0-9]+)")
	res := re.FindStringIndex(content)
	if res != nil {
		electronVersion := content[res[0]:res[1]]
		return true, electronVersion
	}
	return false, ""
}

func link(logFile string, execFile string, appFile string) {
	serviceLog, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("日志文件打开失败", err)
	}
	cmd := exec.Command("python", execFile, appFile)
	cmd.Stderr = serviceLog
	cmd.Stdout = serviceLog
	err = cmd.Start()
	if err != nil {
		log.Println("处理进程启动失败")
	}
}
