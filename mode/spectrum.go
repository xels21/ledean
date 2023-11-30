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
	// hueFrom720    float64
	// hueTo720      float64
	// hueCurrent720 float64
	// hueDistance   float64
	// pPercentStep  *float64
	// percent       float64
}

type ModeSpectrum struct {
	ModeSuper
	parameter ModeSpectrumParameter
	limits    ModeSpectrumLimits
	ledsHSV   []color.HSV
	// percentStep  float64
	// positions    []ModeSpectrumPosition
	// posDistances []float64
}

type ModeSpectrumParameter struct {
	HueFrom720 float64 `json:"hueFrom720"`
	HueTo720   float64 `json:"hueTo720"`
	Brightness float64 `json:"brightness"`
	// Count       uint32  `json:"count"`
	// RoundTimeMs uint32  `json:"roundTimeMs"`
}
type ModeSpectrumLimits struct {
	// MinRoundTimeMs uint32  `json:"minRoundTimeMs"`
	// MaxRoundTimeMs uint32  `json:"maxRoundTimeMs"`
	MinBrightness float64 `json:"minBrightness"`
	MaxBrightness float64 `json:"maxBrightness"`
	// MinCount       uint32  `json:"minCount"`
	// MaxCount       uint32  `json:"maxCount"`
}

func NewModeSpectrum(dbdriver *dbdriver.DbDriver, display *display.Display) *ModeSpectrum {
	self := ModeSpectrum{
		limits: ModeSpectrumLimits{
			// MinRoundTimeMs: 1000,
			// MaxRoundTimeMs: 10000,
			MinBrightness: 0.01,
			MaxBrightness: 1.0,
			// MinCount:       2,
			// MaxCount:       6,
		},
	}

	self.ModeSuper = *NewModeSuper(dbdriver, display, "ModeSpectrum", RenderTypeDynamic, self.calcDisplay)

	// self.posDistances = make([]float64, self.limits.MaxCount-1)
	// self.positions = make([]ModeSpectrumPosition, self.limits.MaxCount)
	// for i := range self.positions {
	// self.positions[i].pPercentStep = &self.percentStep
	// }

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

func (self *ModeSpectrum) SetParameter(parm ModeSpectrumParameter) {
	if parm.HueFrom720 > parm.HueTo720 {
		hueTo720 := parm.HueFrom720
		parm.HueFrom720 = parm.HueTo720
		parm.HueTo720 = hueTo720
	}
	self.parameter = parm
	self.dbdriver.Write(self.GetName(), "parameter", self.parameter)
	self.postSetParameter()
}

// func (self *ModeSpectrumPosition) StepForward() {
// 	self.hueCurrent720 = self.hueFrom720 + self.hueDistance*self.percent/100

// 	self.percent += *self.pPercentStep
// 	if self.percent > 100 {
// 		self.percent -= 100
// 		self.hueFrom720 = self.hueTo720
// 		self.randomizeWoFrom()
// 	}
// }
// func (self *ModeSpectrumPosition) randomizeWoFrom() {
// 	self.hueTo720 = rand.Float64() * 720.0

// 	self.hueDistance = self.hueTo720 - self.hueFrom720
// }

// func (self *ModeSpectrumPosition) Randomize() {
// self.percent = rand.Float64() * 100.0
// self.hueFrom720 = rand.Float64() * 720.0
// self.randomizeWoFrom()
// }

func (self *ModeSpectrum) postSetParameter() {
	// self.percentStep = 100 / (float64(self.parameter.RoundTimeMs) / 1000) * (float64(RefreshIntervalNs) / 1000 / 1000 / 1000)
	//
	// for i := range self.positions {
	// self.positions[i].Randomize()
	// }

	for i := 0; i < len(self.ledsHSV); i++ {
		self.ledsHSV[i].V = self.parameter.Brightness
	}
}

// func (self *ModeSpectrum) calcDisplayWoStep() {
// for i := 0; i < len(self.ledsHSV); i++ {
// absPos := float64(i) / float64(len(self.ledsHSV)) * float64(self.parameter.Count-1) //eg count 2 -> 1.0, count=3 ->2
// huePos := int(absPos)
//
// relPos := absPos - float64(huePos)
//
// self.ledsHSV[i].H = self.positions[huePos].hueCurrent720 + self.posDistances[huePos]*relPos
// }

// self.display.ApplySingleRowHSV(self.ledsHSV)
// }

func (self *ModeSpectrum) calcDisplay() {
	// for i := 0; i < int(self.parameter.Count); i++ {
	// self.positions[i].StepForward()
	// }
	// for i := 0; i < int(self.parameter.Count-1); i++ {
	// self.posDistances[i] = self.positions[1+i].hueCurrent720 - self.positions[i].hueCurrent720
	// }
	// self.calcDisplayWoStep()

	hueDist := self.parameter.HueTo720 - self.parameter.HueFrom720

	for i := 0; i < len(self.ledsHSV); i++ {
		x := (float64(i) / float64(len(self.ledsHSV))) * 2 * math.Pi
		self.ledsHSV[i].H = self.parameter.HueFrom720 + ((math.Sin(2*x+2)*math.Sin(.5*x+4)*0.5)+.5)*hueDist
		// absPos := x * float64(self.parameter.Count-1) //eg count 2 -> 1.0, count=3 ->2
		// huePos := int(absPos)

		// relPos := absPos - float64(huePos)

		// self.ledsHSV[i].H = self.positions[huePos].hueCurrent720 + self.posDistances[huePos]*relPos
	}

	self.display.ApplySingleRowHSV(self.ledsHSV)
}

func (self *ModeSpectrum) Randomize() {
	rand.Seed(time.Now().UnixNano())
	self.SetParameter(ModeSpectrumParameter{
		Brightness: rand.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
		HueFrom720: rand.Float64() * 720.0,
		HueTo720:   rand.Float64() * 720.0,
		// Count:       rand.Uint32()%(self.limits.MaxCount-self.limits.MinCount) + self.limits.MinCount,
		// RoundTimeMs: uint32(rand.Float32()*float32(self.limits.MaxRoundTimeMs-self.limits.MinRoundTimeMs)) + self.limits.MinRoundTimeMs,
	})
}
