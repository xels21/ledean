package picscaler

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

type PicScaler struct {
	inDir      string
	pixelCount int
	picNames   []string
	isReadDone bool
	asBytes    bool
	name       string
	outDir     string
}

func NewPicScaler(inDir string, outDir string, name string, pixelCount int, asBytes bool) *PicScaler {
	self := PicScaler{
		inDir:      inDir,
		pixelCount: pixelCount,
		asBytes:    asBytes,
		name:       name,
		outDir:     outDir,
	}

	return &self
}

// func (self *PicScaler) ClearOldFiles() {
// 	os.RemoveAll(self.outDir)
// }

func FilenameGoComform(filename string) string {
	filename = RemoveFileExtension(filename)
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, "-", "_")
	filename = strings.ReplaceAll(filename, "(", "_")
	filename = strings.ReplaceAll(filename, ")", "_")
	return filename
}

func (self *PicScaler) CreateController() {
	output, err := os.Create(filepath.Join(self.outDir, "pics_"+self.name+".go"))
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	fmt.Fprint(output, "//go:build poi\n")
	fmt.Fprint(output, "// +build poi\n\n")

	fmt.Fprint(output, "package "+self.name+"\n\n")
	if !self.asBytes {
		fmt.Fprint(output, "import \"image\"\n\n")
	}
	fmt.Fprint(output, "var PixelCount = "+strconv.Itoa(self.pixelCount)+"\n\n")

	if self.asBytes {
		fmt.Fprint(output, "var Pics = [][]string{")
	} else {
		fmt.Fprint(output, "var Pics = []*image.NRGBA{")
	}
	for _, picName := range self.picNames {
		picGo := FilenameGoComform(picName)
		if self.asBytes {
			fmt.Fprint(output, "\n\t"+self.name+"_"+picGo+",")
		} else {
			fmt.Fprint(output, "\n\t"+"&"+self.name+"_"+picGo+",")
		}
	}
	fmt.Fprint(output, "\n}\n")
	self.createControlleStub()
}

func (self *PicScaler) createControlleStub() {
	output, err := os.Create(filepath.Join(self.outDir, "pics_"+self.name+"_stub.go"))
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	fmt.Fprint(output, "//go:build !poi\n")
	fmt.Fprint(output, "// +build !poi\n\n")

	fmt.Fprint(output, "package "+self.name+"\n\n")
	// if !self.asBytes {
	// 	fmt.Fprint(output, "import \"image\"\n\n")
	// }
	fmt.Fprint(output, "var PixelCount = 0\n\n")

	if self.asBytes {
		fmt.Fprint(output, "var Pics = [][]string{")
	} else {
		fmt.Fprint(output, "var Pics = []*image.NRGBA{")
	}
	// for _, picName := range self.picNames {
	// 	picGo := FilenameGoComform(picName)
	// 	if self.asBytes {
	// 		fmt.Fprint(output, "\n\t"+self.name+"_"+picGo+",")
	// 	} else {
	// 		fmt.Fprint(output, "\n\t"+"&"+self.name+"_"+picGo+",")
	// 	}
	// }
	fmt.Fprint(output, "\n}\n")
}

func (self *PicScaler) Scale() {
	os.RemoveAll(self.outDir)
	err := os.Mkdir(self.outDir, os.ModeDir)
	if err != nil {
		log.Fatal(err)
	}
	self.ScaleToPixel()
}

func (self *PicScaler) readInDir() {
	entries, err := os.ReadDir(self.inDir)
	if err != nil {
		log.Fatal(err)
	}
	self.picNames = make([]string, 0, len(entries))
	for _, e := range entries {
		switch filepath.Ext(e.Name()) {
		case ".png", ".jpeg", ".jpg", ".bmp", ".gif":
			// if strings.HasPrefix(e.Name(), "_") {
			// continue
			// }
			self.picNames = append(self.picNames, e.Name())
			log.Print("Found image: " + e.Name())
		}
	}

}

func (self *PicScaler) ScaleToPixel() {
	if !self.isReadDone {
		self.readInDir()
	}
	for _, picName := range self.picNames {
		self.ScaleSingleToPixel(picName)
	}
}

func RemoveFileExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func (self *PicScaler) ScaleSingleToPixel(picName string) {
	// src, err := imgconv.Open(filepath.Join(self.inDir, picName))
	src, err := imaging.Open(filepath.Join(self.inDir, picName), imaging.AutoOrientation(true))
	// defer src.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Lanczos           - A high-quality resampling filter for photographic images yielding sharp results.
	// CatmullRom        - A sharp cubic filter that is faster than Lanczos filter while providing similar results.
	// MitchellNetravali - A cubic filter that produces smoother results with less ringing artifacts than CatmullRom.
	// Linear            - Bilinear resampling filter, produces smooth output. Faster than cubic filters.
	// Box               - Simple and fast averaging filter appropriate for downscaling. When upscaling it's similar to NearestNeighbor.
	// NearestNeighbor   - Fastest resampling filter, no antialiasing.
	resized := imaging.Resize(src, 0, self.pixelCount, imaging.Lanczos)
	resized = imaging.Sharpen(resized, 0.5)
	resized = imaging.AdjustContrast(resized, 20)
	rotated := imaging.Rotate270(resized)

	// Write the resulting image as TIFF.
	outFile := self.name + FilenameGoComform(picName) + ".jpg"
	os.MkdirAll(self.outDir, os.ModePerm)
	self.ConvertToGo(rotated, picName)

	err = imaging.Save(resized, filepath.Join(self.outDir, outFile))
	// err = imgconv.Write(output, resized, &imgconv.FormatOption{Format: imgconv.TIFF})
	if err != nil {
		log.Fatalf("failed to write image: %v", err)
	}
}

