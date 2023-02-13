package mode

import (
	"encoding/json"
	"ledean/color"
	"ledean/dbdriver"
	"ledean/display"
	"math/rand"
	"time"
)

type ModeGradient struct {
	ModeSuper
	parameter           ModeGradientParameter
	limits              ModeGradientLimits
	ledsHSV             []color.HSV
	hues                [3]float64
	clockwise           [3]bool
	progressDegStepSize float64
	progressDegStep     float64
}

type ModeGradientParameter struct {
	Brightness  float64 `json:"brightness"`
	RoundTimeMs uint32  `json:"roundTimeMs"`
	Reverse     bool    `json:"reverse"`
}
type ModeGradientLimits struct {
	MinRoundTimeMs uint32  `json:"minRoundTimeMs"`
	MaxRoundTimeMs uint32  `json:"maxRoundTimeMs"`
	MinBrightness  float64 `json:"minBrightness"`
	MaxBrightness  float64 `json:"maxBrightness"`
}

func NewModeGradient(dbdriver *dbdriver.DbDriver, display *display.Display) *ModeGradient {
	self := ModeGradient{
		limits: ModeGradientLimits{
			MinRoundTimeMs: 500,   //500ms
			MaxRoundTimeMs: 30000, //30sek
			MinBrightness:  0.1,
			MaxBrightness:  1.0,
		},
	}

	self.ModeSuper = *NewModeSuper(dbdriver, display, "ModeGradient", RenderTypeDynamic, self.calcDisplay)

	self.ledsHSV = make([]color.HSV, self.display.GetRowLedCount())
	for i := 0; i < len(self.ledsHSV); i++ {
		self.ledsHSV[i] = color.HSV{H: 0.0, S: 1.0, V: 0.0}
	}

	for i := 0; i < 3; i++ {
		self.shiftGradiant()
	}

	self.progressDegStep = rand.Float64() * 360.0

	err := dbdriver.Read(self.GetName(), "parameter", &self.parameter)
	if err != nil {
		self.Randomize()
	} else {
		self.postSetParameter()
	}

	return &self
}

func (self *ModeGradient) GetParameter() interface{} { return &self.parameter }
func (self *ModeGradient) GetLimits() interface{}    { return &self.limits }

func (self *ModeGradient) TrySetParameter(b []byte) error {
	var tempPar ModeGradientParameter
	err := json.Unmarshal(b, &tempPar)

	if err != nil {
		return err
	}

	self.setParameter(tempPar)
	return nil
}

func (self *ModeGradient) setParameter(parm ModeGradientParameter) {
	self.parameter = parm
	self.dbdriver.Write(self.GetName(), "parameter", self.parameter)
	self.postSetParameter()
}

func (self *ModeGradient) postSetParameter() {
	self.progressDegStepSize = 360 / (float64(self.parameter.RoundTimeMs) / 1000) * (float64(RefreshIntervalNs) / 1000 / 1000 / 1000)
	for i := 0; i < len(self.ledsHSV); i++ {
		self.ledsHSV[i].V = self.parameter.Brightness
	}
}

func (self *ModeGradient) shiftGradiant() {
	for i := 2; i > 0; i-- {
		self.hues[i] = self.hues[i-1]
		self.clockwise[i] = self.clockwise[i-1]
	}
	self.hues[0] = rand.Float64() * 360.0
	self.clockwise[0] = rand.Int()%2 == 1
}

func (self *ModeGradient) calcDisplay() {
	self.progressDegStep += self.progressDegStepSize
	if self.progressDegStep >= 360.0 {
		self.shiftGradiant()
		self.progressDegStep -= 360.0
	}

	for i := 0; i < len(self.ledsHSV); i++ {
		SCALE := 1.0
		pos := self.progressDegStep/360.0 - SCALE*float64(i)/float64(len(self.ledsHSV)-1)
		hueI := int(pos)
		if hueI < 0 {
			hueI *= -1
		}
		posRel := pos + float64(hueI)
		hue := (self.hues[hueI]*(posRel) + self.hues[hueI+1]*(1.0-posRel))
		if hue > 360.0 {
			hue -= 360.0
		} else if hue < 0 {
			hue += 360.0
		}
		self.ledsHSV[i].H = hue
	}

	self.display.ApplySingleRowHSV(self.ledsHSV)
}

func (self *ModeGradient) Randomize() {
	rand.Seed(time.Now().UnixNano())
	self.setParameter(ModeGradientParameter{
		RoundTimeMs: uint32(rand.Float32()*float32(self.limits.MaxRoundTimeMs-self.limits.MinRoundTimeMs)) + self.limits.MinRoundTimeMs,
		Brightness:  rand.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
		Reverse:     rand.Int()%2 == 1,
	})
}
