package logger

import (
	"log"
	"os"
)

var Logger log.Logger

var Println = Logger.Println
var Fatalln = Logger.Fatalln

func Init(logFilePath string) error {
	// 0644: 用户具有读写权限，组用户和其它用户具有只读权限；
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	Logger = *log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}
