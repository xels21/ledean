package led

import (
	"LEDean/led/color"
	"LEDean/led/mode"
	"LEDean/pi/button"
	"LEDean/pi/ws28x"
	"encoding/json"

	"github.com/sdomino/scribble"
	log "github.com/sirupsen/logrus"
)

type LedController struct {
	cUpdate          chan bool
	led_rows         int
	piWs28xConnector *ws28x.PiWs28xConnector
	piButton         *button.PiButton
	leds             []color.RGB
	buffer           []byte
	active           bool
	modeController   *mode.ModeController //[]mode.Mode
}

func NewLedController(led_count int, led_rows int, piWs28xConnector *ws28x.PiWs28xConnector, piButton *button.PiButton, dbDriver *scribble.Driver) *LedController {
	var self LedController = LedController{
		cUpdate:          make(chan bool, 1),
		led_rows:         led_rows,
		piWs28xConnector: piWs28xConnector,
		piButton:         piButton,
		leds:             make([]color.RGB, led_count),
		buffer:           make([]byte, 9*led_count),
		active:           false,
	}

	self.modeController = mode.NewModeController(dbDriver, &self.cUpdate, self.leds, self.led_rows)

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
func (self *LedController) GetLedCount() int {
	return len(self.leds)
}
func (self *LedController) GetLedRows() int {
	return self.led_rows
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

func (self *LedController) Start() {
	if !self.active {
		log.Trace("start")
		self.modeController.ActivateCurrentMode()
		self.active = true
	}
}
func (self *LedController) Stop() {
	if self.active {
		log.Trace("stop")
		self.modeController.DeactivateCurrentMode()
		self.Clear()
		self.Render()
		self.active = false
	}
}
func (self *LedController) Restart() {
	if self.active {
		log.Trace("restart")
		self.modeController.DeactivateCurrentMode()
		self.modeController.ActivateCurrentMode()
	}
}

func (self *LedController) NextMode() {
	log.Info("nextMode")
	if self.active {
		self.modeController.DeactivateCurrentMode()
	}
	self.modeController.NextMode()
	// self.modeController.RandomizeCurrentMode()
	if self.active {
		self.modeController.ActivateCurrentMode()
	}
}

func (self *LedController) Randomize() {
	log.Info("Randomize")
	if self.active {
		self.modeController.DeactivateCurrentMode()
	}
	self.modeController.RandomizeCurrentMode()
	if self.active {
		self.modeController.ActivateCurrentMode()
	}
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
	self.buffer = make([]byte, 0, 9*len(self.leds))
	for _, led := range self.leds {
		self.buffer = append(self.buffer, led.ToSpi()...)
	}
}
