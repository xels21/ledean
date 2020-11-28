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
	led_rows            int
	led_per_row         int
	cUpdate             *chan bool
	parameter           ModeTransitionRainbowParameter
	limits              ModeTransitionRainbowLimits
	singleRowledsHSV    []color.HSV
	progressDeg         float64 //from 0.0 to 360.0
	progressDegStepSize float64
	cExit               chan bool

	// stepSizeHue float64
}

type ModeTransitionRainbowParameter struct {
	Brightness  float64 `json:"brightness"`
	Spectrum    float64 `json:"spectrum"`
	RoundTimeMs uint32  `json:"roundTimeMs"`
	Reverse     bool    `json:"reverse"`
}
type ModeTransitionRainbowLimits struct {
	MinRoundTimeMs uint32  `json:"minRoundTimeMs"`
	MaxRoundTimeMs uint32  `json:"maxRoundTimeMs"`
	MinBrightness  float64 `json:"minBrightness"`
	MaxBrightness  float64 `json:"maxBrightness"`
	MinSpectrum    float64 `json:"minSpectrum"`
	MaxSpectrum    float64 `json:"maxSpectrum"`
}

func NewModeTransitionRainbow(dbDriver *scribble.Driver, cUpdate *chan bool, leds []color.RGB, led_rows int) *ModeTransitionRainbow {
	self := ModeTransitionRainbow{
		dbDriver:    dbDriver,
		leds:        leds,
		led_rows:    led_rows,
		led_per_row: len(leds) / led_rows,
		cUpdate:     cUpdate,
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
	}

	self.singleRowledsHSV = make([]color.HSV, self.led_per_row)
	for i := 0; i < self.led_per_row; i++ {
		self.singleRowledsHSV[i] = color.HSV{H: 0.0, S: 1.0, V: 0.0}
	}

	err := dbDriver.Read(self.GetFriendlyName(), "parameter", &self.parameter)
	if err != nil {
		self.Randomize()
	} else {
		self.postSetParameter()
	}

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
		self.postSetParameter()

	}
}

func (self *ModeTransitionRainbow) postSetParameter() {
	self.progressDegStepSize = 360 / (float64(self.parameter.RoundTimeMs) / 1000) * (float64(REFRESH_INTERVAL_NS) / 1000 / 1000 / 1000)
	for ri := 0; ri < self.led_per_row; ri++ {
		self.singleRowledsHSV[ri].H = self.singleRowledsHSV[0].H + float64(ri)/float64(self.led_per_row)*self.parameter.Spectrum*360.0
		self.singleRowledsHSV[ri].V = self.parameter.Brightness
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
	for ri := 0; ri < len(self.leds)/self.led_rows; ri++ {
		if !self.parameter.Reverse {
			self.singleRowledsHSV[ri].H += self.progressDegStepSize
			if self.singleRowledsHSV[ri].H > 360.0 {
				self.singleRowledsHSV[ri].H -= 360.0
			}
		} else {
			self.singleRowledsHSV[ri].H -= self.progressDegStepSize
			if self.singleRowledsHSV[ri].H < 0.0 {
				self.singleRowledsHSV[ri].H += 360.0
			}
		}
	}

	for r := 0; r < self.led_rows; r++ {
		for ri := 0; ri < len(self.leds)/self.led_rows; ri++ {
			i := r*self.led_per_row + ri
			self.leds[i] = self.singleRowledsHSV[ri].ToRGB()
		}
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
		Reverse:     rand.Int()%2 == 1,
	})
}
