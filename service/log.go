package service

import (
	"bufio"
	"fmt"
	log "github.com/wonderivan/logger"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"ticket-service/api/apimodel"
	"ticket-service/pkg/utils"
)

const LogPath = "./logs"

func (operator *ResourceOperator) QueryLogList() (*apimodel.LogListResponse, error) {
	var err error
	var files []string
	var resp apimodel.LogListResponse
	err = utils.GetFiles(LogPath, true, &files)
	if err != nil {
		log.Error("读取路径失败. pathL:[%v],err:[%v]", LogPath, err)
		return nil, err
	}
	if files == nil {
		log.Error("目录为空 pathL:[%v],err:[%v]", LogPath, err)
		return nil, err
	}

	var fileNames []string
	for _, v := range files {
		fileNames = append(fileNames, filepath.Base(v))
	}

	resp.Load(len(fileNames), fileNames)
	return &resp, nil
}

func (operator *ResourceOperator) QueryLogData(req *apimodel.LogInfoRequest) (*apimodel.LogInfoResponse, error) {
	var logName string

	if req.LogName == "" {
		logName = "ticket-service.log"
	} else {
		logName = req.LogName
	}

	path := LogPath + "/" + logName
	fmt.Println(path)
	var resp apimodel.LogInfoResponse
	if !utils.Exists(path) {
		log.Error("路径下文件不存在. pathL:[%v],err:[%v]", path, nil)
		return nil, nil
	}
	file, err := os.Open(path)
	if err != nil {
		log.Error("打开文件失败: [%v]", err)
		return nil, err
	}
	defer file.Close()

	var logInfos []apimodel.LogInfo
	// 按行读取文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 输出每行内容
		logMsg := parseLog(scanner.Text(), logName)
		if logMsg.Request != "" {
			//剔除内部操作日志
			logInfos = append(logInfos, logMsg)
		}
	}

	resp.Load(len(logInfos), logInfos)
	return &resp, nil
}

func parseLog(log string, logName string) apimodel.LogInfo {
	logInfo := apimodel.LogInfo{}
	lines := strings.Split(log, "\n")
	for _, line := range lines {
		if strings.Contains(line, "[EROR]") {
			logInfo.ErrorType = "EROR"
		} else if strings.Contains(line, "[WARN]") {
			logInfo.ErrorType = "WARN"
		} else if strings.Contains(line, "[DEBUG]") {
			logInfo.ErrorType = "DEBUG"
		}

		// 提取时间戳
		reTime := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})`)
		times := reTime.FindStringSubmatch(line)
		if len(times) > 1 {
			logInfo.CreatedAt = times[1]
		}

		// 提取文件路径
		reFile := regexp.MustCompile(`(.*?\.go:\d+)`)
		files := reFile.FindStringSubmatch(line)
		if len(files) > 1 {
			logInfo.OriginDetail = files[1]
		}

		// 提取请求信息
		reRequest := regexp.MustCompile(`Request \[(.*?)] \[(.*?)]`)
		requests := reRequest.FindStringSubmatch(line)
		if len(requests) > 2 {
			logInfo.Request = requests[1]
			logInfo.Url = requests[2]
		}

		// 提取错误码和消息
		reError := regexp.MustCompile(`Error\.HttpCode\[(\d+)] BusinessCode\[(\d+)] Message\[([^]]+)] ErrDetail\[([^]]+)]`)
		errors := reError.FindStringSubmatch(line)
		if len(errors) > 4 {
			logInfo.ErrorCode = errors[2]
			logInfo.Msg = errors[3]
			logInfo.ErrorDetail = errors[4]
		}
	}
	logInfo.LogName = logName
	return logInfo
}
