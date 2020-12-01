package mode

import (
	"fmt"
	"ledean/display"
	"ledean/pi/button"
	"time"

	"github.com/sdomino/scribble"
	log "github.com/sirupsen/logrus"
)

const (
	REFRESH_INTERVAL_NS = time.Duration((1000 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/ / 30) * time.Nanosecond)
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
	}
	self.modes = []Mode{self.modeSolid, self.modeSolidRainbow, self.modeTransitionRainbow, self.modeRunningLed}
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
		modesString = append(modesString, mode.GetFriendlyName())
	}
	return modesString
}

func (self *ModeController) GetModeRef(friendlyName string) (*Mode, error) {
	// tempModes := self.GetModes()
	for i, _ := range self.modes {
		if self.modes[i].GetFriendlyName() == friendlyName {
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
		self.Stop()
	}
	self.SetIndex(index)
	if resume {
		self.Start()
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
