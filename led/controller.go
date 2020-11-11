package led

import (
	"LEDean/led/color"
	"LEDean/led/mode"
	"LEDean/pi/button"
	"LEDean/pi/ws28x"

	log "github.com/sirupsen/logrus"
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
	modesLength      uint8
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
		modeIndex:        0,
	}

	self.modes = mode.GetAllModes(self.leds, &self.cUpdate)
	self.modesLength = uint8(len(self.modes))
	// self.modes

	self.registerEvents()
	self.Clear()
	go self.run()

	return &self
}

func (self *LedController) GetLeds() []color.RGB {
	return self.leds
}

func (self *LedController) run() {
	for {
		<-self.cUpdate
		self.Render()
	}
}

func (self *LedController) StartStop() {
	// TODO: On/Off
	if self.active {
		self.Stop()
	} else {
		self.Start()
	}
	self.active = !self.active
}

func (self *LedController) Stop() {
	if !self.active {
		return
	}
	log.Trace("stop")
	self.modes[self.modeIndex].Deactivate()
	self.Clear()
	self.Render()
}
func (self *LedController) Start() {
	if self.active {
		return
	}
	log.Trace("start")
	self.modes[self.modeIndex].Activate()
}

func (self *LedController) NextMode() {
	if !self.active {
		return
	}
	self.modes[self.modeIndex].Deactivate()
	self.modeIndex += 1
	if self.modeIndex >= self.modesLength {
		self.modeIndex = 0
	}
	log.Trace("Next mode: ", self.modeIndex)
	self.modes[self.modeIndex].Activate()
}

func (self *LedController) registerEvents() {
	self.piButton.AddCbSinglePress(self.NextMode)
	// self.piButton.CbDoublePress = append(self.piButton.CbDoublePress, func() {
	// TODO: randomize parameter
	// })
	self.piButton.AddCbLongPress(self.StartStop)
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
