package mode

import (
	"LEDean/led/color"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

type ModeSolid struct {
	leds          []color.RGB
	cUpdate       *chan bool
	minBrightness float64
	maxBrightness float64
	brightness    float64
	shouldExit    bool
	rgb           color.RGB
}

func NewModeSolid(leds []color.RGB, cUpdate *chan bool) *ModeSolid {
	self := ModeSolid{
		leds:          leds,
		cUpdate:       cUpdate,
		minBrightness: 0.3,
		maxBrightness: 1.0,
		shouldExit:    false,
	}

	self.Randomize()

	return &self
}

func (self *ModeSolid) GetFriendlyName() string {
	return "ModeSolid"
}

func (self *ModeSolid) Activate() {
	log.Info("wat")

	log.Debugf("start ModeSolid with:\n -self.rgb: %s\n", self.rgb)
	self.shouldExit = false

	for i := 0; i < len(self.leds); i++ {
		self.leds[i] = self.rgb
	}
	*self.cUpdate <- true

}
func (self *ModeSolid) Deactivate() {
	self.shouldExit = true
}

func (self *ModeSolid) Randomize() {
	rand.Seed(time.Now().UnixNano())
	self.brightness = rand.Float64()*(self.maxBrightness-self.minBrightness) + self.minBrightness
	self.rgb = color.RGB{
		R: uint8(rand.Float64() * self.brightness * 255),
		G: uint8(rand.Float64() * self.brightness * 255),
		B: uint8(rand.Float64() * self.brightness * 255),
	}
}
