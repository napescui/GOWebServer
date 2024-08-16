package logger

import (
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.Formatter = &logrus.TextFormatter{
		TimestampFormat: "01-02 15:04:05",
		FullTimestamp:   true,
	}
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}
