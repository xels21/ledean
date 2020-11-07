package led

import (
	"LEDean/pi/button"
	"LEDean/pi/ws28x"
)

type LedController struct {
	led_count        int64
	piWs28xConnector *ws28x.PiWs28xConnector
	piButton         *button.PiButton
	leds             []ColorRGB
	buffer           []byte
}

func NewLedController(led_count int64, piWs28xConnector *ws28x.PiWs28xConnector, piButton *button.PiButton) *LedController {
	var self LedController = LedController{
		led_count:        led_count,
		piWs28xConnector: piWs28xConnector,
		piButton:         piButton,
		leds:             make([]ColorRGB, led_count),
		buffer:           make([]byte, 9*led_count),
	}
	self.registerEvents()

	// self.piButton.CbSinglePress
	// self.piWs28xConnector.
	self.Clear()
	// go func() {
	// 	time.Sleep(100 * time.Millisecond)
	// 	self.Render()
	// 	}()

	return &self
}

func (self *LedController) registerEvents() {
	self.piButton.CbSinglePress = append(self.piButton.CbSinglePress, func() {
		// self.AllSolid(ColorRGB{R: 255, G: 0, B: 0})
		self.AllSolid(ColorRGB{R: 0, G: 255, B: 0})
		// self.AllSolid(ColorRGB{R: 0, G: 0, B: 255})
		self.Render()
	})
	self.piButton.CbLongPress = append(self.piButton.CbLongPress, func() {
		self.Clear()
		self.Render()
	})
}

func (self *LedController) Clear() {
	self.AllSolid(ColorRGB{R: 0, G: 0, B: 0})
}

func (self *LedController) AllSolid(color ColorRGB) {
	for i, _ := range self.leds {
		self.leds[i] = color
	}
}

func (self *LedController) Render() {
	self.leds2Buffer()
	self.piWs28xConnector.Write(self.buffer)
}

func (self *LedController) leds2Buffer() {
	self.buffer = make([]byte, 0, 9*self.led_count)
	for _, led := range self.leds {
		self.buffer = append(self.buffer, led.ToSpi()...)
	}
}
