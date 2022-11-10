//go:build !tinygo && !pi
// +build !tinygo,!pi

package display

type Display struct {
	DisplayBase
}

func NewDisplay(ledCount int, ledRows int, gpioLedData string, reverseRowsRaw string) *Display {
	self := Display{
		DisplayBase: *NewDisplayBase(ledCount, ledRows, reverseRowsRaw),
	}
	return &self
}

func (self *Display) Render() {
}
