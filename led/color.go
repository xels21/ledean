package led

type ColorRGB struct {
	R byte
	G byte
	B byte
}

func (self *ColorRGB) ToSpi() []byte {
	var ret = make([]byte, 0, 9)
	return append(append(append(ret, Channel2Spi(self.R)...), Channel2Spi(self.B)...), Channel2Spi(self.G)...) //RGB -> RBG
}

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
