//go:build pi
// +build pi

package display

type Display struct {
	DisplayBase
	piWs28xConnector *ws28x.PiWs28xConnector
}

func NewDisplay(led_count int, led_rows int, gpioLedData string, reverse_rows_raw string) *Display {
	self := DisplayBase{
		DisplayBase: *NewDisplayBase(led_count, led_rows, reverse_rows_raw),
	}
	self.piWs28xConnector = ws28x.NewPiWs28xConnector(gpioLedData)
	self.piWs28xConnector.Connect(led_count)
	return &self
}

func (self *Display) Render() {
	self.leds2Buffer()
	self.piWs28xConnector.Write(self.buffer)
}
