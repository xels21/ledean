//go:build pi
// +build pi

package display

import (
	"ledean/driver/ws28x"
	"ledean/websocket"
)

type Display struct {
	DisplayBase
	piWs28xConnector *ws28x.PiWs28xConnector
}

func NewDisplay(ledCount int, ledRows int, gpioLedData string, reverseRowsRaw string, fps int, hub *websocket.Hub) *Display {
	self := Display{
		DisplayBase: *NewDisplayBase(ledCount, ledRows, reverseRowsRaw, fps, hub),
	}
	self.piWs28xConnector = ws28x.NewPiWs28xConnector(gpioLedData)
	self.piWs28xConnector.Connect(ledCount)
	return &self
}

func (self *Display) Render() {
	self.leds2Buffer()
	self.piWs28xConnector.Write(self.buffer)
}
