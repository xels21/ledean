//go:build tinygo
// +build tinygo

package display

import (
	"image/color"
	"machine"
	"strconv"
	"time"

	"github.com/aykevl/ledsgo"
	"tinygo.org/x/drivers/ws2812"
)

type Display struct {
	DisplayBase
	pin   machine.Pin
	ws    ws2812.Device
	strip ledsgo.Strip
}

func NewDisplay(led_count int, led_rows int, gpioLedData string, reverse_rows_raw string) *Display {
	gpioLedDataInt, _ := strconv.Atoi(gpioLedData)
	self := Display{
		DisplayBase: *NewDisplayBase(led_count, led_rows, reverse_rows_raw),
		pin:         machine.Pin(gpioLedDataInt),
		strip:       make(ledsgo.Strip, led_count),
	}
	self.pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	self.ws = ws2812.New(self.pin)

	now := time.Now().UnixNano()
	for i := range self.strip {
		self.strip[i] = ledsgo.Color{uint16(now>>15) - uint16(i)<<12, 0xff, 0x44}.Spectrum()
	}

	return &self
}

func (self *Display) Render() {
	for i := range self.strip {
		self.strip[i] = color.RGBA{self.leds[i].R, self.leds[i].G, self.leds[i].B, 0}
	}
	self.ws.WriteColors(self.strip)
}
