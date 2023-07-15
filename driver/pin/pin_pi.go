//go:build pi
// +build pi

package pin

import (
	"ledean/log"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"time"
)

type Pin struct {
	pin gpio.PinIO
}

func NewPin(gpioName string) *Pin {
	pin := gpioreg.ByName(gpioName)
	if pin == nil {
		log.Fatal("Failed to find pin by name: ", gpioName)
	}
	self := Pin{pin: pin}
	// var p PinIn
	if err := self.pin.In(gpio.PullDown, gpio.BothEdges); err != nil {
		log.Fatal("Could not init pin:", err)
	}

	return &self
}

func (self *Pin) WaitForEdge(timeout time.Duration) bool {
	return self.pin.WaitForEdge(timeout)
}

func (self *Pin) Read() bool {
	return self.pin.Read() == true
}
