package mode

import (
	"LEDean/led/color"
	"math/rand"
	"time"

	"github.com/sdomino/scribble"
	log "github.com/sirupsen/logrus"
)

type ModeRainBowSolid struct {
	dbDriver       *scribble.Driver
	leds           []color.RGB
	cUpdate        *chan bool
	minRoundTimeMs uint16 //time for one rainbow round
	maxRoundTimeMs uint16
	roundTimeMs    uint16
	minBrightness  float64
	maxBrightness  float64
	brightness     float64
	stepSizeHue    float64
	refreshRateNs  time.Duration
	shouldExit     bool
	hsv            color.HSV
}

func NewModeRainBowSolid(leds []color.RGB, cUpdate *chan bool, dbDriver *scribble.Driver) *ModeRainBowSolid {
	self := ModeRainBowSolid{
		dbDriver: dbDriver,
		leds:     leds,
		cUpdate:  cUpdate,
		// minRoundTimeMs: 2000, //2s
		// maxRoundTimeMs: 60000, //1min
		minRoundTimeMs: 5000, //10s
		maxRoundTimeMs: 5000, //1min
		minBrightness:  0.3,
		maxBrightness:  1.0,
		refreshRateNs:  time.Duration((1000 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/ / 30) * time.Nanosecond), //30fps
		shouldExit:     false,
	}

	self.Randomize()

	return &self
}

func (ModeRainBowSolid) GetFriendlyName() string {
	return "ModeRainBowSolid"
}

func (self *ModeRainBowSolid) GetParameterJson() []byte {
	// msg, _ := json.Marshal(ModeSolidParameter{rgb: self.rgb})

	return []byte{}
}

func (self *ModeRainBowSolid) SetParameter(parm interface{}) {

}

func (self *ModeRainBowSolid) Activate() {
	log.Debugf("start ModeRainBowSolid with:\n -self.roundTimeMs: %d\n -self.stepSizeHue: %f\n- self.refreshRateNs: %d", self.roundTimeMs, self.stepSizeHue, self.refreshRateNs)
	self.shouldExit = false
	go func() {
		rgb := color.RGB{}
		for !self.shouldExit {
			self.hsv.H += self.stepSizeHue
			for self.hsv.H > 360.0 {
				self.hsv.H -= 360.0
			}
			rgb = self.hsv.ToRGB()
			for i := 0; i < len(self.leds); i++ {
				self.leds[i] = rgb
			}
			*self.cUpdate <- true
			time.Sleep(self.refreshRateNs)
		}
	}()
}
func (self *ModeRainBowSolid) Deactivate() {
	self.shouldExit = true
}

func (self *ModeRainBowSolid) Randomize() {
	rand.Seed(time.Now().UnixNano())
	self.roundTimeMs = uint16(rand.Float32()*float32(self.maxRoundTimeMs-self.minRoundTimeMs)) + self.minRoundTimeMs
	self.brightness = rand.Float64()*(self.maxBrightness-self.minBrightness) + self.minBrightness
	self.hsv = color.HSV{
		H: rand.Float64() * 360.0,
		S: 1.0,
		V: self.brightness,
	}
	self.stepSizeHue = 360.0 / (float64(self.roundTimeMs) / 1000) * (float64(self.refreshRateNs) / 1000 / 1000 / 1000)
}
