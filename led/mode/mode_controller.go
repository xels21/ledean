package mode

import (
	"LEDean/led/color"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type ModeController struct {
	modes            []Mode
	modeSolid        *ModeSolid
	modeRainBowSolid *ModeRainBowSolid
	index            uint8
}

func NewModeController(leds []color.RGB, cUpdate *chan bool) *ModeController {
	modes := ModeController{
		modeSolid:        NewModeSolid(leds, cUpdate),
		modeRainBowSolid: NewModeRainBowSolid(leds, cUpdate),
		index:            0,
	}
	modes.modes = []Mode{modes.modeSolid, modes.modeRainBowSolid}

	return &modes
}

func (self *ModeController) NextMode() {
	self.index = (self.index + 1) % self.GetLength()
	log.Info("Next mode: ", self.index)
}

func (self *ModeController) ActivateCurrentMode() {
	self.modes[self.index].Activate()
}
func (self *ModeController) DeactivateCurrentMode() {
	self.modes[self.index].Deactivate()
}
func (self *ModeController) RandomizeCurrentMode() {
	self.modes[self.index].Randomize()
}

func (self *ModeController) GetModeResolver() []string {
	// m = make(map[int]string)
	modesString := make([]string, 0, 10)
	for _, mode := range self.modes {
		modesString = append(modesString, mode.GetFriendlyName())
	}
	return modesString
}

func (self *ModeController) GetModeRef(friendlyName string) (*Mode, error) {
	// tempModes := self.modeController.GetModes()
	for i, _ := range self.modes {
		if self.modes[i].GetFriendlyName() == friendlyName {
			return &(self.modes[i]), nil
		}
	}
	return nil, fmt.Errorf("mode '%s' not found", friendlyName)
}

func (self *ModeController) GetIndex() uint8 {
	return self.index
}
func (self *ModeController) SetIndex(index uint8) {
	self.index = index
}

func (self *ModeController) GetLength() uint8 {
	return uint8(len(self.modes))
}

func (self *ModeController) GetModes() []Mode {
	return self.modes
}
