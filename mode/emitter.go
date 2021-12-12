package mode

import (
	"encoding/json"
	"ledean/color"
	"ledean/display"
	"math"
	"math/rand"
	"time"

	"github.com/sdomino/scribble"
)

type EmitStyle uint8

const (
	MaxCooldown float64 = 0.2 //20%
)

const (
	EmitStyleDrop  EmitStyle = iota
	EmitStylePulse EmitStyle = iota
	EmitStyleCount EmitStyle = iota
)

type ModeEmit struct {
	HueFrom         float64
	HueTo           float64
	MinLifetimeMs   uint32
	MaxLifetimeMs   uint32
	LifetimeMs      uint32
	PositionPer     float64
	ImpactPer       float64
	ProgressPer     float64
	ProgressPerStep float64
}

func (self *ModeEmit) addPulseToLeds(leds []color.HSV) {
	self.addDropToLeds(leds)
	// startI := int(self.PositionPer * float64(len(leds)))

	// affectedLedsCount := self.ProgressPer * len(leds) * self.ImpactPer
	// impactLedCount := self.ImpactPer * float64(len(leds)) / 2

	// leds[startI].Add(color.HSV{H: self.HueFrom, S: 1.0, V: self.ProgressPer})
}

func (self *ModeEmit) addDropToLeds(leds []color.HSV) {
	startI := int(self.PositionPer * float64(len(leds)))

	affectedLedsCountf := self.ProgressPer * float64(len(leds)) * self.ImpactPer / 2
	for i := 0; i <= int(affectedLedsCountf); i++ {
		rest := math.Min(affectedLedsCountf-float64(i), 1.0)
		hsv := color.HSV{H: self.HueFrom, S: 1.0, V: rest}
		if i == 0 {
			leds[startI+i].Add(hsv)
			continue
		}
		if startI+i < len(leds) {
			leds[startI+i].Add(hsv)
		}
		if startI-i >= 0 {
			leds[startI-i].Add(hsv)
		}
	}
	// impactLedCount := self.ImpactPer * float64(len(leds)) / 2

}

func (self *ModeEmit) stepForward() {
	self.ProgressPer += self.ProgressPerStep
	if self.ProgressPer > 1.0 {
		self.randomize()
	}
}
func (self *ModeEmit) randomize() {
	self.HueFrom = rand.Float64() * 360.0
	self.HueTo = rand.Float64() * 360.0
	self.ImpactPer = rand.Float64()
	self.LifetimeMs = rand.Uint32()%(self.MaxLifetimeMs-self.MinLifetimeMs) + self.MinLifetimeMs
	self.PositionPer = rand.Float64()
	self.ProgressPerStep = 1.0 / (float64(self.LifetimeMs) / 1000) * (float64(RefreshIntervalNs) / 1000 / 1000 / 1000)
	self.ProgressPer = -rand.Float64() * MaxCooldown
}

type ModeEmitter struct {
	ModeSuper
	parameter ModeEmitterParameter
	limits    ModeEmitterLimits
	emits     []ModeEmit
	ledsHSV   []color.HSV
}

type ModeEmitterParameter struct {
	EmitCount uint32    `json:"emitCount"`
	EmitStyle EmitStyle `json:"emitStyle"`
	// RGB        color.RGB `json:"rgb"`
	// Brightness float64   `json:"brightness"`
}

type ModeEmitterLimits struct {
	MinEmitCount      uint32 `json:"minEmitCount"`
	MaxEmitCount      uint32 `json:"maxEmitCount"`
	MinEmitLifetimeMs uint32 `json:"minEmitLifetimeMs"`
	MaxEmitLifetimeMs uint32 `json:"maxEmitLifetimeMs"`
}

func NewModeEmitter(dbDriver *scribble.Driver, display *display.Display) *ModeEmitter {
	self := ModeEmitter{
		limits: ModeEmitterLimits{
			MinEmitCount:      1,
			MaxEmitCount:      5,
			MinEmitLifetimeMs: 500,
			MaxEmitLifetimeMs: 7000,
		},
	}
	self.ModeSuper = *NewModeSuper(dbDriver, display, "ModeEmitter", RenderTypeDynamic, self.calcDisplay)
	self.ledsHSV = make([]color.HSV, display.GetRowLedCount())

	err := dbDriver.Read(self.name, "parameter", &self.parameter)
	if err != nil {
		self.Randomize()
	} else {
		self.postSetParameter()
	}

	return &self
}

func (self *ModeEmitter) GetParameter() interface{} { return &self.parameter }
func (self *ModeEmitter) GetLimits() interface{}    { return &self.limits }

func (self *ModeEmitter) calcDisplay() {
	color.HsvArrClear(self.ledsHSV)
	for i := range self.emits {
		self.emits[i].stepForward()
		if self.emits[i].ProgressPer < 0 {
			continue
		}

		switch self.parameter.EmitStyle {
		case EmitStylePulse:
			self.emits[i].addPulseToLeds(self.ledsHSV)
		case EmitStyleDrop:
			self.emits[i].addDropToLeds(self.ledsHSV)
		}

	}
	self.display.ApplySingleRowHSV(self.ledsHSV)
}

func (self *ModeEmitter) TrySetParameter(b []byte) error {
	var tempPar ModeEmitterParameter
	err := json.Unmarshal(b, &tempPar)

	if err != nil {
		return err
	}

	self.setParameter(tempPar)
	return nil
}

func (self *ModeEmitter) postSetParameter() {
	self.emits = make([]ModeEmit, self.parameter.EmitCount)
	for i := range self.emits {
		self.emits[i].MinLifetimeMs = self.limits.MinEmitLifetimeMs
		self.emits[i].MaxLifetimeMs = self.limits.MaxEmitLifetimeMs
		self.emits[i].randomize()
	}
}

func (self *ModeEmitter) setParameter(parm ModeEmitterParameter) {
	self.parameter = parm
	self.dbDriver.Write(self.name, "parameter", self.parameter)
	self.postSetParameter()
}

func (self *ModeEmitter) Randomize() {
	rand.Seed(time.Now().UnixNano())
	parameter := ModeEmitterParameter{
		EmitCount: rand.Uint32()%(self.limits.MaxEmitCount-self.limits.MinEmitCount) + self.limits.MinEmitCount,
		EmitStyle: EmitStyle(rand.Uint32() % uint32(EmitStyleCount)),
	}
	self.setParameter(parameter)
}
