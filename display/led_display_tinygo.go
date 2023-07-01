//go:build tinygo
// +build tinygo

package display

import (
	"machine"
	"strconv"

	"tinygo.org/x/drivers/ws2812"
)

type Display struct {
	DisplayBase
	ws ws2812.Device
}

func NewDisplay(ledCount int, ledRows int, gpioLedData string, reverseRowsRaw string, hub *webserver.Hub) *Display {
	gpioLedDataInt, _ := strconv.Atoi(gpioLedData)
	pin := machine.Pin(gpioLedDataInt)
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	self := Display{
		DisplayBase: *NewDisplayBase(led_count, led_rows, reverse_rows_raw, hub),
		ws:          ws2812.New(pin),
	}

	return &self
}

func (self *Display) Render() {
	self.leds2Buffer()
	self.ws.Write(self.buffer)
}
