package mode

import (
	"encoding/json"
	"ledean/color"
	"ledean/display"
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
	display             *display.Display
	parameter           ModeRunningLedParameter
	limits              ModeRunningLedLimits
	cExit               chan bool
	ledsRGB             []color.RGB
	activatedLeds       []float64
	positionDeg         float64
	positionDegStepSize float64
	darkenStepSize      float64
	lightenStepSize     float64
	hueDistance         float64
	hueDistanceFct      float64
}

type ModeRunningLedParameter struct {
	Brightness  float64         `json:"brightness"`
	RoundTimeMs float64         `json:"roundTimeMs"`
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

func NewModeRunningLed(dbDriver *scribble.Driver, display *display.Display) *ModeRunningLed {
	self := ModeRunningLed{
		dbDriver: dbDriver,
		display:  display,
		limits: ModeRunningLedLimits{
			MinRoundTimeMs: 1000,  //1s
			MaxRoundTimeMs: 30000, //30s
			MinBrightness:  0.3,
			MaxBrightness:  1.0,
			MinFadePct:     0.0,
			MaxFadePct:     1.0,
		},
		positionDeg:   0.0,
		ledsRGB:       make([]color.RGB, display.GetRowLedCount()),
		activatedLeds: make([]float64, display.GetRowLedCount()),
		cExit:         make(chan bool, 1),
	}

	err := dbDriver.Read(self.GetFriendlyName(), "parameter", &self.parameter)
	if err != nil {
		self.Randomize()
	} else {
		self.postSetParameter()
	}

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
		self.postSetParameter()
	}
}

func (self *ModeRunningLed) postSetParameter() {
	self.hueDistance = math.Abs(self.parameter.HueFrom - self.parameter.HueTo)
	self.hueDistanceFct = 1.0
	if self.parameter.HueFrom > self.parameter.HueTo {
		self.hueDistanceFct = -1.0
	}
	self.positionDegStepSize = 360.0 / (float64(self.parameter.RoundTimeMs) / 1000.0 /*s*/) * (float64(FPS80) / 1000.0 /*s*/ / 1000.0 /*ms*/ / 1000.0 /*us*/)
	self.darkenStepSize = (1 / self.parameter.FadePct) / (float64(self.parameter.RoundTimeMs) / 1000.0 /*s*/) * (float64(FPS80) / 1000.0 /*s*/ / 1000.0 /*ms*/ / 1000.0 /*us*/)
	self.lightenStepSize = 2 * self.parameter.Brightness * (float64(self.display.GetRowLedCount())) / (float64(self.parameter.RoundTimeMs) / 1000.0 /*s*/) * (float64(FPS80) / 1000.0 /*s*/ / 1000.0 /*ms*/ / 1000.0 /*us*/)
}

func (self *ModeRunningLed) Activate() {
	log.Debugf("start "+self.GetFriendlyName()+" with:\n %s\n", self.GetParameterJson())
	ticker := time.NewTicker(FPS80)

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

func (self *ModeRunningLed) darkenLeds() {
	for i := 0; i < len(self.activatedLeds); i++ {
		if self.activatedLeds[i] != 0.0 {
			if self.activatedLeds[i] <= self.darkenStepSize {
				self.activatedLeds[i] = 0.0
			} else {
				self.activatedLeds[i] -= self.darkenStepSize
			}
		}
	}
}

func (self *ModeRunningLed) getActiveLedIdx() int {
	var activeLedIdx int
	switch self.parameter.Style {
	case RunningLedStyleLinear:
		position := self.positionDeg / 180.0
		if position > 1.0 {
			position = 2.0 - position
		}
		activeLedIdx = int(position * float64(len(self.activatedLeds)))
		break
	case RunningLedStyleTrigonometric:
		activeLedIdx = int(((math.Cos((self.positionDeg+180.0)*math.Pi/180.0) + 1.0) / 2) * float64(len(self.activatedLeds)))
		break
	}
	return activeLedIdx
}

func (self *ModeRunningLed) stepForward() {
	self.positionDeg += self.positionDegStepSize
	if self.positionDeg >= 360.0 {
		self.positionDeg -= 360.0
	}
}

func (self *ModeRunningLed) renderLoop() {
	self.stepForward()
	activeLedIdx := self.getActiveLedIdx()
	//prevent darken active led while lighning it up
	self.activatedLeds[activeLedIdx] += self.darkenStepSize
	self.darkenLeds()

	self.activatedLeds[activeLedIdx] += self.lightenStepSize
	if self.activatedLeds[activeLedIdx] > 1.0 {
		self.activatedLeds[activeLedIdx] = 1.0
	}

	c := color.HSV{H: 0.0, S: 1.0, V: 0.0}
	for i, activatedLed := range self.activatedLeds {
		// if i == activeLedIdx {
		// c.H = self.parameter.HueTo
		// } else {
		c.H = self.parameter.HueFrom + (self.hueDistanceFct * self.hueDistance * activatedLed)
		// }
		c.V = activatedLed
		self.ledsRGB[i] = c.ToRGB()
	}

	self.display.ApplySingleRowRGB(self.ledsRGB)
	self.display.Render()
}

func (self *ModeRunningLed) AddRunningLed(position float64, speed float64) {

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
		RoundTimeMs: rand.Float64()*float64(self.limits.MaxRoundTimeMs-self.limits.MinRoundTimeMs) + float64(self.limits.MinRoundTimeMs),
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
	return style
}