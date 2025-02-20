//go:build !tinygo && !pi
// +build !tinygo,!pi

package display

import "ledean/websocket"

type Display struct {
	DisplayBase
}

func NewDisplay(ledCount int, ledRows int, gpioLedData string, reverseRowsRaw string, fps int, order int, device int, hub *websocket.Hub) Display {
	self := Display{
		DisplayBase: NewDisplayBase(ledCount, ledRows, reverseRowsRaw, fps, order, hub),
	}
	return self
}

func (self *Display) Render() {
}
