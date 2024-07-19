package poi

import (
	"image"
	"ledean/color"
	"ledean/dbdriver"
	"ledean/display"
	"ledean/json"
	"ledean/mode"
	poi "ledean/poi/gen"
	"log"
	"time"

	"math/rand"
)

// type PicRGB struct {
// col []PicCol
// }
//
// type PicCol struct {
// row []color.RGB
// }
//
// type PictureController struct {
// pictures []PicRGB
//
// display *display.Display
// leds    []color.RGB
// }

// display.ApplySingleRowRGB

// activatedLeds: make([]float64, display.GetRowLedCount()),

type ModePoi struct {
	mode.ModeSuper
	parameter          ModePoiParameter
	limits             ModePoiLimits
	poiPics            []*image.NRGBA
	ledsRGB            []color.RGB
	pixelCount         int
	picIndex           uint8
	picProgress        float32
	picProgressPerStep float32
	colIndex           uint32
	colProgress        float32
	colProgressPerStep float32
	// ColumnPerStep int
}

type ModePoiParameter struct {
	PictureColumnMs          uint32 `json:"pictureColumnMs"`
	PictureChangeIntervallMs uint32 `json:"pictureChangeIntervallMs"`
	// PicturePath              string `json:"picturePath"`
}

type ModePoiLimits struct {
	MinPictureColumnMs          uint32 `json:"minPictureColumnMs"`
	MaxPictureColumnMs          uint32 `json:"maxPictureColumnMs"`
	MinPictureChangeIntervallMs uint32 `json:"minPictureChangeIntervallMs"`
	MaxPictureChangeIntervallMs uint32 `json:"maxPictureChangeIntervallMs"`
}

func NewModePoi(dbdriver *dbdriver.DbDriver, display *display.Display) *ModePoi {
	if display.GetRowLedCount() != poi.PixelCount {
		log.Fatalf("Display led size[%d] not matching to generated picture size[%d]", display.GetRowLedCount(), poi.PixelCount)
	}

	self := ModePoi{
		// name: "ModePoi",
		limits: ModePoiLimits{
			MinPictureColumnMs:          5,
			MaxPictureColumnMs:          100,
			MinPictureChangeIntervallMs: 1000,
			MaxPictureChangeIntervallMs: 60000,
		}, //here must
		pixelCount:         poi.PixelCount,
		ledsRGB:            make([]color.RGB, poi.PixelCount),
		poiPics:            poi.PoiPics,
		colIndex:           0,
		colProgress:        0.0,
		colProgressPerStep: 0.0,
		picProgress:        0.0,
		picProgressPerStep: 0.0,
		picIndex:           uint8(rand.Uint32() % uint32(len(poi.PoiPics))),
	}
	self.ModeSuper = *mode.NewModeSuper(dbdriver, display, "ModePoi", mode.RenderTypeDynamic, self.calcDisplay)

	err := dbdriver.Read(self.GetName(), "parameter", &self.parameter)
	if err != nil {
		self.Randomize()
	} else {
		self.postSetParameter()
	}

	return &self
}

// func (self *ModePoi) stepForward() {
// self.ProgressPer += self.ProgressPerStep
//
//	if self.ProgressPer > 1.0 {
//		self.randomize()
//	}
//
// }

func (self *ModePoi) GetParameter() interface{} { return &self.parameter }
func (self *ModePoi) GetLimits() interface{}    { return &self.limits }

func getPixel(pic *image.NRGBA, col int, row int) color.RGB {
	rgba := pic.NRGBAAt(row, col)
	return color.RGB{R: rgba.R, G: rgba.G, B: rgba.B}
}

func (self *ModePoi) calcDisplay() {
	self.picProgress += self.picProgressPerStep
	if self.picProgress > 1.0 {
		self.picProgress -= 1.0
		self.picIndex = (self.picIndex + 1) % (uint8(len(self.poiPics) - 1))
		self.colIndex = 0
	}
	self.colProgress += self.colProgressPerStep
	if self.colProgress > 1.0 {
		self.colProgress -= 1.0
		self.colIndex = (self.colIndex + 1) % (uint32(self.poiPics[self.picIndex].Rect.Dy()) - 1)
	}

	for i := range self.pixelCount {
		self.ledsRGB[i] = getPixel(self.poiPics[self.picIndex], int(self.colIndex), i)
	}

	self.GetDisplay().ApplySingleRowRGB(self.ledsRGB)
}

func (self *ModePoi) TrySetParameter(b []byte) error {
	var tempPar ModePoiParameter
	err := json.Unmarshal(b, &tempPar)

	if err != nil {
		return err
	}

	self.SetParameter(tempPar)
	return nil
}

func (self *ModePoi) postSetParameter() {
	self.colProgressPerStep = 1.0 / (float32(self.parameter.PictureColumnMs) / 1000) * (float32(self.GetDisplay().GetRefreshIntervalNs()) / 1000 / 1000 / 1000)
	self.picProgressPerStep = 1.0 / (float32(self.parameter.PictureChangeIntervallMs) / 1000) * (float32(self.GetDisplay().GetRefreshIntervalNs()) / 1000 / 1000 / 1000)
}

func (self *ModePoi) SetParameter(parm ModePoiParameter) {
	self.parameter = parm
	self.GetDbDriver().Write(self.GetName(), "parameter", self.parameter)
	self.postSetParameter()
}

func (self *ModePoi) Randomize() {
	rand.Seed(time.Now().UnixNano())
	parameter := ModePoiParameter{
		PictureColumnMs:          (rand.Uint32())%(self.limits.MaxPictureColumnMs-self.limits.MinPictureColumnMs) + self.limits.MinPictureColumnMs,
		PictureChangeIntervallMs: (rand.Uint32())%(self.limits.MaxPictureChangeIntervallMs-self.limits.MinPictureChangeIntervallMs) + self.limits.MinPictureChangeIntervallMs,
	}
	self.SetParameter(parameter)
}
