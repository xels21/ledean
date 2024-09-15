//go:build tinygo
// +build tinygo

package display

import (
	"ledean/log"
	"ledean/websocket"
	"machine"
	"strconv"

	"tinygo.org/x/drivers/apa102"
	"tinygo.org/x/drivers/ws2812"
)

type Display struct {
	DisplayBase
	ws      ws2812.Device
	apa     *apa102.Device
	cbWrite func(buf []byte) (n int, err error)
}

func NewDisplay(ledCount int, ledRows int, gpioLedData string, reverseRowsRaw string, fps int, order int, device int, hub *websocket.Hub) *Display {
	self := Display{
		DisplayBase: *NewDisplayBase(ledCount, ledRows, reverseRowsRaw, fps, order, hub),
	}
	switch device {
	case LED_DEVICE_WS2812:
		gpioLedDataInt, _ := strconv.Atoi(gpioLedData)
		pin := machine.Pin(gpioLedDataInt)
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})

		self.ws = ws2812.New(pin)
		self.cbWrite = self.ws.Write
		break
	case LED_DEVICE_APA102:
		// machine.SPI2.Configure(machine.SPIConfig{
		// Frequency: 500000,
		// Mode:      0})
		//
		// self.apa = apa102.New(machine.SPI2)
		sckPin := machine.Pin(4)
		sckPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		sdoPin := machine.Pin(6)
		sdoPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		self.apa = apa102.NewSoftwareSPI(sckPin, sdoPin, 10)
		self.cbWrite = self.apa.Write
		break
	default:
		log.Error("not supported SPI device:")
	}
	return &self
}

func (self *Display) Render() {
	self.leds2Buffer()
	self.cbWrite(self.buffer)
}
