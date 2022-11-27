package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func New(fileName string, level string) *logrus.Logger {
	logger := logrus.New()
	var mylevel logrus.Level
	switch strings.ToUpper(level) {
	case "FATAL":
		mylevel = logrus.FatalLevel
	case "ERROR":
		mylevel = logrus.ErrorLevel
	case "WARNING":
		mylevel = logrus.WarnLevel
	case "INFO":
		mylevel = logrus.InfoLevel
	case "DEBUG":
		mylevel = logrus.DebugLevel
	default:
		mylevel = logrus.TraceLevel
	}
	logger.Level = mylevel
	if fileName != "" {
		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
		if err == nil {
			logger.Out = file
		} else {
			logger.Out = os.Stdout
		}
	} else {
		logger.Out = os.Stdout
	}
	return logger
}
