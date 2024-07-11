package mode

import (
	"ledean/dbdriver"
	"ledean/display"
	"ledean/json"
)

type ModePicture struct {
	ModeSuper
	parameter     ModePictureParameter
	limits        ModePictureLimits
	Column        int
	ColumnPerStep int
}

type ModePictureParameter struct {
	PictureColumnMs int    `json:"pictureColumnMs"`
	PicturePath     string `json:"picturePath"`
}

type ModePictureLimits struct {
	MinPictureColumnMs int      `json:"minPictureColumnMs"`
	MaxPictureColumnMs int      `json:"maxPictureColumnMs"`
	PicturePaths       []string `json:"picturePath"`
}

func NewModePicture(dbdriver *dbdriver.DbDriver, display *display.Display) *ModePicture {
	self := ModePicture{
		limits: ModePictureLimits{}, //here must
	}
	self.ModeSuper = *NewModeSuper(dbdriver, display, "ModePicture", RenderTypeDynamic, self.calcDisplay)

	err := dbdriver.Read(self.name, "parameter", &self.parameter)
	if err != nil {
		self.Randomize()
	} else {
		self.postSetParameter()
	}

	return &self
}

func (self *ModePicture) stepForward() {
	// self.ProgressPer += self.ProgressPerStep
	// if self.ProgressPer > 1.0 {
	// 	self.randomize()
	// }
}
func (self *ModePicture) randomize() {
}
func (self *ModePicture) GetParameter() interface{} { return &self.parameter }
func (self *ModePicture) GetLimits() interface{}    { return &self.limits }

func (self *ModePicture) calcDisplay() {
	// color.HsvArrClear(self.ledsHSV)
	// for i := uint8(0); i < self.parameter.PictureCount; i++ {
	// 	self.emits[i].stepForward()
	// 	if self.emits[i].ProgressPer < 0 {
	// 		continue
	// 	}

	// 	switch self.parameter.PictureStyle {
	// 	case PictureStylePulse:
	// 		self.emits[i].addPulseToLeds(self.ledsHSV)
	// 	case PictureStyleDrop:
	// 		self.emits[i].addDropToLeds(self.ledsHSV)
	// 	}

	// }
	// self.display.ApplySingleRowHSV(self.ledsHSV)
}

func (self *ModePicture) TrySetParameter(b []byte) error {
	var tempPar ModePictureParameter
	err := json.Unmarshal(b, &tempPar)

	if err != nil {
		return err
	}

	self.SetParameter(tempPar)
	return nil
}

func (self *ModePicture) postSetParameter() {
	// for i := uint8(0); i < self.parameter.PictureCount; i++ {
	// 	self.emits[i].randomize()
	// }
}

func (self *ModePicture) SetParameter(parm ModePictureParameter) {
	self.parameter = parm
	self.dbdriver.Write(self.name, "parameter", self.parameter)
	self.postSetParameter()
}

func (self *ModePicture) Randomize() {
	// rand.Seed(time.Now().UnixNano())
	// minBrightness := self.limits.MinBrightness + (rand.Float64() * (self.limits.MaxBrightness - self.limits.MinBrightness))
	// minPictureLifetimeMs := self.limits.MinPictureLifetimeMs + (rand.Uint32() % (self.limits.MaxPictureLifetimeMs - self.limits.MinPictureLifetimeMs))
	// parameter := ModePictureParameter{
	// 	PictureCount:         uint8(rand.Uint32())%(self.limits.MaxPictureCount-self.limits.MinPictureCount+1) + self.limits.MinPictureCount,
	// 	PictureStyle:         self.getRandomStyle(),
	// 	MinBrightness:        minBrightness,
	// 	MaxBrightness:        minBrightness + (rand.Float64() * (self.limits.MaxBrightness - minBrightness)),
	// 	MinPictureLifetimeMs: minPictureLifetimeMs,
	// 	MaxPictureLifetimeMs: minPictureLifetimeMs + (rand.Uint32() % (self.limits.MaxPictureLifetimeMs - minPictureLifetimeMs)),
	// 	WaveSpeedFac:         1.0, //TODO
	// 	WaveWidthFac:         1.0, //TODO
	// }
	// self.SetParameter(parameter)
}
