package mode

/*

import (
	"LEDean/led/color"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/sdomino/scribble"
	log "github.com/sirupsen/logrus"
)

type ModeTransitionRainbow struct {
	dbDriver         *scribble.Driver
	leds             []color.RGB
	cUpdate          *chan bool
	parameter        ModeTransitionRainbowParameter
	limits           ModeTransitionRainbowLimits
	shouldExit       bool
	ledsHsv          []color.HSV
	position         float64 //from 0.0 to 1.0
	positionStepSize float64
	// stepSizeHue float64
}

type ModeTransitionRainbowParameter struct {
	Brightness  float64 `json:"brightness"`
	RoundTimeMs uint32  `json:"roundTimeMs"`
	FromHue     float64 `json:"fromHue"`
	ToHue       float64 `json:"toHue"`
	FadeMs      uint32  `json:"FadeMs"`
}
type ModeTransitionRainbowLimits struct {
	MinRoundTimeMs uint32  `json:"minRoundTimeMs"`
	MaxRoundTimeMs uint32  `json:"maxRoundTimeMs"`
	MinBrightness  float64 `json:"minBrightness"`
	MaxBrightness  float64 `json:"maxBrightness"`
	MinFadeMs      float64 `json:"minFadeMs"`
	MaxFadeMs      float64 `json:"maxFadeMs"`
}

func NewModeTransitionRainbow(leds []color.RGB, cUpdate *chan bool, dbDriver *scribble.Driver) *ModeTransitionRainbow {
	self := ModeTransitionRainbow{
		dbDriver: dbDriver,
		leds:     leds,
		cUpdate:  cUpdate,
		limits: ModeTransitionRainbowLimits{
			MinRoundTimeMs: 3000,   //2s
			MaxRoundTimeMs: 300000, //5min
			MinBrightness:  0.1,
			MaxBrightness:  1.0,
			MinFadeMs:      1000.0,
			MaxFadeMs:      1000.0,
		},
		position:   0.0,
		shouldExit: false,
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
		self.positionStepSize = 360 / (float64(self.parameter.RoundTimeMs) / 1000) * (float64(REFRESH_INTERVAL_NS) / 1000 / 1000 / 1000)
		// self.hsv.V = self.parameter.Brightness
		// self.stepSizeHue = 360.0 / (float64(self.parameter.RoundTimeMs) / 1000) * (float64(REFRESH_INTERVAL_NS) / 1000 / 1000 / 1000)
	}
}

func (self *ModeTransitionRainbow) Activate() {
	log.Debugf("start "+self.GetFriendlyName()+" with:\n %s\n", self.GetParameterJson())

	self.shouldExit = false
	go func() {
		self.ledsHsv = color.RgbArr2HsvArr(self.leds)
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
			time.Sleep(REFRESH_INTERVAL_NS)
		}
	}()
}
func (self *ModeTransitionRainbow) Deactivate() {
	self.shouldExit = true
}

func (self *ModeTransitionRainbow) Randomize() {
	rand.Seed(time.Now().UnixNano())
	self.SetParameter(ModeTransitionRainbowParameter{
		RoundTimeMs: uint32(rand.Float32()*float32(self.limits.MaxRoundTimeMs-self.limits.MinRoundTimeMs)) + self.limits.MinRoundTimeMs,
		Brightness:  rand.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
		FromHue:     rand.Float64(),
		ToHue:       rand.Float64(),
	})
}
*/
