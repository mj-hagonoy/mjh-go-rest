package logger

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
)

var (
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
)

func InitLoggers() {
	logdir := config.GetConfig().Log.LogDir
	if err := os.MkdirAll(logdir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	initInfoLogger()
	initErrorLogger()
}

func initErrorLogger() {
	logdir := config.GetConfig().Log.LogDir
	logpath := path.Clean(fmt.Sprintf("%s/error.log", logdir))
	if errorFile, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		log.Fatal(err)
	} else {
		ErrorLogger = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

func initInfoLogger() {
	logdir := config.GetConfig().Log.LogDir
	logpath := path.Clean(fmt.Sprintf("%s/info.log", logdir))
	if errorFile, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		log.Fatal(err)
	} else {
		InfoLogger = log.New(errorFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
}
