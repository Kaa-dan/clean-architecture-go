package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Init(level string) {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{})

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	log.SetLevel(logLevel)
}

func GetLogger() *logrus.Logger {
	if log == nil {
		Init("info")
	}
	return log
}

func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}
