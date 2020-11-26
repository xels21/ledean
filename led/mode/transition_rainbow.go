package mode

import (
	"LEDean/led/color"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/sdomino/scribble"
	log "github.com/sirupsen/logrus"
)

type ModeTransitionRainbow struct {
	dbDriver            *scribble.Driver
	leds                []color.RGB
	cUpdate             *chan bool
	parameter           ModeTransitionRainbowParameter
	limits              ModeTransitionRainbowLimits
	ledsHsv             []color.HSV
	progressDeg         float64 //from 0.0 to 360.0
	progressDegStepSize float64
	cExit               chan bool

	// stepSizeHue float64
}

type ModeTransitionRainbowParameter struct {
	Brightness  float64 `json:"brightness"`
	Spectrum    float64 `json:"spectrum"`
	RoundTimeMs uint32  `json:"roundTimeMs"`
}
type ModeTransitionRainbowLimits struct {
	MinRoundTimeMs uint32  `json:"minRoundTimeMs"`
	MaxRoundTimeMs uint32  `json:"maxRoundTimeMs"`
	MinBrightness  float64 `json:"minBrightness"`
	MaxBrightness  float64 `json:"maxBrightness"`
	MinSpectrum    float64 `json:"minSpectrum"`
	MaxSpectrum    float64 `json:"maxSpectrum"`
}

func NewModeTransitionRainbow(leds []color.RGB, cUpdate *chan bool, dbDriver *scribble.Driver) *ModeTransitionRainbow {
	self := ModeTransitionRainbow{
		dbDriver: dbDriver,
		leds:     leds,
		cUpdate:  cUpdate,
		limits: ModeTransitionRainbowLimits{
			MinRoundTimeMs: 500,
			MaxRoundTimeMs: 30000,
			MinBrightness:  0.1,
			MaxBrightness:  1.0,
			MinSpectrum:    0.1,
			MaxSpectrum:    2.0,
		},
		progressDeg: 0.0,
		cExit:       make(chan bool, 1),
		ledsHsv:     make([]color.HSV, len(leds)),
	}

	for i := 0; i < len(self.ledsHsv); i++ {
		self.ledsHsv[i] = color.HSV{H: 0.0, S: 1.0, V: 0.0}
	}

	self.Randomize()

	return &self
}

func (ModeTransitionRainbow) GetFriendlyName() string {
	return "ModeTransitionRainbow"
}

func (self *ModeTransitionRainbow) GetParameterJson() []byte {
	msg, _ := json.Marshal(self.parameter)
	return msg
}

func (self *ModeTransitionRainbow) GetLimitsJson() []byte {
	msg, _ := json.Marshal(self.limits)
	return msg
}

func (self *ModeTransitionRainbow) SetParameter(parm interface{}) {
	switch parm.(type) {
	case ModeTransitionRainbowParameter:
		self.parameter = parm.(ModeTransitionRainbowParameter)
		self.dbDriver.Write(self.GetFriendlyName(), "parameter", self.parameter)
		self.progressDegStepSize = 360 / (float64(self.parameter.RoundTimeMs) / 1000) * (float64(REFRESH_INTERVAL_NS) / 1000 / 1000 / 1000)
		for i := 0; i < len(self.ledsHsv); i++ {
			self.ledsHsv[i].H = self.ledsHsv[0].H + float64(i)/float64(len(self.ledsHsv))*self.parameter.Spectrum*360.0
			self.ledsHsv[i].V = self.parameter.Brightness
		}

		// self.hsv.V = self.parameter.Brightness
		// self.stepSizeHue = 360.0 / (float64(self.parameter.RoundTimeMs) / 1000) * (float64(REFRESH_INTERVAL_NS) / 1000 / 1000 / 1000)
	}
}

func (self *ModeTransitionRainbow) Activate() {
	log.Debugf("start "+self.GetFriendlyName()+" with:\n %s\n", self.GetParameterJson())
	ticker := time.NewTicker(REFRESH_INTERVAL_NS)

	go func() {
		for {
			select {
			case <-self.cExit:
				ticker.Stop()
				return
			case <-ticker.C:
				self.renderLoop()
			}
		}
	}()
}

func (self *ModeTransitionRainbow) renderLoop() {
	for i := 0; i < len(self.ledsHsv); i++ {
		self.ledsHsv[i].H += self.progressDegStepSize
		if self.ledsHsv[i].H > 360.0 {
			self.ledsHsv[i].H -= 360.0
		}
		self.leds[i] = self.ledsHsv[i].ToRGB()
	}
	*self.cUpdate <- true
}

func (self *ModeTransitionRainbow) Deactivate() {
	self.cExit <- true
}

func (self *ModeTransitionRainbow) Randomize() {
	rand.Seed(time.Now().UnixNano())
	self.SetParameter(ModeTransitionRainbowParameter{
		RoundTimeMs: uint32(rand.Float32()*float32(self.limits.MaxRoundTimeMs-self.limits.MinRoundTimeMs)) + self.limits.MinRoundTimeMs,
		Brightness:  rand.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
		Spectrum:    rand.Float64()*(self.limits.MaxSpectrum-self.limits.MinSpectrum) + self.limits.MinSpectrum,
	})
}
