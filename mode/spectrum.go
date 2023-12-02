package mode

import (
	"ledean/color"
	"ledean/dbdriver"
	"ledean/display"
	"ledean/json"
	"math"
	"math/rand"
	"time"
)

type ModeSpectrumPosition struct {
	factor            float64
	factorPercent     float64
	factorPercentStep float64
	offset            float64
	offsetPercent     float64
	offsetPercentStep float64
	parm              *ModeSpectrumParameterPosition
}

func (self *ModeSpectrumPosition) StepForward() {
	self.factorPercent += self.factorPercentStep
	for self.factorPercent >= 1.0 {
		self.factorPercent -= 1.0
	}
	self.factor = ((self.parm.FacTo - self.parm.FacFrom) * (math.Sin(self.factorPercent*2*math.Pi)*.5 + .5)) + self.parm.FacFrom

	self.offsetPercent += self.offsetPercentStep
	for self.offsetPercent >= 1.0 {
		self.offsetPercent -= 1.0
	}
	self.offset = ((self.parm.OffTo - self.parm.OffFrom) * (math.Sin(self.offsetPercent*2*math.Pi)*.5 + .5)) + self.parm.OffFrom

}

func (self *ModeSpectrumPosition) postSetParameter() {
	self.factorPercentStep = 1.0 / float64(self.parm.FacRoundTimeMs) * float64(RefreshIntervalNs) / 1000.0 /*us*/ / 1000.0 /*ms*/
	self.offsetPercentStep = 1.0 / float64(self.parm.OffRoundTimeMs) * float64(RefreshIntervalNs) / 1000.0 /*us*/ / 1000.0 /*ms*/
	self.factorPercent = rand.Float64()
	self.offsetPercent = rand.Float64()
}

type ModeSpectrum struct {
	ModeSuper
	parameter ModeSpectrumParameter
	limits    ModeSpectrumLimits
	ledsHSV   []color.HSV
	positions [2]ModeSpectrumPosition
}

type ModeSpectrumParameterPosition struct {
	FacFrom        float64 `json:"facFrom"`
	FacTo          float64 `json:"facTo"`
	FacRoundTimeMs uint32  `json:"facRoundTimeMs"`
	OffFrom        float64 `json:"offFrom"`
	OffTo          float64 `json:"offTo"`
	OffRoundTimeMs uint32  `json:"offRoundTimeMs"`
}

type ModeSpectrumParameter struct {
	HueFrom720 float64                          `json:"hueFrom720"`
	HueTo720   float64                          `json:"hueTo720"`
	Brightness float64                          `json:"brightness"`
	Positions  [2]ModeSpectrumParameterPosition `json:"positions"`
}
type ModeSpectrumLimits struct {
	MaxRoundTimeMs uint32  `json:"maxRoundTimeMs"`
	MinRoundTimeMs uint32  `json:"minRoundTimeMs"`
	MinBrightness  float64 `json:"minBrightness"`
	MaxBrightness  float64 `json:"maxBrightness"`
	MinFactor      float64 `json:"minFactor"`
	MaxFactor      float64 `json:"maxFactor"`
	MinOffset      float64 `json:"minOffset"`
	MaxOffset      float64 `json:"maxOffset"`
}

func NewModeSpectrum(dbdriver *dbdriver.DbDriver, display *display.Display) *ModeSpectrum {
	self := ModeSpectrum{
		limits: ModeSpectrumLimits{
			MinBrightness:  0.01,
			MaxBrightness:  1.0,
			MinRoundTimeMs: 10000.0,
			MaxRoundTimeMs: 60000.0,
			MinFactor:      0.0,
			MaxFactor:      10.0,
			MinOffset:      0.0,
			MaxOffset:      math.Pi,
		},
	}

	self.ModeSuper = *NewModeSuper(dbdriver, display, "ModeSpectrum", RenderTypeDynamic, self.calcDisplay)

	self.ledsHSV = make([]color.HSV, self.display.GetRowLedCount())
	for i := 0; i < len(self.ledsHSV); i++ {
		self.ledsHSV[i] = color.HSV{H: 0.0, S: 1.0, V: 0.0}
	}

	err := dbdriver.Read(self.GetName(), "parameter", &self.parameter)
	if err != nil {
		self.Randomize()
	} else {
		self.postSetParameter()
	}

	return &self
}

