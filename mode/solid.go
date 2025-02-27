package mode

import (
	"encoding/json"
	"ledean/color"
	"ledean/dbdriver"
	"ledean/display"
)

type ModeSolid struct {
	ModeSuper
	parameter ModeSolidParameter
	limits    ModeSolidLimits
}

type ModeSolidParameter struct {
	RGB        color.RGB `json:"rgb"`
	Brightness float64   `json:"brightness"`
}

type ModeSolidLimits struct {
	MinBrightness float64 `json:"minBrightness"`
	MaxBrightness float64 `json:"maxBrightness"`
}

func NewModeSolid(dbdriver *dbdriver.DbDriver, display *display.Display, isRandDeterministic bool) *ModeSolid {
	self := ModeSolid{
		limits: ModeSolidLimits{
			MinBrightness: 0.0,
			MaxBrightness: 1.0,
		},
	}
	self.ModeSuper = *NewModeSuper(dbdriver, display, "ModeSolid", RenderTypeStatic, self.calcDisplay, self.calcDisplayDelta, isRandDeterministic)

	err := dbdriver.Read(self.name, "parameter", &self.parameter)
	if err != nil {
		self.Randomize()
	}

	return &self
}

func (self *ModeSolid) calcDisplayDelta(deltaTimeNs int64) {
	self.calcDisplay()
}

func (self *ModeSolid) calcDisplay() {
	rgb := color.RGB{
		R: uint8(float64(self.parameter.RGB.R) * self.parameter.Brightness),
		G: uint8(float64(self.parameter.RGB.G) * self.parameter.Brightness),
		B: uint8(float64(self.parameter.RGB.B) * self.parameter.Brightness),
	}

	self.display.AllSolid(rgb)
}

func (self *ModeSolid) GetParameter() interface{} { return &self.parameter }
func (self *ModeSolid) GetLimits() interface{}    { return &self.limits }

func (self *ModeSolid) TrySetParameter(b []byte) error {
	var tempPar ModeSolidParameter
	err := json.Unmarshal(b, &tempPar)

	if err != nil {
		return err
	}

	self.SetParameter(tempPar)
	return nil
}

func (self *ModeSolid) SetParameter(parm ModeSolidParameter) {
	self.parameter = parm
	self.dbdriver.Write(self.name, "parameter", self.parameter)
}

func (self *ModeSolid) RandomizePreset() {
	self.Randomize()
}

func (self *ModeSolid) Randomize() {
	parameter := ModeSolidParameter{
		Brightness: self.rand.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
		RGB: color.RGB{
			R: uint8(self.rand.Intn(255)),
			G: uint8(self.rand.Intn(255)),
			B: uint8(self.rand.Intn(255)),
		},
	}
	self.SetParameter(parameter)
}
