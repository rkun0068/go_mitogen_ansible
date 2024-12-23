package output

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Creates log files and retunrns the log recorder
func CreateLogFile(ip string) (*os.File, *log.Logger, error) {
	timestamp := time.Now().Format("20020202_222222")

	logDir := fmt.Sprintf("/tmp/fast-ansible-log/%s", timestamp)

	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		return nil, nil, fmt.Errorf("func CreateLogFile error: %v", err)
	}
	ip = strings.ReplaceAll(ip, ".", "-")
	logFilePath := fmt.Sprintf("%s/%s.log", logDir, ip)
	// open log file
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("func CreateLogFile error: %v", err)
	}
	// create a log recorder
	logger := log.New(logFile, "", log.LstdFlags)

	return logFile, logger, nil

}

func CloseLogFile(logFile *os.File) error {
	if err := logFile.Close(); err != nil {
		return fmt.Errorf("func CloseLogFile error: %v", err)
	}
	return nil
}
