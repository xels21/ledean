package led

import (
	"LEDean/led/color"
	"LEDean/led/mode"
	"LEDean/pi/button"
	"LEDean/pi/ws28x"
	"encoding/json"

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
	modeController   *mode.ModeController //[]mode.Mode
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

	self.modeController = mode.NewModeController(self.leds, &self.cUpdate)

	self.registerEvents()
	self.Clear()
	go self.listen()

	return &self
}

func (self *LedController) SwitchModeIndex(index uint8) {
	if self.modeController.GetIndex() == index {
		return
	}

	resume := self.active

	if self.active {
		self.Stop()
	}
	self.modeController.SetIndex(index)
	if resume {
		self.Start()
	}
}

func (self *LedController) GetModeLength() uint8 {
	return self.modeController.GetLength()
}

func (self *LedController) GetModeIndex() uint8 {
	return self.modeController.GetIndex()
}

func (self *LedController) GetModeRef(friendlyName string) (*mode.Mode, error) {
	return self.modeController.GetModeRef(friendlyName)
}

func (self *LedController) GetModeResolver() []string {
	return self.modeController.GetModeResolver()
}

func (self *LedController) GetLeds() []color.RGB {
	return self.leds
}

func (self *LedController) GetLedsJson() []byte {
	msg, err := json.Marshal(self.leds)
	if err != nil {
		msg = []byte(err.Error())
	}
	return msg
}

func (self *LedController) listen() {
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
}

func (self *LedController) Stop() {
	if !self.active {
		return
	}
	log.Trace("stop")
	self.modeController.DeactivateCurrentMode()
	self.Clear()
	self.Render()
	self.active = false
}
func (self *LedController) Start() {
	if self.active {
		return
	}
	log.Trace("start")
	self.modeController.ActivateCurrentMode()
	self.active = true
}

func (self *LedController) NextMode() {
	log.Info("nextMode")
	if !self.active {
		return
	}
	self.modeController.DeactivateCurrentMode()
	self.modeController.NextMode()
	self.modeController.RandomizeCurrentMode()
	self.modeController.ActivateCurrentMode()
}

func (self *LedController) Randomize() {
	log.Info("Randomize")
	if !self.active {
		return
	}
	self.modeController.DeactivateCurrentMode()
	self.modeController.RandomizeCurrentMode()
	self.modeController.ActivateCurrentMode()
}

func (self *LedController) registerEvents() {
	self.piButton.AddCbPressSingle(self.NextMode)
	self.piButton.AddCbPressDouble(self.Randomize)
	self.piButton.AddCbPressLong(self.StartStop)
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
