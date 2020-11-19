package mode

import (
	"LEDean/led/color"
	"encoding/json"
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
	rgb           color.RGB
}

type ModeSolidParameter struct {
	RGB color.RGB `json:"rgb"`
}

func NewModeSolid(leds []color.RGB, cUpdate *chan bool) *ModeSolid {
	self := ModeSolid{
		leds:          leds,
		cUpdate:       cUpdate,
		minBrightness: 0.3,
		maxBrightness: 1.0,
	}

	self.Randomize()

	return &self
}

func (ModeSolid) GetFriendlyName() string {
	return "ModeSolid"
}

func (self *ModeSolid) GetParameterJson() []byte {
	json, _ := json.Marshal(ModeSolidParameter{RGB: self.rgb})
	return json
}

func (self *ModeSolid) SetParameter(parm interface{}) {
	switch parm.(type) {
	case ModeSolidParameter:
		solidParm := parm.(ModeSolidParameter)
		self.rgb = solidParm.RGB
		self.Activate()
	}
}

func (self *ModeSolid) Activate() {
	log.Debugf("start ModeSolid with:\n -self.rgb: %s\n", self.rgb)

	for i := 0; i < len(self.leds); i++ {
		self.leds[i] = self.rgb
	}
	*self.cUpdate <- true

}
func (self *ModeSolid) Deactivate() {
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
