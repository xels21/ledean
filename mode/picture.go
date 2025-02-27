package mode

import (
	"image"
	"ledean/color"
	"ledean/dbdriver"
	"ledean/display"
	"ledean/json"
	"ledean/log"
	picture "ledean/mode/gen_picture"
	"sync"
)

type ModePicture struct {
	ModeSuper
	parameter ModePictureParameter
	limits    ModePictureLimits
	// poiPics            []image.NRGBA
	muPic              sync.Mutex
	currentPic         [][]color.RGB
	ledsRGB            []color.RGB
	pixelCount         int
	picIndex           uint8
	picProgress        float64
	picProgressPerStep float64
	colIndex           uint32
	colProgress        float64
	colProgressPerStep float64
	// ColumnPerStep int
}

type ModePictureParameter struct {
	PictureColumnUs          uint32  `json:"pictureColumnUs"`
	PictureChangeIntervallMs uint32  `json:"pictureChangeIntervallMs"`
	Brightness               float64 `json:"brightness"`
	// PicturePath              string `json:"picturePath"`
}

type ModePictureLimits struct {
	MinPictureColumnUs          uint32  `json:"minPictureColumnUs"`
	MaxPictureColumnUs          uint32  `json:"maxPictureColumnUs"`
	MinPictureChangeIntervallMs uint32  `json:"minPictureChangeIntervallMs"`
	MaxPictureChangeIntervallMs uint32  `json:"maxPictureChangeIntervallMs"`
	MinBrightness               float64 `json:"minBrightness"`
	MaxBrightness               float64 `json:"maxBrightness"`
}

func NewModePicture(dbdriver *dbdriver.DbDriver, display *display.Display, isRandDeterministic bool) *ModePicture {
	if display.GetRowLedCount() != picture.PixelCount {
		log.Warningf("Display led size[%d] not matching to generated picture size[%d]", display.GetRowLedCount(), picture.PixelCount)
	}
	self := ModePicture{
		// name: "ModePicture",
		limits: ModePictureLimits{
			MinPictureColumnUs:          1000,
			MaxPictureColumnUs:          20000,
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
	}
	self.ModeSuper = *NewModeSuper(dbdriver, display, "ModePicture", RenderTypeDynamic, self.calcDisplay, self.calcDisplayDelta, isRandDeterministic)

	self.picIndex = uint8(self.rand.Uint32() % uint32(len(picture.Pics)))
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
		PictureColumnUs:          10000,
		PictureChangeIntervallMs: 3000,
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
	self.muPic.Lock()
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
	self.muPic.Unlock()
}

func (self *ModePicture) calcDisplayFinal(picProgressPerStep float64, colProgressPerStep float64) {
	self.picProgress += picProgressPerStep
	if self.picProgress > 1.0 {
		self.picProgress -= 1.0
		self.picIndex = (self.picIndex + 1) % (uint8(len(picture.Pics)))
		self.updateCurrentPic()
	}
	self.colProgress += colProgressPerStep
	if self.colProgress > 1.0 {
		self.colProgress -= 1.0
		self.colIndex = (self.colIndex + 1) % (uint32(len(self.currentPic)) - 1)
	}

	self.muPic.Lock()
	self.ledsRGB = self.currentPic[self.colIndex]
	self.muPic.Unlock()

	// for i := range self.pixelCount {
	// self.ledsRGB[i] = getPixel(&self.poiPics[self.picIndex], int(self.colIndex), i)
	// }

	self.GetDisplay().ApplySingleRowRGB(self.ledsRGB)
}

func (self *ModePicture) calcDisplayDelta(deltaTimeNs int64) {
	self.calcDisplayFinal(
		self.getPicProgressPerStep(float64(deltaTimeNs)),
		self.getColProgressPerStep(float64(deltaTimeNs)),
	)

}

func (self *ModePicture) calcDisplay() {
	self.calcDisplayFinal(self.picProgressPerStep, self.colProgressPerStep)
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

func (self *ModePicture) getColProgressPerStep(timeNs float64) float64 {
	//  1.0 / (float32(self.parameter.PictureColumnUs / 1000 / 1000)) * (float32(self.GetDisplay().GetRefreshIntervalNs()) / 1000 / 1000 / 1000)
	return 1.0 / (float64(self.parameter.PictureColumnUs) / 1000 / 1000) * (timeNs / 1000 / 1000 / 1000)
}
func (self *ModePicture) getPicProgressPerStep(timeNs float64) float64 {
	return 1.0 / (float64(self.parameter.PictureChangeIntervallMs) / 1000) * (timeNs / 1000 / 1000 / 1000)
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
	self.colProgressPerStep = self.getColProgressPerStep(float64(self.display.GetRefreshIntervalNs()))
	self.picProgressPerStep = self.getPicProgressPerStep(float64(self.display.GetRefreshIntervalNs()))
}

func (self *ModePicture) SetParameter(parm ModePictureParameter) {
	self.parameter = parm
	self.GetDbDriver().Write(self.GetName(), "parameter", self.parameter)
	self.postSetParameter()
}
func (self *ModePicture) RandomizePreset() {
	self.Randomize()
}
func (self *ModePicture) Randomize() {
	self.picIndex = uint8(self.rand.Uint32() % uint32(len(picture.Pics)))
	parameter := ModePictureParameter{
		PictureColumnUs:          (self.rand.Uint32())%(self.limits.MaxPictureColumnUs-self.limits.MinPictureColumnUs) + self.limits.MinPictureColumnUs,
		PictureChangeIntervallMs: (self.rand.Uint32())%(self.limits.MaxPictureChangeIntervallMs-self.limits.MinPictureChangeIntervallMs) + self.limits.MinPictureChangeIntervallMs,
		Brightness:               self.rand.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
	}
	self.SetParameter(parameter)
}
