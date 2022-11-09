//go:build !tinygo
// +build !tinygo

package log

import (
	"github.com/sirupsen/logrus"
)

// var log *logrus.Logger

func SetLogger(logLevelStr string) error {
	// log = &logrus.Logger{
	// 	Formatter: &easy.Formatter{
	// 		TimestampFormat: "15:04",
	// 		LogFormat:       "[%lvl%]: %time% - %msg%\n",
	// 	},
	// }
	logLevel, err := logrus.ParseLevel(logLevelStr)
	if err != nil {
		return err
	}
	logrus.SetLevel(logLevel)
	return nil
}

func Panic(args ...interface{}) {
	logrus.Panic(args...)
}
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}
func Info(args ...interface{}) {
	logrus.Info(args...)
}
func Trace(args ...interface{}) {
	logrus.Trace(args...)
}

func Debug(args ...interface{}) {
	logrus.Debug(args...)
}
func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}
