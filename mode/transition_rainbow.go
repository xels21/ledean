package mode

import (
	"encoding/json"
	"ledean/color"
	"ledean/display"
	"math/rand"
	"time"

	"github.com/sdomino/scribble"
	log "github.com/sirupsen/logrus"
)

type ModeTransitionRainbow struct {
	dbDriver            *scribble.Driver
	display             *display.Display
	parameter           ModeTransitionRainbowParameter
	limits              ModeTransitionRainbowLimits
	cUpdate             *chan bool
	cExit               chan bool
	ledsHSV             []color.HSV
	progressDeg         float64 //from 0.0 to 360.0
	progressDegStepSize float64

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

func NewModeTransitionRainbow(dbDriver *scribble.Driver, display *display.Display) *ModeTransitionRainbow {
	self := ModeTransitionRainbow{
		dbDriver: dbDriver,
		display:  display,
		limits: ModeTransitionRainbowLimits{
			MinRoundTimeMs: 500,
			MaxRoundTimeMs: 30000,
			MinBrightness:  0.1,
			MaxBrightness:  1.0,
			MinSpectrum:    0.1,
			MaxSpectrum:    2.0,
		},
		cExit:       make(chan bool, 1),
		progressDeg: 0.0,
	}

	self.ledsHSV = make([]color.HSV, self.display.GetRowLedCount())
	for i := 0; i < len(self.ledsHSV); i++ {
		self.ledsHSV[i] = color.HSV{H: 0.0, S: 1.0, V: 0.0}
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
	for i := 0; i < len(self.ledsHSV); i++ {
		self.ledsHSV[i].H = self.ledsHSV[0].H + float64(i)/float64(len(self.ledsHSV))*self.parameter.Spectrum*360.0
		self.ledsHSV[i].V = self.parameter.Brightness
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
	for i := 0; i < len(self.ledsHSV); i++ {
		if !self.parameter.Reverse {
			self.ledsHSV[i].H += self.progressDegStepSize
			if self.ledsHSV[i].H > 360.0 {
				self.ledsHSV[i].H -= 360.0
			}
		} else {
			self.ledsHSV[i].H -= self.progressDegStepSize
			if self.ledsHSV[i].H < 0.0 {
				self.ledsHSV[i].H += 360.0
			}
		}
	}

	self.display.ApplySingleRowHSV(self.ledsHSV)
	self.display.Render()
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
