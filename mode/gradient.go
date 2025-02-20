package mode

import (
	"ledean/color"
	"ledean/dbdriver"
	"ledean/display"
	"ledean/json"
	"math/rand"
	"time"
)

type ModeGradientPosition struct {
	hueFrom720    float64
	hueTo720      float64
	hueCurrent720 float64
	hueDistance   float64
	pPercentStep  *float64
	percent       float64
}

type ModeGradient struct {
	ModeSuper
	parameter    ModeGradientParameter
	presets      []ModeGradientParameter
	limits       ModeGradientLimits
	percentStep  float64
	ledsHSV      []color.HSV
	positions    []ModeGradientPosition
	posDistances []float64
}

type ModeGradientParameter struct {
	Brightness  float64 `json:"brightness"`
	Count       uint32  `json:"count"`
	RoundTimeMs uint32  `json:"roundTimeMs"`
}
type ModeGradientLimits struct {
	MinRoundTimeMs uint32  `json:"minRoundTimeMs"`
	MaxRoundTimeMs uint32  `json:"maxRoundTimeMs"`
	MinBrightness  float64 `json:"minBrightness"`
	MaxBrightness  float64 `json:"maxBrightness"`
	MinCount       uint32  `json:"minCount"`
	MaxCount       uint32  `json:"maxCount"`
}

func NewModeGradient(dbdriver *dbdriver.DbDriver, display *display.Display) *ModeGradient {
	self := ModeGradient{
		limits: ModeGradientLimits{
			MinRoundTimeMs: 1000,
			MaxRoundTimeMs: 10000,
			MinBrightness:  0.01,
			MaxBrightness:  1.0,
			MinCount:       2,
			MaxCount:       6,
		},
	}

	self.ModeSuper = *NewModeSuper(dbdriver, display, "ModeGradient", RenderTypeDynamic, self.calcDisplay)
	self.presets = self.getPresets()

	self.posDistances = make([]float64, self.limits.MaxCount-1)
	self.positions = make([]ModeGradientPosition, self.limits.MaxCount)
	for i := range self.positions {
		self.positions[i].pPercentStep = &self.percentStep
	}

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

func (self *ModeGradient) getPresets() []ModeGradientParameter {
	return []ModeGradientParameter{
		{
			Brightness:  0.6,
			Count:       3,
			RoundTimeMs: 600,
		},
		{
			Brightness:  0.7,
			Count:       4,
			RoundTimeMs: 1000,
		},
	}
}
func (self *ModeGradient) GetParameter() interface{} { return &self.parameter }
func (self *ModeGradient) GetLimits() interface{}    { return &self.limits }

func (self *ModeGradient) TrySetParameter(b []byte) error {
	var tempPar ModeGradientParameter
	err := json.Unmarshal(b, &tempPar)

	if err != nil {
		return err
	}

	self.SetParameter(tempPar)
	return nil
}

func (self *ModeGradient) SetParameter(parm ModeGradientParameter) {
	self.parameter = parm
	self.dbdriver.Write(self.GetName(), "parameter", self.parameter)
	self.postSetParameter()
}

func (self *ModeGradientPosition) StepForward() {
	self.hueCurrent720 = self.hueFrom720 + self.hueDistance*self.percent/100

	self.percent += *self.pPercentStep
	if self.percent > 100 {
		self.percent -= 100
		self.hueFrom720 = self.hueTo720
		self.randomizeWoFrom()
	}
}
func (self *ModeGradientPosition) randomizeWoFrom() {
	self.hueTo720 = rand.Float64() * 720.0

	self.hueDistance = self.hueTo720 - self.hueFrom720
}

func (self *ModeGradientPosition) Randomize() {
	self.percent = rand.Float64() * 100.0
	self.hueFrom720 = rand.Float64() * 720.0
	self.randomizeWoFrom()
}

func (self *ModeGradient) postSetParameter() {
	self.percentStep = 100 / (float64(self.parameter.RoundTimeMs) / 1000) * (float64(self.display.GetRefreshIntervalNs()) / 1000 / 1000 / 1000)

	for i := range self.positions {
		self.positions[i].Randomize()
	}

	for i := 0; i < len(self.ledsHSV); i++ {
		self.ledsHSV[i].V = self.parameter.Brightness
	}
}

func (self *ModeGradient) calcDisplayWoStep() {
	for i := 0; i < len(self.ledsHSV); i++ {
		absPos := float64(i) / float64(len(self.ledsHSV)) * float64(self.parameter.Count-1) //eg count 2 -> 1.0, count=3 ->2
		huePos := int(absPos)

		relPos := absPos - float64(huePos)

		self.ledsHSV[i].H = self.positions[huePos].hueCurrent720 + self.posDistances[huePos]*relPos
		// if self.ledsHSV[i].H > 360 {
		// 	self.ledsHSV[i].H -= 360
		// } else if self.ledsHSV[i].H < 0 {
		// 	self.ledsHSV[i].H += 360
		// }
	}

	self.display.ApplySingleRowHSV(self.ledsHSV)
}

func (self *ModeGradient) calcDisplay() {
	for i := 0; i < int(self.parameter.Count); i++ {
		self.positions[i].StepForward()
	}
	for i := 0; i < int(self.parameter.Count-1); i++ {
		self.posDistances[i] = self.positions[1+i].hueCurrent720 - self.positions[i].hueCurrent720
	}
	self.calcDisplayWoStep()
}

func (self *ModeGradient) RandomizePreset() {
	self.SetParameter(self.presets[rand.Uint32()%uint32(len(self.presets))])
}

func (self *ModeGradient) Randomize() {
	rand.Seed(time.Now().UnixNano())
	self.SetParameter(ModeGradientParameter{
		Brightness:  rand.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
		Count:       rand.Uint32()%(self.limits.MaxCount-self.limits.MinCount) + self.limits.MinCount,
		RoundTimeMs: uint32(rand.Float32()*float32(self.limits.MaxRoundTimeMs-self.limits.MinRoundTimeMs)) + self.limits.MinRoundTimeMs,
	})
}
