package mode

import (
	"LEDean/led/color"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/sdomino/scribble"
	log "github.com/sirupsen/logrus"
)

type ModeSolid struct {
	dbDriver  *scribble.Driver
	parameter ModeSolidParameter
	limits    ModeSolidLimits
	leds      []color.RGB
	cUpdate   *chan bool
}

type ModeSolidParameter struct {
	RGB        color.RGB `json:"rgb"`
	Brightness float64   `json:"brightness"`
}
type ModeSolidLimits struct {
	MinBrightness float64 `json:"minBrightness"`
	MaxBrightness float64 `json:"maxBrightness"`
}

func NewModeSolid(dbDriver *scribble.Driver, cUpdate *chan bool, leds []color.RGB) *ModeSolid {
	self := ModeSolid{
		dbDriver: dbDriver,
		leds:     leds,
		cUpdate:  cUpdate,
		limits: ModeSolidLimits{
			MinBrightness: 0.0,
			MaxBrightness: 1.0,
		},
	}

	err := dbDriver.Read(self.GetFriendlyName(), "parameter", &self.parameter)
	if err != nil {
		self.Randomize()
	} else {
		self.postSetParameter()
	}

	return &self
}

func (ModeSolid) GetFriendlyName() string {
	return "ModeSolid"
}

func (self *ModeSolid) GetParameterJson() []byte {
	json, _ := json.Marshal(self.parameter)
	return json
}

func (self *ModeSolid) GetLimitsJson() []byte {
	json, _ := json.Marshal(self.limits)
	return json
}

func (self *ModeSolid) SetParameter(parm interface{}) {
	switch parm.(type) {
	case ModeSolidParameter:
		self.parameter = parm.(ModeSolidParameter)
		self.dbDriver.Write(self.GetFriendlyName(), "parameter", self.parameter)
		self.postSetParameter()
	}
}

func (self *ModeSolid) postSetParameter() {
}

func (self *ModeSolid) Activate() {
	log.Debugf("start "+self.GetFriendlyName()+" with:\n %s\n", self.GetParameterJson())

	rgb := color.RGB{
		R: uint8(float64(self.parameter.RGB.R) * self.parameter.Brightness),
		G: uint8(float64(self.parameter.RGB.G) * self.parameter.Brightness),
		B: uint8(float64(self.parameter.RGB.B) * self.parameter.Brightness),
	}

	for i := 0; i < len(self.leds); i++ {
		self.leds[i] = rgb
	}
	*self.cUpdate <- true

}
func (self *ModeSolid) Deactivate() {
}

func (self *ModeSolid) Randomize() {
	rand.Seed(time.Now().UnixNano())
	parameter := ModeSolidParameter{
		Brightness: rand.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
		RGB: color.RGB{
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
		},
	}
	self.SetParameter(parameter)
}
