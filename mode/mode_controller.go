package mode

import (
	"encoding/json"
	"fmt"
	"ledean/display"
	"ledean/pi/button"
	"time"

	"github.com/sdomino/scribble"
	log "github.com/sirupsen/logrus"
)

const (
	FPS30             = time.Duration((1000 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/ / 30) * time.Nanosecond)
	FPS40             = time.Duration((1000 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/ / 40) * time.Nanosecond)
	FPS50             = time.Duration((1000 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/ / 50) * time.Nanosecond)
	FPS60             = time.Duration((1000 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/ / 60) * time.Nanosecond)
	FPS70             = time.Duration((1000 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/ / 70) * time.Nanosecond)
	FPS80             = time.Duration((1000 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/ / 80) * time.Nanosecond)
	FPS90             = time.Duration((1000 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/ / 90) * time.Nanosecond)
	FPS100            = time.Duration((1000 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/ / 100) * time.Nanosecond)
	FPS110            = time.Duration((1000 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/ / 110) * time.Nanosecond)
	FPS120            = time.Duration((1000 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/ / 120) * time.Nanosecond)
	RefreshIntervalNs = FPS40
)

type ModeController struct {
	dbDriver              *scribble.Driver
	display               *display.Display
	piButton              *button.PiButton
	active                bool
	modes                 []Mode
	modesIndex            uint8
	modesLength           uint8
	modeSolid             *ModeSolid
	modeSolidRainbow      *ModeSolidRainbow
	modeTransitionRainbow *ModeTransitionRainbow
	modeRunningLed        *ModeRunningLed
	modeEmitter           *ModeEmitter
}

func NewModeController(dbDriver *scribble.Driver, display *display.Display, piButton *button.PiButton) *ModeController {
	self := ModeController{
		dbDriver:              dbDriver,
		display:               display,
		piButton:              piButton,
		active:                false,
		modeSolid:             NewModeSolid(dbDriver, display),
		modeSolidRainbow:      NewModeSolidRainbow(dbDriver, display),
		modeTransitionRainbow: NewModeTransitionRainbow(dbDriver, display),
		modeRunningLed:        NewModeRunningLed(dbDriver, display),
		modeEmitter:           NewModeEmitter(dbDriver, display),
	}
	self.modes = []Mode{self.modeSolid, self.modeSolidRainbow, self.modeTransitionRainbow, self.modeRunningLed, self.modeEmitter}
	self.modesLength = uint8(len(self.modes))

	err := dbDriver.Read("modeController", "modesIndex", &self.modesIndex)
	if err != nil {
		self.SetIndex(0)
	}

	self.registerEvents()

	return &self
}

func (self *ModeController) NextMode() {
	log.Info("nextMode")
	if self.active {
		self.DeactivateCurrentMode()
	}
	self.SetIndex((self.modesIndex + 1) % self.modesLength)
	// self.RandomizeCurrentMode()
	if self.active {
		self.ActivateCurrentMode()
	}
}

func (self *ModeController) ActivateCurrentMode() {
	parm, _ := json.Marshal(self.modes[self.modesIndex].GetParameter())
	log.Debugf("Start: `%s` with parameter: `%s`", self.modes[self.modesIndex].GetName(), parm)
	self.modes[self.modesIndex].Activate()
}
func (self *ModeController) DeactivateCurrentMode() {
	self.modes[self.modesIndex].Deactivate()
}
func (self *ModeController) RandomizeCurrentMode() {
	self.modes[self.modesIndex].Randomize()
}

func (self *ModeController) GetModeResolver() []string {
	// m = make(map[int]string)
	modesString := make([]string, 0, self.modesLength)
	for _, mode := range self.modes {
		modesString = append(modesString, mode.GetName())
	}
	return modesString
}

func (self *ModeController) GetModeRef(friendlyName string) (*Mode, error) {
	// tempModes := self.GetModes()
	for i := range self.modes {
		if self.modes[i].GetName() == friendlyName {
			return &(self.modes[i]), nil
		}
	}
	return nil, fmt.Errorf("mode '%s' not found", friendlyName)
}

func (self *ModeController) SwitchIndex(index uint8) {
	if self.modesIndex == index {
		return
	}

	resume := self.active

	if self.active {
		self.DeactivateCurrentMode()
	}
	self.SetIndex(index)
	if resume {
		self.ActivateCurrentMode()
	}
}

func (self *ModeController) GetIndex() uint8 {
	return self.modesIndex
}
func (self *ModeController) SetIndex(modesIndex uint8) {
	self.modesIndex = modesIndex
	self.dbDriver.Write("modeController", "modesIndex", self.modesIndex)
	log.Info("Current mode: ", self.modesIndex)
}

func (self *ModeController) GetLength() uint8 {
	return self.modesLength
}

func (self *ModeController) GetModes() []Mode {
	return self.modes
}

func (self *ModeController) StartStop() {
	// TODO: On/Off
	if self.active {
		self.Stop()
	} else {
		self.Start()
	}
}

func (self *ModeController) Start() {
	if !self.active {
		log.Trace("start")
		self.ActivateCurrentMode()
		self.active = true
	}
}
func (self *ModeController) Stop() {
	if self.active {
		log.Trace("stop")
		self.DeactivateCurrentMode()
		self.display.Clear()
		self.display.Render()
		self.active = false
	}
}
func (self *ModeController) Restart() {
	if self.active {
		log.Trace("restart")
		self.DeactivateCurrentMode()
		self.ActivateCurrentMode()
	}
}

func (self *ModeController) Randomize() {
	log.Info("Randomize")
	if self.active {
		self.DeactivateCurrentMode()
	}
	self.RandomizeCurrentMode()
	if self.active {
		self.ActivateCurrentMode()
	}
}

func (self *ModeController) registerEvents() {
	self.piButton.AddCbPressSingle(self.NextMode)
	self.piButton.AddCbPressDouble(self.Randomize)
	self.piButton.AddCbPressLong(self.StartStop)
}
