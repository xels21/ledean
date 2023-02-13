package mode

import (
	"fmt"
	"ledean/dbdriver"
	"ledean/display"
	"ledean/driver/button"
	"ledean/json"
	"time"

	"ledean/log"
)

const (
	FPS               = 40
	RefreshIntervalNs = time.Duration((1000 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/ / FPS) * time.Nanosecond)
)

type ModeController struct {
	dbdriver              *dbdriver.DbDriver
	display               *display.Display
	button                *button.Button
	active                bool
	modes                 []Mode
	modesIndex            uint8
	modesLength           uint8
	modeSolid             *ModeSolid
	modeSolidRainbow      *ModeSolidRainbow
	modeTransitionRainbow *ModeTransitionRainbow
	modeRunningLed        *ModeRunningLed
	modeEmitter           *ModeEmitter
	modeGradient          *ModeGradient
}

func NewModeController(dbdriver *dbdriver.DbDriver, display *display.Display, button *button.Button) *ModeController {
	self := ModeController{
		dbdriver:              dbdriver,
		display:               display,
		button:                button,
		active:                false,
		modeSolid:             NewModeSolid(dbdriver, display),
		modeSolidRainbow:      NewModeSolidRainbow(dbdriver, display),
		modeTransitionRainbow: NewModeTransitionRainbow(dbdriver, display),
		modeRunningLed:        NewModeRunningLed(dbdriver, display),
		modeEmitter:           NewModeEmitter(dbdriver, display),
		modeGradient:          NewModeGradient(dbdriver, display),
	}
	self.modes = []Mode{self.modeSolid, self.modeSolidRainbow, self.modeTransitionRainbow, self.modeRunningLed, self.modeEmitter, self.modeGradient}
	self.modesLength = uint8(len(self.modes))

	err := dbdriver.Read("modeController", "modesIndex", &self.modesIndex)
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
	self.dbdriver.Write("modeController", "modesIndex", self.modesIndex)
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
	self.button.AddCbPressSingle(self.NextMode)
	self.button.AddCbPressDouble(self.Randomize)
	self.button.AddCbPressLong(self.StartStop)
}
