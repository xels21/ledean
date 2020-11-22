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
	dbDriver    *scribble.Driver
	leds        []color.RGB
	cUpdate     *chan bool
	parameter   ModeSolidRainbowParameter
	limits      ModeSolidRainbowLimits
	cExit       chan bool
	hsv         color.HSV
	stepSizeHue float64
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

func NewModeSolidRainbow(leds []color.RGB, cUpdate *chan bool, dbDriver *scribble.Driver) *ModeSolidRainbow {
	self := ModeSolidRainbow{
		dbDriver: dbDriver,
		leds:     leds,
		cUpdate:  cUpdate,
		limits: ModeSolidRainbowLimits{
			MinRoundTimeMs: 2000,   //2s
			MaxRoundTimeMs: 300000, //5min
			MinBrightness:  0.1,
			MaxBrightness:  1.0,
		},
		cExit: make(chan bool, 1),
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
		self.stepSizeHue = 360.0 / (float64(self.parameter.RoundTimeMs) / 1000) * (float64(REFRESH_RATE_NS) / 1000 / 1000 / 1000)
	}
}

func (self *ModeSolidRainbow) Activate() {
	log.Debugf("start "+self.GetFriendlyName()+" with:\n %s\n", self.GetParameterJson())
	ticker := time.NewTicker(REFRESH_RATE_NS)

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
func (self *ModeSolidRainbow) renderLoop() {
	rgb := color.RGB{}
	self.hsv.H += self.stepSizeHue
	for self.hsv.H > 360.0 {
		self.hsv.H -= 360.0
	}
	rgb = self.hsv.ToRGB()
	for i := 0; i < len(self.leds); i++ {
		self.leds[i] = rgb
	}
	*self.cUpdate <- true
}

func (self *ModeSolidRainbow) Deactivate() {
	self.cExit <- true
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
