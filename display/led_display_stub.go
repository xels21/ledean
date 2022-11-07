//go:build !tinygo && !pi
// +build !tinygo,!pi

package display

type Display struct {
	DisplayBase
}

func NewDisplay(led_count int, led_rows int, gpioLedData string, reverse_rows_raw string) *Display {
	self := Display{
		DisplayBase: *NewDisplayBase(led_count, led_rows, reverse_rows_raw),
	}
	return &self
}

func (self *Display) Render() {
}
