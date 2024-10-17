package color

import (
	"ledean/helper"
	"ledean/json"
	"math"
)

type RGB struct {
	R byte `json:"r"`
	G byte `json:"g"`
	B byte `json:"b"`
}

func (self RGB) String() string {
	msg, err := json.Marshal(self)
	if err != nil {
		msg = []byte(err.Error())
	}
	return string(msg)
}

func (self *RGB) ToHsv() HSV {
	r := float64(self.R) / 255.0
	g := float64(self.G) / 255.0
	b := float64(self.B) / 255.0

	maxC := math.Max(math.Max(r, g), b)
	minC := math.Min(math.Min(r, g), b)

	return HSV{
		H: rgb2hue(r, b, g, maxC, minC),
		S: rgb2sHsv(r, b, g, maxC, minC),
		V: maxC,
	}
}

func rgb2sHsv(r float64, b float64, g float64, maxC float64, minC float64) float64 {
	if maxC == 0.0 {
		return 0.0
	}
	return ((maxC - minC) / maxC)
}

func rgb2hue(r float64, b float64, g float64, maxC float64, minC float64) float64 {
	deltaC := maxC - minC
	hue := float64(0)
	switch maxC {
	case r:
		hue = 60.0 * (0.0 + ((g - b) / deltaC))
		break
	case g:
		hue = 60.0 * (2.0 + ((b - r) / deltaC))
		break
	case b:
		hue = 60.0 * (4.0 + ((r - g) / deltaC))
		break
	}
	for hue < 0.0 {
		hue += 360.0
	}
	return hue
}

const (
	// BGR aka "Blue Green Red" is the current APA102 LED color order.
	SPI_ORDER_BGR = iota
	// BRG aka "Blue Red Green" is the typical APA102 color order from 2015-2017.
	SPI_ORDER_BRG
	// GRB aka "Green Red Blue" is the typical APA102 color order from pre-2015.
	// ALSO WS2812 144 leds/m
	SPI_ORDER_GRB

	SPI_ORDER_GBR

	SPI_ORDER_RGB

	SPI_ORDER_RBG
)

func OrderStr2int(ledOrder string) int {
	switch ledOrder {

	case "BGR":
		return SPI_ORDER_BGR
	case "BRG":
		return SPI_ORDER_BRG
	case "GRB":
		return SPI_ORDER_GRB
	case "GBR":
		return SPI_ORDER_GBR
	case "RGB":
		return SPI_ORDER_RGB
	case "RBG":
		return SPI_ORDER_RBG
	default:
		return SPI_ORDER_RGB //error
	}
}

// DEFAULT
//
//	func (self *RGB) ToSpi() []byte {
//		return []byte{self.R, self.G, self.B}
//	}

// POI
func (self *RGB) ToSpi(order int) []byte {
	switch order {

	case SPI_ORDER_BGR:
		return []byte{self.B, self.G, self.R}
	case SPI_ORDER_BRG:
		return []byte{self.B, self.R, self.G}
	case SPI_ORDER_GRB:
		return []byte{self.G, self.R, self.B}
	case SPI_ORDER_GBR:
		return []byte{self.G, self.B, self.R}
	case SPI_ORDER_RGB:
		return []byte{self.R, self.G, self.B}
	case SPI_ORDER_RBG:
		return []byte{self.R, self.B, self.G}
	default:
		return []byte{0, 0, 0} //error
	}
}

func (self *RGB) Add(toAdd RGB) {
	self.R = uint8(helper.MinInt16(int16(self.R)+int16(toAdd.R), 255))
	self.G = uint8(helper.MinInt16(int16(self.G)+int16(toAdd.G), 255))
	self.B = uint8(helper.MinInt16(int16(self.B)+int16(toAdd.B), 255))
}
func (self *RGB) Sub(toAdd RGB) {
	self.R = uint8(helper.MaxInt16(int16(self.R)-int16(toAdd.R), 0))
	self.G = uint8(helper.MaxInt16(int16(self.G)-int16(toAdd.G), 0))
	self.B = uint8(helper.MaxInt16(int16(self.B)-int16(toAdd.B), 0))
}
