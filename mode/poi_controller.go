package mode

import (
	"image"
	"ledean/color"
	"ledean/dbdriver"
	"ledean/display"
	"ledean/json"
	poi "ledean/mode/gen_poi"
	"time"

	"math/rand"
)

type ModePoi struct {
	ModeSuper
	parameter ModePoiParameter
	limits    ModePoiLimits
	// getPoiPics         []poi.GetPoiPic
	poiPic             image.NRGBA
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
	PictureColumnNs          uint32  `json:"pictureColumnNs"`
	PictureChangeIntervallMs uint32  `json:"pictureChangeIntervallMs"`
	Brightness               float64 `json:"brightness"`
	// PicturePath              string `json:"picturePath"`
}

type ModePoiLimits struct {
	MinPictureColumnNs          uint32  `json:"minPictureColumnNs"`
	MaxPictureColumnNs          uint32  `json:"maxPictureColumnNs"`
	MinPictureChangeIntervallMs uint32  `json:"minPictureChangeIntervallMs"`
	MaxPictureChangeIntervallMs uint32  `json:"maxPictureChangeIntervallMs"`
	MinBrightness               float64 `json:"minBrightness"`
	MaxBrightness               float64 `json:"maxBrightness"`
}

func NewModePoi(dbdriver *dbdriver.DbDriver, display *display.Display) *ModePoi {
	self := ModePoi{
		// name: "ModePoi",
		limits: ModePoiLimits{
			MinPictureColumnNs:          1,
			MaxPictureColumnNs:          5000,
			MinPictureChangeIntervallMs: 1000,
			MaxPictureChangeIntervallMs: 60000,
			MinBrightness:               0.0,
			MaxBrightness:               1.0,
		}, //here must
		pixelCount: poi.PixelCount,
		ledsRGB:    make([]color.RGB, poi.PixelCount),
		// poiPics:            poi.PoiPicsGetter, //make([]image.NRGBA, len(poi.PoiPics)),
		colIndex:           0,
		colProgress:        0.0,
		colProgressPerStep: 0.0,
		picProgress:        0.0,
		picProgressPerStep: 0.0,
		picIndex:           0,
		// picIndex:           uint8(rand.Uint32() % uint32(len(poi.PoiPics))),
	}
	self.ModeSuper = *NewModeSuper(dbdriver, display, "ModePoi", RenderTypeDynamic, self.calcDisplay)

	err := dbdriver.Read(self.GetName(), "parameter", &self.parameter)
	if err != nil {
		// self.Randomize()
		self.Default()
	} else {
		self.postSetParameter()
	}

	return &self
}

func (self *ModePoi) Default() {
	parameter := ModePoiParameter{
		PictureColumnNs:          1,
		PictureChangeIntervallMs: 5000,
		Brightness:               0.1,
	}
	self.SetParameter(parameter)
}

func (self *ModePoi) GetParameter() interface{} { return &self.parameter }
func (self *ModePoi) GetLimits() interface{}    { return &self.limits }

func getPixel(pic *image.NRGBA, col int, row int) color.RGB {
	rgba := pic.NRGBAAt(row, col)
	return color.RGB{R: rgba.R, G: rgba.G, B: rgba.B}
}

// func PixelsScale(pixels            []color.RGB){
// 	for i := range pixels {
// 		pixels[i].B = pixels[i].B
// 	}

// }

func (self *ModePoi) loadPic(i uint8) {
	self.poiPic = poi.PoiPicsGetter[i]()
}

func (self *ModePoi) calcDisplay() {
	self.picProgress += self.picProgressPerStep
	if self.picProgress > 1.0 {
		self.picProgress -= 1.0
		self.picIndex = (self.picIndex + 1) % (uint8(poi.PoiPicsCount))
		self.colIndex = 0
		self.loadPic(self.picIndex)
	}
	self.colProgress += self.colProgressPerStep
	if self.colProgress > 1.0 {
		self.colProgress -= 1.0
		self.colIndex = (self.colIndex + 1) % (uint32(self.poiPic.Rect.Dy()) - 1)
	}

	for i := range self.pixelCount {
		self.ledsRGB[i] = getPixel(&self.poiPic, int(self.colIndex), i)
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

	// for iPic := range len(poi.PoiPics) {
	// self.poiPics[iPic].Pix = make([]uint8, len(poi.PoiPics[iPic].Pix))
	// self.poiPics[iPic].Rect = poi.PoiPics[iPic].Rect
	// self.poiPics[iPic].Stride = poi.PoiPics[iPic].Stride
	// for iPix := range poi.PoiPics[iPic].Pix {
	// self.poiPics[iPic].Pix[iPix] = uint8(float64(poi.PoiPics[iPic].Pix[iPix]) * self.parameter.Brightness)
	// }
	// }
	self.loadPic(self.picIndex)
	self.colProgressPerStep = 1.0 / (float32(self.parameter.PictureColumnNs / 1000 / 1000)) * (float32(self.GetDisplay().GetRefreshIntervalNs()) / 1000 / 1000 / 1000)
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
		PictureColumnNs:          (rand.Uint32())%(self.limits.MaxPictureColumnNs-self.limits.MinPictureColumnNs) + self.limits.MinPictureColumnNs,
		PictureChangeIntervallMs: (rand.Uint32())%(self.limits.MaxPictureChangeIntervallMs-self.limits.MinPictureChangeIntervallMs) + self.limits.MinPictureChangeIntervallMs,
		Brightness:               rand.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
	}
	self.SetParameter(parameter)
}
