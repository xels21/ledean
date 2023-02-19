package mode

import (
	"encoding/json"
	"ledean/color"
	"ledean/dbdriver"
	"ledean/display"
	"math/rand"
	"time"
)

type ModeGradientPosition struct {
	hueFrom     float64
	hueTo       float64
	hueCurrent  float64
	hueDistance float64
	positive    bool
	roundTimeMs uint32
	percentStep float64
	percent     float64
	limits      *ModeGradientLimits
}

type ModeGradient struct {
	ModeSuper
	parameter ModeGradientParameter
	limits    ModeGradientLimits
	ledsHSV   []color.HSV
	positions []ModeGradientPosition
}

type ModeGradientParameter struct {
	Brightness float64 `json:"brightness"`
	Count      uint32  `json:"count"`
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
			MinRoundTimeMs: 5000,
			MaxRoundTimeMs: 5000,
			MinBrightness:  0.1,
			MaxBrightness:  1.0,
			MinCount:       2,
			MaxCount:       6,
		},
	}

	self.ModeSuper = *NewModeSuper(dbdriver, display, "ModeGradient", RenderTypeDynamic, self.calcDisplay)

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

func (self *ModeGradientPosition) StepForward() {
	self.hueCurrent = self.hueFrom + self.hueDistance*self.percent/100
	if self.hueCurrent < 0 {
		self.hueCurrent += 360
	} else if self.hueCurrent > 360 {
		self.hueCurrent -= 360
	}

	self.percent += self.percentStep
	if self.percent > 100 {
		self.percent -= 100
		self.hueFrom = self.hueTo
		self.randomizeWoFrom()
	}
}
func (self *ModeGradientPosition) randomizeWoFrom() {
	self.hueTo = rand.Float64() * 360.0
	self.positive = rand.Uint32()&1 == 1
	self.roundTimeMs = uint32(rand.Float32()*float32(self.limits.MaxRoundTimeMs-self.limits.MinRoundTimeMs)) + self.limits.MinRoundTimeMs
	self.percentStep = 100 / (float64(self.roundTimeMs) / 1000) * (float64(RefreshIntervalNs) / 1000 / 1000 / 1000)

	self.hueDistance = self.hueTo - self.hueFrom
	if self.hueDistance < 0 && self.positive {
		self.hueDistance += 360.0
	} else if self.hueDistance > 0 && !self.positive {
		self.hueDistance -= 360.0
	}
}
func (self *ModeGradientPosition) Randomize() {
	self.percent = rand.Float64() * 100.0
	self.hueFrom = rand.Float64() * 360.0
	self.randomizeWoFrom()
}

func (self *ModeGradient) postSetParameter() {
	self.positions = make([]ModeGradientPosition, self.parameter.Count)
	for i := range self.positions {
		self.positions[i].percent = 0
		self.positions[i].limits = &self.limits
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
		if huePos != 0 {
			relPos = 0
		}
		self.ledsHSV[i].H = self.positions[huePos].hueCurrent*(1.0-relPos) + self.positions[huePos+1].hueCurrent*relPos

		// self.ledsHSV[i].H = self.positions[huePos].hueCurrent

		// if !(self.ledsHSV[i].H >= self.positions[huePos].hueCurrent && self.ledsHSV[i].H <= self.positions[huePos+1].hueCurrent) && !(self.ledsHSV[i].H <= self.positions[huePos].hueCurrent && self.ledsHSV[i].H >= self.positions[huePos+1].hueCurrent) {
		// fmt.Printf("%d_%d,%d\n", self.ledsHSV[i].H, self.positions[huePos].hueCurrent, self.positions[huePos+1].hueCurrent)
		// }
	}
	// self.percent += self.percentStep
	// if self.percent >= 360.0 {
	// 	self.shiftGradiant()
	// 	self.percent -= 360.0
	// }

	// for i := 0; i < len(self.ledsHSV); i++ {
	// 	SCALE := 1.0
	// 	pos := self.percent/360.0 - SCALE*float64(i)/float64(len(self.ledsHSV)-1)
	// 	hueI := int(pos)
	// 	if hueI < 0 {
	// 		hueI *= -1
	// 	}
	// 	posRel := pos + float64(hueI)
	// 	hue := (self.hues[hueI]*(posRel) + self.hues[hueI+1]*(1.0-posRel))
	// 	if hue > 360.0 {
	// 		hue -= 360.0
	// 	} else if hue < 0 {
	// 		hue += 360.0
	// 	}
	// 	self.ledsHSV[i].H = hue
	// }

	self.display.ApplySingleRowHSV(self.ledsHSV)
}

func (self *ModeGradient) calcDisplay() {
	for i := range self.positions {
		self.positions[i].StepForward()
	}
	self.calcDisplayWoStep()
}

func (self *ModeGradient) Randomize() {
	rand.Seed(time.Now().UnixNano())
	self.setParameter(ModeGradientParameter{
		Brightness: rand.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
		Count:      rand.Uint32()%(self.limits.MaxCount-self.limits.MinCount) + self.limits.MinCount,
	})
}
