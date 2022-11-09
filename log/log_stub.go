//go:build tinygo
// +build tinygo

package log

import (
	"log"
)

func SetLogger(logLevelStr string) error {
	return nil
}

func Panic(args ...interface{}) {
	log.Panic(args...)
}
func Debugf(format string, args ...interface{}) {
	// log.Print("[Debug]")
	log.Printf(format, args...)
}
func Info(args ...interface{}) {
	// log.Print("[Info]")
	log.Println(args...)
}
func Trace(args ...interface{}) {
	// log.Print("[Trace]")
	log.Println(args...)
}
func Debug(args ...interface{}) {
	// log.Print("[Debug]")
	log.Println(args...)
}
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}
