package mode

import (
	"LEDean/led/color"
	"encoding/json"
	"math"
	"math/rand"
	"time"

	"github.com/sdomino/scribble"
	log "github.com/sirupsen/logrus"
)

const (
	RunningLedStyleLinear        RunningLedStyle = "linear"
	RunningLedStyleTrigonometric RunningLedStyle = "trigonometric"
)

type RunningLedStyle string

type ModeRunningLed struct {
	dbDriver            *scribble.Driver
	leds                []color.RGB
	cUpdate             *chan bool
	parameter           ModeRunningLedParameter
	limits              ModeRunningLedLimits
	cExit               chan bool
	ticker              time.Ticker
	ledsHsv             []color.HSV
	positionDeg         float64
	positionDegStepSize float64
	brightnessStepSize  float64
}

type ModeRunningLedParameter struct {
	Brightness  float64         `json:"brightness"`
	RoundTimeMs uint32          `json:"roundTimeMs"`
	HueFrom     float64         `json:"hueFrom"`
	HueTo       float64         `json:"hueTo"`
	FadePct     float64         `json:"fadePct"`
	Style       RunningLedStyle `json:"style"`
}
type ModeRunningLedLimits struct {
	MinRoundTimeMs uint32  `json:"minRoundTimeMs"`
	MaxRoundTimeMs uint32  `json:"maxRoundTimeMs"`
	MinBrightness  float64 `json:"minBrightness"`
	MaxBrightness  float64 `json:"maxBrightness"`
	MinFadePct     float64 `json:"minFadePct"`
	MaxFadePct     float64 `json:"maxFadePct"`
}

func NewModeRunningLed(leds []color.RGB, cUpdate *chan bool, dbDriver *scribble.Driver) *ModeRunningLed {
	self := ModeRunningLed{
		dbDriver: dbDriver,
		leds:     leds,
		cUpdate:  cUpdate,
		limits: ModeRunningLedLimits{
			MinRoundTimeMs: 1000,  //1s
			MaxRoundTimeMs: 30000, //30s
			MinBrightness:  0.5,
			MaxBrightness:  1.0,
			MinFadePct:     0.0,
			MaxFadePct:     1.0,
		},
		positionDeg: 0.0,
		ledsHsv:     make([]color.HSV, len(leds)),
		cExit:       make(chan bool, 1),
	}

	self.Randomize()

	return &self
}

func (ModeRunningLed) GetFriendlyName() string {
	return "ModeRunningLed"
}

func (self *ModeRunningLed) GetParameterJson() []byte {
	msg, _ := json.Marshal(self.parameter)
	return msg
}

func (self *ModeRunningLed) GetLimitsJson() []byte {
	msg, _ := json.Marshal(self.limits)
	return msg
}

func (self *ModeRunningLed) SetParameter(parm interface{}) {
	switch parm.(type) {
	case ModeRunningLedParameter:
		self.parameter = parm.(ModeRunningLedParameter)
		self.dbDriver.Write(self.GetFriendlyName(), "parameter", self.parameter)
		self.positionDegStepSize = 360.0 / (float64(self.parameter.RoundTimeMs) / 1000.0 /*s*/) * (float64(REFRESH_RATE_NS) / 1000.0 /*s*/ / 1000.0 /*ms*/ / 1000.0 /*us*/)
		// self.brightnessStepSize = self.positionDegStepSize * self.parameter.FadePct
	}
}

func (self *ModeRunningLed) Activate() {
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

func (self *ModeRunningLed) renderLoop() {

	self.positionDeg += self.positionDegStepSize
	if self.positionDeg >= 360.0 {
		self.positionDeg -= 360.0
	}

	color.HsvArrClear(self.ledsHsv)

	switch self.parameter.Style {
	case RunningLedStyleLinear:
		fadeDirectionInverse := false
		idxNorm := self.positionDeg / 180.0
		if self.positionDeg >= 180 {
			idxNorm = math.Abs(idxNorm - 2.0)
			fadeDirectionInverse = true
		}
		idxf := idxNorm * float64(len(self.ledsHsv))
		idx, idxRest := math.Modf(idxf)
		if fadeDirectionInverse {
			idxRest = 1 - idxRest
		}
		self.ledsHsv[uint(idx)] = color.HSV{H: self.parameter.HueFrom, S: 1.0, V: self.parameter.Brightness * idxRest}
		break
	case RunningLedStyleTrigonometric:
		// idx := (-math.Cos(self.positionDeg) + 1.0) / 2.0
		break
	}
	// self.hsv.H += self.stepSizeHue
	// for self.hsv.H > 360.0 {
	// 	self.hsv.H -= 360.0
	// }
	// rgb = self.hsv.ToRGB()
	// for i := 0; i < len(self.leds); i++ {
	// 	self.leds[i] = rgb
	// }

	for i, _ := range self.ledsHsv {
		self.leds[i] = self.ledsHsv[i].ToRGB()
	}
	*self.cUpdate <- true
}

func (self *ModeRunningLed) Deactivate() {
	self.cExit <- true
}

func (self *ModeRunningLed) Randomize() {
	rand.Seed(time.Now().UnixNano())

	self.SetParameter(ModeRunningLedParameter{
		Brightness:  rand.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
		FadePct:     rand.Float64()*(self.limits.MaxFadePct-self.limits.MinFadePct) + self.limits.MinFadePct,
		HueFrom:     rand.Float64() * 360.0,
		HueTo:       rand.Float64() * 360.0,
		RoundTimeMs: uint32(rand.Float32()*float32(self.limits.MaxRoundTimeMs-self.limits.MinRoundTimeMs)) + self.limits.MinRoundTimeMs,
		Style:       self.getRandomStyle(),
	})
}

func (self *ModeRunningLed) getRandomStyle() RunningLedStyle {
	styleSwitch := rand.Uint32() % 2
	var style RunningLedStyle
	switch styleSwitch {
	case 0:
		style = RunningLedStyleLinear
	case 1:
		style = RunningLedStyleTrigonometric
	}

	style = RunningLedStyleLinear //TODO: remove
	return style
}