func NRGBAToGo(self *image.NRGBA) string {
	return fmt.Sprintf(`image.NRGBA{
	Pix:    []uint8{%s},
	Stride: %d,
	Rect:   image.Rect(%d, %d, %d, %d),
}`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(self.Pix)), ", "), "[]"), self.Stride, self.Rect.Min.X, self.Rect.Min.Y, self.Rect.Max.X, self.Rect.Max.Y)
}

// func NRGBAToStringArray(self *image.NRGBA) string {
// 	return fmt.Sprintf(`image.NRGBA{
// 	Pix:    []uint8{%s},
// 	Stride: %d,
// 	Rect:   image.Rect(%d, %d, %d, %d),
// }`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(self.Pix)), ", "), "[]"), self.Stride, self.Rect.Min.X, self.Rect.Min.Y, self.Rect.Max.X, self.Rect.Max.Y)
// }

func RgbaToRgbString(rgba []uint8, pixelPerRow int) string {
	pixelCount := len(rgba) / 4
	rowCount := pixelCount / pixelPerRow
	pixelAsString := "[]string{"
	for r := 0; r < rowCount; r++ {
		pixelAsString += "\n\t\""
		for i := 0; i < pixelPerRow; i++ {
			pixelAsString += fmt.Sprintf("\\x%02x\\x%02x\\x%02x", rgba[r*pixelPerRow*4+i*4+0], rgba[r*pixelPerRow*4+i*4+1], rgba[r*pixelPerRow*4+i*4+2])
			// pixelAsString += string(rgba[r*pixelPerRow*4+i*4 : r*pixelPerRow*4+i*4+3])
		}
		pixelAsString += "\","
	}
	return pixelAsString + "\n}\n"
}

func NRGBAToStringArray(self *image.NRGBA) string {
	return RgbaToRgbString(self.Pix, self.Rect.Max.X-self.Rect.Min.X)
}

func NRGBAToString(self *image.NRGBA) string {
	ret := "PicRGB{col: []PicCol{\n"
	// var pic_1 = PicRGB{col: []PicCol{PicCol{row: []color.RGB{
	// color.RGB{R: 1, G: 1, B: 1},
	// color.RGB{R: 1, G: 1, B: 1},
	// color.RGB{R: 1, G: 1, B: 1}},
	// }, PicCol{row: []color.RGB{
	// color.RGB{R: 1, G: 1, B: 1},
	// color.RGB{R: 1, G: 1, B: 1},
	// color.RGB{R: 1, G: 1, B: 1}},
	// }}}`+
	for y := range self.Rect.Max.Y {
		// output.Write("{\n")
		// fmt.Fprint(output, "PicCol{row:")
		ret += "	PicCol{row:"
		for x := range self.Rect.Max.X {
			off := y*self.Rect.Max.X*4 + x*4
			// fmt.Fprintf(output, " []color.RGB{%d,%d,%d},", self.Pix[off], self.Pix[off+1], self.Pix[off+2])
			ret += fmt.Sprintf(" []color.RGB{%d,%d,%d},", self.Pix[off], self.Pix[off+1], self.Pix[off+2])
			// output.Write("{}")
		}
		// output.Write("\n},")
		// fmt.Fprint(output, "},\n")
		ret += "},\n"

	}
	// fmt.Fprint(output, "}}")`
	ret += "	}\n}"
	return ret
}

func (self *PicScaler) ConvertToGo(resized *image.NRGBA, picName string) {
	/*
		format is:
		data.Pix -> Array R, G, B, A
		data.Rect.Max.X -> col
		data.Rect.Max.Y -> row
	*/

	picGo := FilenameGoComform(picName)
	output, err := os.Create(filepath.Join(self.outDir, strings.ToLower(self.name)+"_"+picGo+".go"))
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	fmt.Fprint(output, "//go:build poi\n")
	fmt.Fprint(output, "// +build poi\n\n")

	fmt.Fprint(output, "package "+self.name+"\n")
	if !self.asBytes {
		fmt.Fprint(output, `
import (
	"image"
)
`)
	}
	fmt.Fprint(output, "\nvar ", self.name, "_", picGo, " = ")
	if self.asBytes {
		fmt.Fprint(output, NRGBAToStringArray(resized))
	} else {
		fmt.Fprint(output, NRGBAToGo(resized))
	}
	fmt.Fprint(output, "\n")
}
