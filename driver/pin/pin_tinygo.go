//go:build tinygo
// +build tinygo

package pin

import (
	"machine"
	"strconv"
	"time"
)

type Pin struct {
	pin machine.Pin
}

const SAMPLING_RATE_MS = 70

func NewPin(gpioName string) *Pin {
	gpioInt, _ := strconv.Atoi(gpio)
	self := Pin{
		machine.Pin(gpioInt),
	}

	self.pin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	return &self
}

func (self *Pin) WaitForEdgeInfinite(returnVal chan bool) {
	currentState := self.Read()
	for {
		time.Sleep(SAMPLING_RATE_MS * time.Millisecond)
		if currentState != self.Read() {
			returnVal <- true
		}
	}
}

func (self *Pin) WaitForEdge(timeout time.Duration) bool {
	returnVal := make(chan bool, 1)

	go self.WaitForEdgeInfinite(returnVal)
	if timeout != -1 {
		select {
		case <-returnVal:
			return true
		case <-time.After(timeout):
			return false
		}
	}
	return <-returnVal
}

func (self *Pin) Read() bool {
	return self.pin.Get()
}
