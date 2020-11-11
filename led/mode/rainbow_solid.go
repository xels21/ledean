package mode

import (
	"LEDean/led/color"
	"math/rand"
	"time"
)

type ModeRainBowSolid struct {
	leds           []color.RGB
	cUpdate        *chan bool
	minRoundTimeMs uint16 //time for one rainbow round
	maxRoundTimeMs uint16
	roundTimeMs    uint16
	minBrightness  float64
	maxBrightness  float64
	brightness     float64
	stepSizeHue    float64
	refreshRate    time.Duration
	shouldExit     bool
	hsv            color.HSV
}

func NewModeRainBowSolid(leds []color.RGB, cUpdate *chan bool) *ModeRainBowSolid {
	self := ModeRainBowSolid{
		leds:    leds,
		cUpdate: cUpdate,
		// minRoundTimeMs: 2000, //2s
		// maxRoundTimeMs: 60000, //1min
		minRoundTimeMs: 1000, //10s
		maxRoundTimeMs: 1000, //1min
		minBrightness:  0.3,
		maxBrightness:  1.0,
		refreshRate:    time.Duration((1000 / 30) * time.Millisecond), //30fps
		shouldExit:     false,
	}

	self.Randomize()

	return &self
}

func (self *ModeRainBowSolid) Activate() {
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
			time.Sleep(self.refreshRate)
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
	self.stepSizeHue = 360.0 / (30.0 * float64(self.roundTimeMs))
}
