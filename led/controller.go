package led

import (
	"LEDean/led/color"
	"LEDean/led/mode"
	"LEDean/pi/button"
	"LEDean/pi/ws28x"
)

type LedController struct {
	cUpdate          chan bool
	led_count        int64
	piWs28xConnector *ws28x.PiWs28xConnector
	piButton         *button.PiButton
	leds             []color.RGB
	buffer           []byte
	active           bool
	modes            []mode.Mode
	modeIndex        uint8
	modeLength       uint8
}

func NewLedController(led_count int64, piWs28xConnector *ws28x.PiWs28xConnector, piButton *button.PiButton) *LedController {
	var self LedController = LedController{
		cUpdate:          make(chan bool, 1),
		led_count:        led_count,
		piWs28xConnector: piWs28xConnector,
		piButton:         piButton,
		leds:             make([]color.RGB, led_count),
		buffer:           make([]byte, 9*led_count),
		active:           false,
	}

	self.registerEvents()
	self.Clear()

	return &self
}

func (self *LedController) run() {
	for {
		<-self.cUpdate
		self.Render()
	}
}

func (self *LedController) registerEvents() {
	self.piButton.AddCbSinglePress(func() {
		// TODO: switch mode
		// self.AllSolid(RGB{R: 255, G: 0, B: 0})
		self.AllSolid(color.RGB{R: 0, G: 255, B: 0})
		// self.AllSolid(RGB{R: 0, G: 0, B: 255})
		self.Render()
	})
	// self.piButton.CbDoublePress = append(self.piButton.CbDoublePress, func() {
	// TODO: randomize parameter
	// })
	self.piButton.AddCbLongPress(func() {
		// TODO: On/Off
		if self.active {
			self.modes[self.modeIndex].Deactivate()
			self.Clear()
			self.Render()
		} else {
			self.modes[self.modeIndex].Activate()
		}
		self.active = !self.active
	})
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
