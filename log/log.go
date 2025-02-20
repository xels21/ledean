//go:build !tinygo
// +build !tinygo

package log

import (
	log "github.com/sirupsen/logrus"
)

// var log *logrus.Logger

func SetLogger(logLevelStr string) error {
	// log = &logrus.Logger{
	// 	Formatter: &easy.Formatter{
	// 		TimestampFormat: "15:04",
	// 		LogFormat:       "[%lvl%]: %time% - %msg%\n",
	// 	},
	// }

	// logrus.SetFormatter(&easy.Formatter{
	// TimestampFormat: "2006-01-02 15:04:05",
	// LogFormat:       "[%lvl%]: %time% - %msg%",
	// })
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,

		// FullTimestamp: true,
	})

	logLevel, err := log.ParseLevel(logLevelStr)
	if err != nil {
		return err
	}
	log.SetLevel(logLevel)
	return nil
}

// Trace logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	log.Trace(args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	log.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	log.Warning(args...)
}
func Warningf(format string, args ...interface{}) {
	// log.Warning(args...)
	log.Warningf(format, args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	log.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	log.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
