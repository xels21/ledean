//go:build !tinygo && !pi
// +build !tinygo,!pi

package pin

import (
	"time"
)

type Pin struct {
}

func NewPin(gpioName string) *Pin {
	self := Pin{}
	return &self
}

func (self *Pin) WaitForEdge(timeout time.Duration) bool {
	return false
}

func (self *Pin) Read() bool {
	return false
}
