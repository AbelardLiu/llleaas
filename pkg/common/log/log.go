package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"sync"
)

var instance *logrus.Logger
var logOnce sync.Once

var system string = runtime.GOOS
var useStdout bool = false
var logLevel string = "info"
var useCaller bool = false

func UseStdOut() {
	useStdout = true
}

func UseCaller() {
	useCaller = true
}

func SetLogLevel(level string) {
	logLevel = level
}

func GetLogger() *logrus.Logger {
	logOnce.Do( func() {
		var logFile string
		logPath := generateFilePath()

		logFile = logPath + "/" + "instance.log"

		instance = logrus.New()

		switch logLevel {
		case "debug": instance.SetLevel(logrus.DebugLevel)
		case "info": instance.SetLevel(logrus.InfoLevel)
		case "warn": instance.SetLevel(logrus.WarnLevel)
		case "error": instance.SetLevel(logrus.ErrorLevel)
		case "fatal": instance.SetLevel(logrus.FatalLevel)
		case "trace": instance.SetLevel(logrus.TraceLevel)
		default:
			panic("illegal log level input")
		}

		instance.SetReportCaller(useCaller)

		customFormatter := new(logrus.TextFormatter)
		customFormatter.FullTimestamp = true
		customFormatter.TimestampFormat = "2006-01-02 15:04:05:06"
		instance.SetFormatter(customFormatter)

		file, err := os.OpenFile(logFile, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
		if err != nil {
			instance.Fatal(err)
		} else {
			if useStdout {
				instance.SetOutput(os.Stdout)
			} else {
				instance.SetOutput(file)
			}

		}
	} )

	return instance
}

func generateFilePath() string {
	var logPathDir string

	logPathDir = "./llleaas/log"

	exist := pathExist(logPathDir)
	if !exist {
		err := os.MkdirAll(logPathDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return logPathDir
}

func pathExist(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}