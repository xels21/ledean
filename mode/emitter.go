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
	pParameter      *ModeEmitterParameter
	HueFrom         float64
	HueTo           float64
	Brightness      float64
	LifetimeMs      uint32
	PositionPer     float64
	ImpactPer       float64
	ProgressPer     float64
	ProgressPerStep float64
}

func (self *ModeEmit) addPulseToLeds(leds []color.HSV) {
	startI := int(self.PositionPer * float64(len(leds)))
	prog := ((math.Cos(math.Pi+self.ProgressPer*2*math.Pi) + 1) / 2)

	affectedLedsCountf := prog * float64(len(leds)) * self.ImpactPer / 2
	h := self.HueFrom + (self.HueTo-self.HueFrom)*prog
	if h < 0.0 {
		h += 360.0
	}
	if h > 360.0 {
		h -= 360.0
	}

	for i := 0; i <= int(affectedLedsCountf); i++ {
		rest := math.Min(affectedLedsCountf-float64(i), 1.0)
		hsv := color.HSV{H: h, S: 1.0, V: rest * self.Brightness}
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
	// self.addDropToLeds(leds)
}

func (self *ModeEmit) addDropToLeds(leds []color.HSV) {
	self.addPulseToLeds(leds)
	return
	startI := int(self.PositionPer * float64(len(leds)))

	affectedLedsCountf := self.ProgressPer * float64(len(leds)) * self.ImpactPer / 2
	for i := 0; i <= int(affectedLedsCountf); i++ {
		rest := math.Min(affectedLedsCountf-float64(i), 1.0)
		hsv := color.HSV{H: self.HueFrom, S: 1.0, V: rest * self.Brightness}
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
	self.HueTo = self.HueFrom + ((rand.Float64() - 0.5) * 360.0 * 0.5)
	self.Brightness = self.pParameter.MinBrightness + ((self.pParameter.MaxBrightness - self.pParameter.MinBrightness) * rand.Float64())
	self.ImpactPer = rand.Float64()
	if self.pParameter.MinEmitLifetimeMs == self.pParameter.MaxEmitLifetimeMs {
		self.LifetimeMs = self.pParameter.MinEmitLifetimeMs
	} else {
		self.LifetimeMs = rand.Uint32()%(self.pParameter.MaxEmitLifetimeMs-self.pParameter.MinEmitLifetimeMs) + self.pParameter.MinEmitLifetimeMs
	}
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
	EmitCount         uint8     `json:"emitCount"`
	EmitStyle         EmitStyle `json:"emitStyle"`
	MinBrightness     float64   `json:"minBrightness"`
	MaxBrightness     float64   `json:"maxBrightness"`
	MinEmitLifetimeMs uint32    `json:"minEmitLifetimeMs"`
	MaxEmitLifetimeMs uint32    `json:"maxEmitLifetimeMs"`
}

type ModeEmitterLimits struct {
	MinEmitCount      uint8   `json:"minEmitCount"`
	MaxEmitCount      uint8   `json:"maxEmitCount"`
	MinEmitLifetimeMs uint32  `json:"minEmitLifetimeMs"`
	MaxEmitLifetimeMs uint32  `json:"maxEmitLifetimeMs"`
	MinBrightness     float64 `json:"minBrightness"`
	MaxBrightness     float64 `json:"maxBrightness"`
}

func NewModeEmitter(dbDriver *scribble.Driver, display *display.Display) *ModeEmitter {
	self := ModeEmitter{
		limits: ModeEmitterLimits{
			MinEmitCount:      1,
			MaxEmitCount:      5,
			MinEmitLifetimeMs: 500,
			MaxEmitLifetimeMs: 7000,
			MinBrightness:     0.05,
			MaxBrightness:     1.0,
		},
	}
	self.ModeSuper = *NewModeSuper(dbDriver, display, "ModeEmitter", RenderTypeDynamic, self.calcDisplay)
	self.ledsHSV = make([]color.HSV, display.GetRowLedCount())
	self.emits = make([]ModeEmit, self.limits.MaxEmitCount)
	for i := uint8(0); i < self.limits.MaxEmitCount; i++ {
		self.emits[i].pParameter = &self.parameter
	}

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
	for i := uint8(0); i < self.parameter.EmitCount; i++ {
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
	for i := uint8(0); i < self.parameter.EmitCount; i++ {
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
	minBrightness := self.limits.MinBrightness + (rand.Float64() * (self.limits.MaxBrightness - self.limits.MinBrightness))
	minEmitLifetimeMs := self.limits.MinEmitLifetimeMs + (rand.Uint32() % (self.limits.MaxEmitLifetimeMs - self.limits.MinEmitLifetimeMs))
	parameter := ModeEmitterParameter{
		EmitCount:         uint8(rand.Uint32())%(self.limits.MaxEmitCount-self.limits.MinEmitCount+1) + self.limits.MinEmitCount,
		EmitStyle:         EmitStyle(rand.Uint32() % uint32(EmitStyleCount)),
		MinBrightness:     minBrightness,
		MaxBrightness:     minBrightness + (rand.Float64() * (self.limits.MaxBrightness - minBrightness)),
		MinEmitLifetimeMs: minEmitLifetimeMs,
		MaxEmitLifetimeMs: minEmitLifetimeMs + (rand.Uint32() % (self.limits.MaxEmitLifetimeMs - minEmitLifetimeMs)),
	}
	self.setParameter(parameter)
}
