package mode

import (
	"fmt"
	"ledean/dbdriver"
	"ledean/display"
	"ledean/driver/button"
	"ledean/json"
	"ledean/websocket"
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
	hub                   *websocket.Hub
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

func NewModeController(dbdriver *dbdriver.DbDriver, display *display.Display, button *button.Button, hub *websocket.Hub) *ModeController {
	self := ModeController{
		dbdriver:              dbdriver,
		display:               display,
		button:                button,
		hub:                   hub,
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

	if self.hub != nil {
		go self.socketHandler()
		self.hub.AppendInitClientCb(self.initClientCb)
	}

	return &self
}

func (self *ModeController) socketHandler() {
	for {
		select {
		case cmdModeAction := <-self.hub.CmdModeActionChannel:
			log.Info(cmdModeAction)
			switch cmdModeAction.Action {
			case websocket.CmdModeActionRandomize:
				self.Randomize()
			default:
				log.Info("Unknown mode action: ", cmdModeAction)
			}
		case cmdMode := <-self.hub.CmdModeChannel:
			if cmdMode.Parameter != nil {
				self.handleModeParameterUpdate(cmdMode)
			}
			self.SwitchIndexFriendlyName(cmdMode.Id)
		}
	}
}

func (self *ModeController) handleModeParameterUpdate(cmdMode websocket.CmdMode) {
	switch cmdMode.Id {
	case self.modeEmitter.name:
		var modeEmitterParameter ModeEmitterParameter
		err := json.Unmarshal(cmdMode.Parameter, &modeEmitterParameter)
		if err != nil {
			log.Info("could not parse emitter parameter: ", cmdMode.Parameter)
			return
		}
		self.modeEmitter.SetParameter(modeEmitterParameter)
	case self.modeGradient.name:
		var modeGradientParameter ModeGradientParameter
		err := json.Unmarshal(cmdMode.Parameter, &modeGradientParameter)
		if err != nil {
			log.Info("could not parse gradient parameter: ", cmdMode.Parameter)
			return
		}
		self.modeGradient.SetParameter(modeGradientParameter)
	case self.modeRunningLed.name:
		var modeRunningLedParameter ModeRunningLedParameter
		err := json.Unmarshal(cmdMode.Parameter, &modeRunningLedParameter)
		if err != nil {
			log.Info("could not parse running led parameter: ", cmdMode.Parameter)
			return
		}
		self.modeRunningLed.SetParameter(modeRunningLedParameter)
	case self.modeSolid.name:
		var modeSolidParameter ModeSolidParameter
		err := json.Unmarshal(cmdMode.Parameter, &modeSolidParameter)
		if err != nil {
			log.Info("could not parse solid parameter: ", cmdMode.Parameter)
			return
		}
		self.modeSolid.SetParameter(modeSolidParameter)
		self.Restart()
	case self.modeSolidRainbow.name:
		var modeSolidRainbowParameter ModeSolidRainbowParameter
		err := json.Unmarshal(cmdMode.Parameter, &modeSolidRainbowParameter)
		if err != nil {
			log.Info("could not parse solid rainbow parameter: ", cmdMode.Parameter)
			return
		}
		self.modeSolidRainbow.SetParameter(modeSolidRainbowParameter)
	case self.modeTransitionRainbow.name:
		var modeTransitionRainbowParameter ModeTransitionRainbowParameter
		err := json.Unmarshal(cmdMode.Parameter, &modeTransitionRainbowParameter)
		if err != nil {
			log.Info("could not parse transition rainbow parameter: ", cmdMode.Parameter)
			return
		}
		self.modeTransitionRainbow.SetParameter(modeTransitionRainbowParameter)
	}
	self.BroadcastCurrentMode()
}

func (self *ModeController) SwitchIndexFriendlyName(friendlyName string) {
	self.SwitchIndex(self.GetIndexOf(friendlyName))
}

func (self *ModeController) GetIndexOf(friendlyName string) uint8 {
	for i := range self.modes {
		if self.modes[i].GetName() == friendlyName {
			return uint8(i)
		}
	}
	return uint8(255)
}

func (self *ModeController) GetModeSolid() *ModeSolid {
	return self.modeSolid
}
func (self *ModeController) GetModeSolidRainbow() *ModeSolidRainbow {
	return self.modeSolidRainbow
}
func (self *ModeController) GetModeTransitionRainbow() *ModeTransitionRainbow {
	return self.modeTransitionRainbow
}
func (self *ModeController) GetModeRunningLed() *ModeRunningLed {
	return self.modeRunningLed
}
func (self *ModeController) GetModeEmitter() *ModeEmitter {
	return self.modeEmitter
}
func (self *ModeController) GetModeGradient() *ModeGradient {
	return self.modeGradient
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
	self.BroadcastCurrentMode()
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
	// self.BroadcastCurrentMode()
}

func (self *ModeController) getModeCmd() (websocket.Cmd, error) {
	modeParameter := self.modes[self.modesIndex].GetParameter()
	modeParameterJSON, err := json.Marshal(modeParameter)
	if err != nil {
		return websocket.Cmd{}, err
	}
	cmdModeJSON, err := json.Marshal(websocket.CmdMode{Id: self.modes[self.modesIndex].GetName(), Parameter: modeParameterJSON})
	if err != nil {
		return websocket.Cmd{}, err
	}
	return websocket.Cmd{Command: websocket.CmdModeId, Parameter: cmdModeJSON}, nil
}

func (self *ModeController) BroadcastCurrentMode() {
	if self.hub != nil {
		cmdMode, err := self.getModeCmd()
		if err == nil {
			self.hub.Boradcast(cmdMode)
		}
	}

}

func (self *ModeController) initClientCb(client *websocket.Client) {
	cmdModeResolverJSON, err := json.Marshal(websocket.CmdModeResolver{Modes: self.GetModeResolver()})
	if err != nil {
		log.Info(err)
	} else {
		client.SendCmd(websocket.Cmd{
			Command:   websocket.CmdModeResolverId,
			Parameter: cmdModeResolverJSON})
	}

	for _, mode := range self.modes {
		limitsJSON, err := json.Marshal(mode.GetLimits())
		if err != nil {
			log.Info(err)
			continue
		}
		cmdModeLimitJSON, err := json.Marshal(websocket.CmdModeLimits{
			Id:     mode.GetName(),
			Limits: limitsJSON,
			// Limits: mode.GetLimits(),
		})
		if err != nil {
			log.Info(err)
			continue
		}

		client.SendCmd(websocket.Cmd{
			Command:   websocket.CmdModeLimitsId,
			Parameter: cmdModeLimitJSON})
	}

	cmdMode, err := self.getModeCmd()
	if err != nil {
		log.Info(err)
	} else {
		client.SendCmd(cmdMode)
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
