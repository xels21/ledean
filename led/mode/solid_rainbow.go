package mode

import (
	"LEDean/led/color"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/sdomino/scribble"
	log "github.com/sirupsen/logrus"
)

type ModeSolidRainbow struct {
	dbDriver      *scribble.Driver
	leds          []color.RGB
	cUpdate       *chan bool
	parameter     ModeSolidRainbowParameter
	limits        ModeSolidRainbowLimits
	refreshRateNs time.Duration
	shouldExit    bool
	hsv           color.HSV
	stepSizeHue   float64
}

type ModeSolidRainbowParameter struct {
	Brightness  float64 `json:"brightness"`
	RoundTimeMs uint32  `json:"roundTimeMs"`
}
type ModeSolidRainbowLimits struct {
	MinRoundTimeMs uint32  `json:"minRoundTimeMs"`
	MaxRoundTimeMs uint32  `json:"maxRoundTimeMs"`
	MinBrightness  float64 `json:"minBrightness"`
	MaxBrightness  float64 `json:"maxBrightness"`
}

func NewModeRainbowSolid(leds []color.RGB, cUpdate *chan bool, dbDriver *scribble.Driver) *ModeSolidRainbow {
	self := ModeSolidRainbow{
		dbDriver: dbDriver,
		leds:     leds,
		cUpdate:  cUpdate,
		limits: ModeSolidRainbowLimits{
			MinRoundTimeMs: 3000,   //2s
			MaxRoundTimeMs: 300000, //5min
			MinBrightness:  0.1,
			MaxBrightness:  1.0,
		},
		refreshRateNs: time.Duration((1000 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/ / 30) * time.Nanosecond), //30fps -> 33.33ms
		// refreshRateNs: time.Duration(100 * time.Millisecond),
		shouldExit: false,
	}

	self.Randomize()

	return &self
}

func (ModeSolidRainbow) GetFriendlyName() string {
	return "ModeSolidRainbow"
}

func (self *ModeSolidRainbow) GetParameterJson() []byte {
	msg, _ := json.Marshal(self.parameter)
	return msg
}

func (self *ModeSolidRainbow) GetLimitsJson() []byte {
	msg, _ := json.Marshal(self.limits)
	return msg
}

func (self *ModeSolidRainbow) SetParameter(parm interface{}) {
	switch parm.(type) {
	case ModeSolidRainbowParameter:
		self.parameter = parm.(ModeSolidRainbowParameter)
		self.dbDriver.Write(self.GetFriendlyName(), "parameter", self.parameter)
		self.hsv.V = self.parameter.Brightness
		self.stepSizeHue = 360.0 / (float64(self.parameter.RoundTimeMs) / 1000) * (float64(self.refreshRateNs) / 1000 / 1000 / 1000)
	}
}

func (self *ModeSolidRainbow) Activate() {
	log.Debugf("start "+self.GetFriendlyName()+" with:\n %s\n", self.GetParameterJson())

	self.shouldExit = false
	go func() {
		rgb := color.RGB{}
		for !self.shouldExit {
			self.hsv.H += self.stepSizeHue
			for self.hsv.H > 360.0 {
				self.hsv.H -= 360.0
			}
			rgb = self.hsv.ToRGB()
			for i := 0; i < len(self.leds); i++ {
				self.leds[i] = rgb
			}
			*self.cUpdate <- true
			time.Sleep(self.refreshRateNs)
		}
	}()
}
func (self *ModeSolidRainbow) Deactivate() {
	self.shouldExit = true
}

func (self *ModeSolidRainbow) Randomize() {
	rand.Seed(time.Now().UnixNano())
	self.SetParameter(ModeSolidRainbowParameter{
		RoundTimeMs: uint32(rand.Float32()*float32(self.limits.MaxRoundTimeMs-self.limits.MinRoundTimeMs)) + self.limits.MinRoundTimeMs,
		Brightness:  rand.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
	})
	self.hsv = color.HSV{
		H: rand.Float64() * 360.0,
		S: 1.0,
		V: self.parameter.Brightness,
	}
}
