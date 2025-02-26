package mode

import (
	"ledean/color"
	"ledean/dbdriver"
	"ledean/display"
	"ledean/json"
)

type ModeSolidRainbow struct {
	ModeSuper
	parameter   ModeSolidRainbowParameter
	limits      ModeSolidRainbowLimits
	stepSizeHue float64
}

type ModeSolidRainbowParameter struct {
	Brightness  float64   `json:"brightness"`
	RoundTimeMs uint32    `json:"roundTimeMs"`
	Hsv         color.HSV `json:"hsv"`
}
type ModeSolidRainbowLimits struct {
	MinRoundTimeMs uint32  `json:"minRoundTimeMs"`
	MaxRoundTimeMs uint32  `json:"maxRoundTimeMs"`
	MinBrightness  float64 `json:"minBrightness"`
	MaxBrightness  float64 `json:"maxBrightness"`
}

func NewModeSolidRainbow(dbdriver *dbdriver.DbDriver, display *display.Display, isRandDeterministic bool) *ModeSolidRainbow {
	self := ModeSolidRainbow{
		limits: ModeSolidRainbowLimits{
			MinRoundTimeMs: 2000,   //2s
			MaxRoundTimeMs: 300000, //5min
			MinBrightness:  0.01,
			MaxBrightness:  1.0,
		},
	}

	self.ModeSuper = *NewModeSuper(dbdriver, display, "ModeSolidRainbow", RenderTypeDynamic, self.calcDisplay, self.calcDisplayDelta, isRandDeterministic)

	err := dbdriver.Read(self.GetName(), "parameter", &self.parameter)
	if err != nil {
		self.Randomize()
	} else {
		self.postSetParameter()
	}

	return &self
}

func (self *ModeSolidRainbow) GetParameter() interface{} { return &self.parameter }
func (self *ModeSolidRainbow) GetLimits() interface{}    { return &self.limits }

func (self *ModeSolidRainbow) TrySetParameter(b []byte) error {
	var tempPar ModeSolidRainbowParameter
	err := json.Unmarshal(b, &tempPar)

	if err != nil {
		return err
	}

	self.SetParameter(tempPar)
	return nil
}

func (self *ModeSolidRainbow) SetParameter(parm ModeSolidRainbowParameter) {
	self.parameter = parm
	self.dbdriver.Write(self.GetName(), "parameter", self.parameter)
	self.postSetParameter()
}

func (self *ModeSolidRainbow) getStepSizeHue(timeNs float64) float64 {
	// self.stepSizeHue = 360.0 / (float64(self.parameter.RoundTimeMs) / 1000) * (float64(self.display.GetRefreshIntervalNs()) / 1000 / 1000 / 1000)

	return 360.0 / (float64(self.parameter.RoundTimeMs) / 1000 /*s*/) * (timeNs / 1000 /*us*/ / 1000 /*ms*/ / 1000 /*s*/)
}

func (self *ModeSolidRainbow) postSetParameter() {
	self.parameter.Hsv.V = self.parameter.Brightness
	self.stepSizeHue = self.getStepSizeHue(float64(self.display.GetRefreshIntervalNs()))
}

func (self *ModeSolidRainbow) calcDisplayFinal(stepSizeHue float64) {
	self.parameter.Hsv.H += stepSizeHue
	for self.parameter.Hsv.H > 360.0 {
		self.parameter.Hsv.H -= 360.0
	}
	rgb := self.parameter.Hsv.ToRGB()
	self.display.AllSolid(rgb)
}

func (self *ModeSolidRainbow) calcDisplayDelta(deltaTimeNs int64) {
	self.calcDisplayFinal(self.getStepSizeHue(float64(deltaTimeNs)))
}

func (self *ModeSolidRainbow) calcDisplay() {
	self.calcDisplayFinal(self.stepSizeHue)
}

func (self *ModeSolidRainbow) RandomizePreset() {
	self.Randomize()
}
func (self *ModeSolidRainbow) Randomize() {
	self.SetParameter(ModeSolidRainbowParameter{
		RoundTimeMs: uint32(self.rand.Float32()*float32(self.limits.MaxRoundTimeMs-self.limits.MinRoundTimeMs)) + self.limits.MinRoundTimeMs,
		Brightness:  self.rand.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
		Hsv: color.HSV{
			H: self.rand.Float64() * 360.0,
			S: 1.0,
			V: self.parameter.Brightness,
		},
	})
}
