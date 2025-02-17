package mode

import (
	"image"
	"ledean/color"
	"ledean/dbdriver"
	"ledean/display"
	"ledean/json"
	"ledean/log"
	picture "ledean/mode/gen_picture"
	"time"

	"math/rand"
)

type ModePicture struct {
	ModeSuper
	parameter ModePictureParameter
	limits    ModePictureLimits
	// poiPics            []image.NRGBA
	currentPic         [][]color.RGB
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

type ModePictureParameter struct {
	PictureColumnNs          uint32  `json:"pictureColumnNs"`
	PictureChangeIntervallMs uint32  `json:"pictureChangeIntervallMs"`
	Brightness               float64 `json:"brightness"`
	// PicturePath              string `json:"picturePath"`
}

type ModePictureLimits struct {
	MinPictureColumnNs          uint32  `json:"minPictureColumnNs"`
	MaxPictureColumnNs          uint32  `json:"maxPictureColumnNs"`
	MinPictureChangeIntervallMs uint32  `json:"minPictureChangeIntervallMs"`
	MaxPictureChangeIntervallMs uint32  `json:"maxPictureChangeIntervallMs"`
	MinBrightness               float64 `json:"minBrightness"`
	MaxBrightness               float64 `json:"maxBrightness"`
}

func NewModePicture(dbdriver *dbdriver.DbDriver, display *display.Display) *ModePicture {
	if display.GetRowLedCount() != picture.PixelCount {
		log.Fatalf("Display led size[%d] not matching to generated picture size[%d]", display.GetRowLedCount(), picture.PixelCount)
	}
	self := ModePicture{
		// name: "ModePicture",
		limits: ModePictureLimits{
			MinPictureColumnNs:          1,
			MaxPictureColumnNs:          1000,
			MinPictureChangeIntervallMs: 1000,
			MaxPictureChangeIntervallMs: 60000,
			MinBrightness:               0.0,
			MaxBrightness:               1.0,
		}, //here must
		pixelCount: picture.PixelCount,
		ledsRGB:    make([]color.RGB, picture.PixelCount),
		// poiPics:            make([]image.NRGBA, len(picture.PoiPics)),
		colIndex:           0,
		colProgress:        0.0,
		colProgressPerStep: 0.0,
		picProgress:        0.0,
		picProgressPerStep: 0.0,
		// picIndex:           0,
		picIndex: uint8(rand.Uint32() % uint32(len(picture.Pics))),
	}
	self.ModeSuper = *NewModeSuper(dbdriver, display, "ModePicture", RenderTypeDynamic, self.calcDisplay)

	err := dbdriver.Read(self.GetName(), "parameter", &self.parameter)
	if err != nil {
		// self.Randomize()
		self.Default()
	} else {
		self.postSetParameter()
	}

	return &self
}

func (self *ModePicture) Default() {
	parameter := ModePictureParameter{

		PictureColumnNs:          1,
		PictureChangeIntervallMs: 5000,
		Brightness:               0.1,
	}
	self.SetParameter(parameter)
}

func (self *ModePicture) GetParameter() interface{} { return &self.parameter }
func (self *ModePicture) GetLimits() interface{}    { return &self.limits }

func getPixel(pic *image.NRGBA, col int, row int) color.RGB {
	rgba := pic.NRGBAAt(row, col)
	return color.RGB{R: rgba.R, G: rgba.G, B: rgba.B}
}

// func PixelsScale(pixels            []color.RGB){
// 	for i := range pixels {
// 		pixels[i].B = pixels[i].B
// 	}

// }

func (self *ModePicture) updateCurrentPic() {
	self.colIndex = 0
	pPic := picture.Pics[self.picIndex]
	rows := len(pPic)
	self.currentPic = make([][]color.RGB, rows)
	for r := 0; r < rows; r++ {
		self.currentPic[r] = make([]color.RGB, self.pixelCount)
		for i := 0; i < self.pixelCount; i++ {
			// TODO: brightnes
			self.currentPic[r][i] = color.RGB{R: uint8(float64(pPic[r][i*3+0]) * self.parameter.Brightness),
				G: uint8(float64(pPic[r][i*3+1]) * self.parameter.Brightness),
				B: uint8(float64(pPic[r][i*3+2]) * self.parameter.Brightness)}
		}
	}
}

func (self *ModePicture) calcDisplay() {
	self.picProgress += self.picProgressPerStep
	if self.picProgress > 1.0 {
		self.picProgress -= 1.0
		self.picIndex = (self.picIndex + 1) % (uint8(len(picture.Pics)))
		self.updateCurrentPic()
	}
	self.colProgress += self.colProgressPerStep
	if self.colProgress > 1.0 {
		self.colProgress -= 1.0
		self.colIndex = (self.colIndex + 1) % (uint32(len(self.currentPic)) - 1)
	}

	self.ledsRGB = self.currentPic[self.colIndex]
	// for i := range self.pixelCount {
	// self.ledsRGB[i] = getPixel(&self.poiPics[self.picIndex], int(self.colIndex), i)
	// }

	self.GetDisplay().ApplySingleRowRGB(self.ledsRGB)
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
	// for iPic := range len(picture.PoiPics) {
	// self.poiPics[iPic].Pix = make([]uint8, len(picture.PoiPics[iPic].Pix))
	// self.poiPics[iPic].Rect = picture.PoiPics[iPic].Rect
	// self.poiPics[iPic].Stride = picture.PoiPics[iPic].Stride
	// for iPix := range picture.PoiPics[iPic].Pix {
	// self.poiPics[iPic].Pix[iPix] = uint8(float64(picture.PoiPics[iPic].Pix[iPix]) * self.parameter.Brightness)
	// }
	// }
	self.updateCurrentPic()
	self.colProgressPerStep = 1.0 / (float32(self.parameter.PictureColumnNs / 1000 / 1000)) * (float32(self.GetDisplay().GetRefreshIntervalNs()) / 1000 / 1000 / 1000)
	self.picProgressPerStep = 1.0 / (float32(self.parameter.PictureChangeIntervallMs) / 1000) * (float32(self.GetDisplay().GetRefreshIntervalNs()) / 1000 / 1000 / 1000)
}

func (self *ModePicture) SetParameter(parm ModePictureParameter) {
	self.parameter = parm
	self.GetDbDriver().Write(self.GetName(), "parameter", self.parameter)
	self.postSetParameter()
}

func (self *ModePicture) Randomize() {
	rand.Seed(time.Now().UnixNano())
	self.picIndex = uint8(rand.Uint32() % uint32(len(picture.Pics)))
	parameter := ModePictureParameter{
		PictureColumnNs:          (rand.Uint32())%(self.limits.MaxPictureColumnNs-self.limits.MinPictureColumnNs) + self.limits.MinPictureColumnNs,
		PictureChangeIntervallMs: (rand.Uint32())%(self.limits.MaxPictureChangeIntervallMs-self.limits.MinPictureChangeIntervallMs) + self.limits.MinPictureChangeIntervallMs,
		Brightness:               rand.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
	}
	self.SetParameter(parameter)
}
