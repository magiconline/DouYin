package logger

import (
	"log"
	"os"
)

var Logger log.Logger

func Init(logFilePath string) error {
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	Logger = log.Logger{}
	Logger.SetOutput(logFile)
	return nil
}
