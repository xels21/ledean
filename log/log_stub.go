//go:build tinygo
// +build tinygo

package log

import (
	"fmt"
)

func SetLogger(logLevelStr string) error {
	return nil
}

func Panic(args ...interface{}) {
	fmt.Print(args)
}
func Debugf(format string, args ...interface{}) {
	fmt.Print(args)
}
func Info(args ...interface{}) {
	fmt.Print(args)
}
func Trace(args ...interface{}) {
	fmt.Print(args)
}
func Debug(args ...interface{}) {
	fmt.Print(args)
}
func Fatal(args ...interface{}) {
	fmt.Print(args)
}
