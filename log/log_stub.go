//go:build tinygo
// +build tinygo

package log

import (
	"log"
)

func SetLogger(logLevelStr string) error {
	return nil
}

// Trace logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	// log.Trace(args...)
	log.Print(args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	// log.Debug(args...)
	log.Print(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	// log.Debugf(format, args...)
	log.Printf(format, args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	log.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	// log.Info(args...)
	log.Print(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	// log.Warn(args...)
	log.Print(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	// log.Warning(args...)
	log.Print(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	// log.Error(args...)
	log.Print(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	log.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}
