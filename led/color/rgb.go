package color

import "math"

type RGB struct {
	R byte `json:"r"`
	G byte `json:"g"`
	B byte `json:"b"`
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

func (self *RGB) ToSpi() []byte {
	var ret = make([]byte, 0, 9)
	return append(append(append(ret, Channel2Spi(self.R)...), Channel2Spi(self.B)...), Channel2Spi(self.G)...) //RGB -> RBG
}

//public for test purpose
func Channel2Spi(channel byte) []byte {

	var ret uint32 = 2396745 //001 001 001 001 001 001 001 001
	for i := uint32(0); i < 8; i++ {
		if (channel & (1 << i)) != 0 {
			ret |= 1 << (i*3 + 1)
		}
	}
	return []byte{
		byte(ret >> 16 & uint32(255)),
		byte(ret >> 8 & uint32(255)),
		byte(ret >> 0 & uint32(255))}
}
