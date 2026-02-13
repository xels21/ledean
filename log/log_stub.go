//go:build tinygo
// +build tinygo

package log

import (
	"fmt"
	// "log" -> not wokring with uart due to heap allocations, so we use fmt.Print instead and disable log output via SetLogger
)

func SetLogger(logLevelStr string) error {
	return nil
}

// Trace logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	// log.Trace(args...)
	if false {
		fmt.Println(args...)
	}
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	// log.Debug(args...)
	fmt.Println(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	// log.Debugf(format, args...)
	fmt.Printf(format, args...)
	fmt.Println()
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	fmt.Println(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	// log.Info(args...)
	fmt.Println(args...)
}

// Info logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	// log.Debugf(format, args...)
	fmt.Printf(format, args...)
	fmt.Println()
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	// log.Warn(args...)
	fmt.Println(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	// log.Warning(args...)
	fmt.Println(args...)
}

func Warningf(format string, args ...interface{}) {
	// log.Warning(args...)
	fmt.Printf(format, args...)
	fmt.Println()
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	// log.Error(args...)
	fmt.Println(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	fmt.Println(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	fmt.Println(args...)
}

func Fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Println()
}
