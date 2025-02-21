package mode

import (
	"ledean/color"
	"ledean/dbdriver"
	"ledean/display"
	"ledean/json"
	"math"
)

const (
	RunningLedStyleLinear        RunningLedStyle = "linear"
	RunningLedStyleTrigonometric RunningLedStyle = "trigonometric"
)

type RunningLedStyle string

type ModeRunningLed struct {
	ModeSuper
	parameter           ModeRunningLedParameter
	limits              ModeRunningLedLimits
	presets             []ModeRunningLedParameter
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

func NewModeRunningLed(dbdriver *dbdriver.DbDriver, display *display.Display, isRandDeterministic bool) *ModeRunningLed {
	self := ModeRunningLed{
		limits: ModeRunningLedLimits{
			MinRoundTimeMs: 1000,  //1s
			MaxRoundTimeMs: 30000, //30s
			MinBrightness:  0.2,
			MaxBrightness:  1.0,
			MinFadePct:     0.0,
			MaxFadePct:     1.0,
		},
		positionDeg:   0.0,
		ledsRGB:       make([]color.RGB, display.GetRowLedCount()),
		activatedLeds: make([]float64, display.GetRowLedCount()),
	}

	self.presets = self.getPresets()
	self.ModeSuper = *NewModeSuper(dbdriver, display, "ModeRunningLed", RenderTypeDynamic, self.calcDisplay, isRandDeterministic)

	err := dbdriver.Read(self.GetName(), "parameter", &self.parameter)
	if err != nil {
		self.Randomize()
	} else {
		self.postSetParameter()
	}

	return &self
}
func (self *ModeRunningLed) getPresets() []ModeRunningLedParameter {
	return []ModeRunningLedParameter{
		{
			Brightness:  .1,
			RoundTimeMs: 300,
			HueFrom:     259.5,
			HueTo:       74.32,
			FadePct:     0.95,
			Style:       RunningLedStyleLinear,
		},
		{
			Brightness:  .1,
			RoundTimeMs: 400,
			HueFrom:     11.04,
			HueTo:       320.05,
			FadePct:     0.27,
			Style:       RunningLedStyleLinear,
		},
	}
}
func (self *ModeRunningLed) GetParameter() interface{} { return &self.parameter }
func (self *ModeRunningLed) GetLimits() interface{}    { return &self.limits }

func (self *ModeRunningLed) TrySetParameter(b []byte) error {
	var tempPar ModeRunningLedParameter
	err := json.Unmarshal(b, &tempPar)

	if err != nil {
		return err
	}

	self.SetParameter(tempPar)
	return nil
}
func (self *ModeRunningLed) SetParameter(parm ModeRunningLedParameter) {
	self.parameter = parm
	self.dbdriver.Write(self.GetName(), "parameter", self.parameter)
	self.postSetParameter()
}

func (self *ModeRunningLed) postSetParameter() {
	self.hueDistance = math.Abs(self.parameter.HueFrom - self.parameter.HueTo)
	self.hueDistanceFct = 1.0
	if self.parameter.HueFrom > self.parameter.HueTo {
		self.hueDistanceFct = -1.0
	}
	self.positionDegStepSize = 360.0 / (float64(self.parameter.RoundTimeMs) / 1000.0 /*s*/) * (float64(self.display.GetRefreshIntervalNs()) / 1000.0 /*s*/ / 1000.0 /*ms*/ / 1000.0 /*us*/)
	self.darkenStepSize = (1 / self.parameter.FadePct) / (float64(self.parameter.RoundTimeMs) / 1000.0 /*s*/) * (float64(self.display.GetRefreshIntervalNs()) / 1000.0 /*s*/ / 1000.0 /*ms*/ / 1000.0 /*us*/)
	self.lightenStepSize = 2 * self.parameter.Brightness * (float64(self.display.GetRowLedCount())) / (float64(self.parameter.RoundTimeMs) / 1000.0 /*s*/) * (float64(self.display.GetRefreshIntervalNs()) / 1000.0 /*s*/ / 1000.0 /*ms*/ / 1000.0 /*us*/)
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
	case RunningLedStyleTrigonometric:
		activeLedIdx = int(((math.Cos((self.positionDeg+180.0)*math.Pi/180.0) + 1.0) / 2) * float64(len(self.activatedLeds)))
	}
	if activeLedIdx >= len(self.activatedLeds) {
		activeLedIdx = len(self.activatedLeds) - 1 //avoid "index out of range [x] with length x"
	}
	return activeLedIdx
}

func (self *ModeRunningLed) stepForward() {
	self.positionDeg += self.positionDegStepSize
	if self.positionDeg >= 360.0 {
		self.positionDeg -= 360.0
	}
}

func (self *ModeRunningLed) calcDisplay() {
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
}

func (self *ModeRunningLed) AddRunningLed(position float64, speed float64) {

}

func (self *ModeRunningLed) RandomizePreset() {
	self.SetParameter(self.presets[self.rand.Uint32()%uint32(len(self.presets))])
}
func (self *ModeRunningLed) Randomize() {
	self.SetParameter(ModeRunningLedParameter{
		Brightness:  self.rand.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
		FadePct:     self.rand.Float64()*(self.limits.MaxFadePct-self.limits.MinFadePct) + self.limits.MinFadePct,
		HueFrom:     self.rand.Float64() * 360.0,
		HueTo:       self.rand.Float64() * 360.0,
		RoundTimeMs: self.rand.Float64()*float64(self.limits.MaxRoundTimeMs-self.limits.MinRoundTimeMs) + float64(self.limits.MinRoundTimeMs),
		Style:       self.getRandomStyle(),
	})
}

func (self *ModeRunningLed) getRandomStyle() RunningLedStyle {
	styleSwitch := self.rand.Uint32() % 2
	var style RunningLedStyle
	switch styleSwitch {
	case 0:
		style = RunningLedStyleLinear
	case 1:
		style = RunningLedStyleTrigonometric
	}
	return style
}
