package mode

import (
	"LEDean/led/color"
	"math/rand"
	"time"
)

type ModeRainBowSolid struct {
	leds          []color.RGB
	cUpdate       chan bool
	minSpeedMs    uint16 //time for one rainbow round
	maxSpeedMs    uint16
	speedMs       uint16
	minBrightness uint8
	maxBrightness uint8
	brightness    uint8
	// solidColor
	refreshRate time.Duration
	shouldExit  bool
}

func NewModeRainBowSolid(leds []color.RGB, cUpdate chan bool) *ModeRainBowSolid {
	self := ModeRainBowSolid{
		leds:          leds,
		cUpdate:       cUpdate,
		minSpeedMs:    1000, //10s
		maxSpeedMs:    6000, //1min
		minBrightness: 50,
		maxBrightness: 255,                                  //math.MaxUint8,
		refreshRate:   time.Duration(33 * time.Millisecond), //30fps
		shouldExit:    false,
	}

	self.Randomize()

	return &self
}

func (self *ModeRainBowSolid) Activate() {
	self.shouldExit = false
	go func() {
		for !self.shouldExit {

			self.cUpdate <- true
			time.Sleep(self.refreshRate)
		}
	}()
}
func (self *ModeRainBowSolid) Deactivate() {
	self.shouldExit = true

}
func (self *ModeRainBowSolid) Randomize() {
	// rand.Seed(time.Now().UnixNano())
	self.speedMs = uint16(rand.Float32()*float32(self.maxSpeedMs-self.minSpeedMs)) + self.minSpeedMs
	self.brightness = uint8(rand.Float32()*float32(self.maxBrightness-self.minBrightness)) + self.minBrightness
}
