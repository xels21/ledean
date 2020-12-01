package color

type HSV struct {
	H float64 //0..360 deg
	S float64 //0..100 %
	V float64 //0..100 %
}

func (self *HSV) ToRGB() RGB {
	var r, g, b float64

	i := uint16(self.H / 60.0)
	f := self.H/60.0 - float64(i)
	p := self.V * (1.0 - self.S)
	q := self.V * (1.0 - self.S*f)
	t := self.V * (1.0 - self.S*(1.0-f))

	switch i % 6 {
	case 0:
		r = self.V
		g = t
		b = p
		break
	case 1:
		r = q
		g = self.V
		b = p
		break
	case 2:
		r = p
		g = self.V
		b = t
		break
	case 3:
		r = p
		g = q
		b = self.V
		break
	case 4:
		r = t
		g = p
		b = self.V
		break
	case 5:
		r = self.V
		g = p
		b = q
		break
	}

	return RGB{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
	}
}

func (self *HSV) Add(toAdd HSV) {
	self.AddRgb(toAdd.ToRGB())
}

func (self *HSV) AddRgb(toAdd RGB) {
	selfRgb := self.ToRGB()
	selfRgb.Add(toAdd)
	*self = selfRgb.ToHsv()
}

func (self *HSV) Sub(toSub HSV) {
	self.SubRgb(toSub.ToRGB())
}

func (self *HSV) SubRgb(toSub RGB) {
	selfRgb := self.ToRGB()
	selfRgb.Sub(toSub)
	*self = selfRgb.ToHsv()
}