func (self *ModeSpectrum) GetParameter() interface{} { return &self.parameter }
func (self *ModeSpectrum) GetLimits() interface{}    { return &self.limits }

func (self *ModeSpectrum) TrySetParameter(b []byte) error {
	var tempPar ModeSpectrumParameter
	err := json.Unmarshal(b, &tempPar)

	if err != nil {
		return err
	}

	self.SetParameter(tempPar)
	return nil
}

func (self *ModeSpectrum) correctParameter() {
	if self.parameter.HueFrom720 > self.parameter.HueTo720 {
		hueTo720 := self.parameter.HueFrom720
		self.parameter.HueFrom720 = self.parameter.HueTo720
		self.parameter.HueTo720 = hueTo720
	}
	for i := range self.parameter.Positions {
		if self.parameter.Positions[i].FacFrom > self.parameter.Positions[i].FacTo {
			facTo := self.parameter.Positions[i].FacFrom
			self.parameter.Positions[i].FacFrom = self.parameter.Positions[i].FacTo
			self.parameter.Positions[i].FacTo = facTo
		}
		if self.parameter.Positions[i].OffFrom > self.parameter.Positions[i].OffTo {
			offTo := self.parameter.Positions[i].OffFrom
			self.parameter.Positions[i].OffFrom = self.parameter.Positions[i].OffTo
			self.parameter.Positions[i].OffTo = offTo
		}
	}
}

func (self *ModeSpectrum) SetParameter(parm ModeSpectrumParameter) {
	self.parameter = parm
	self.correctParameter()
	self.dbdriver.Write(self.GetName(), "parameter", self.parameter)
	self.postSetParameter()
}

func (self *ModeSpectrum) postSetParameter() {
	for i := range self.positions {
		self.positions[i].parm = &self.parameter.Positions[i]
		self.positions[i].postSetParameter()
	}

	for i := 0; i < len(self.ledsHSV); i++ {
		self.ledsHSV[i].V = self.parameter.Brightness
	}
}

func (self *ModeSpectrum) calcDisplay() {
	for i := range self.positions {
		self.positions[i].StepForward()
	}

	hueDist := self.parameter.HueTo720 - self.parameter.HueFrom720

	for i := 0; i < len(self.ledsHSV); i++ {
		x := (float64(i) / float64(len(self.ledsHSV))) * 2 * math.Pi
		self.ledsHSV[i].H = self.parameter.HueFrom720 + ((math.Sin(self.positions[0].factor*x+self.positions[0].offset)*math.Cos(self.positions[1].factor*x+self.positions[1].offset)*0.5)+.5)*hueDist
	}

	self.display.ApplySingleRowHSV(self.ledsHSV)
}

func (self *ModeSpectrum) getRandomPosition(seed int64) ModeSpectrumParameterPosition {
	r := rand.New(rand.NewSource(seed + time.Now().UnixNano()))
	return ModeSpectrumParameterPosition{
		FacFrom:        r.Float64()*(self.limits.MaxFactor-self.limits.MinFactor) + self.limits.MinFactor,
		FacTo:          r.Float64()*(self.limits.MaxFactor-self.limits.MinFactor) + self.limits.MinFactor,
		FacRoundTimeMs: uint32(r.Float64()*(float64(self.limits.MaxRoundTimeMs)-float64(self.limits.MinRoundTimeMs)) + float64(self.limits.MinRoundTimeMs)),
		OffFrom:        r.Float64()*(self.limits.MaxOffset-self.limits.MinOffset) + self.limits.MinOffset,
		OffTo:          r.Float64()*(self.limits.MaxOffset-self.limits.MinOffset) + self.limits.MinOffset,
		OffRoundTimeMs: uint32(r.Float64()*(float64(self.limits.MaxRoundTimeMs)-float64(self.limits.MinRoundTimeMs)) + float64(self.limits.MinRoundTimeMs)),
	}
}

func (self *ModeSpectrum) Randomize() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	self.SetParameter(ModeSpectrumParameter{
		Brightness: r.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
		HueFrom720: r.Float64() * 720.0,
		HueTo720:   r.Float64() * 720.0,
		Positions:  [2]ModeSpectrumParameterPosition{self.getRandomPosition(0), self.getRandomPosition(1)},
	})
}
