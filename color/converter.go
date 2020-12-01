package color

func RgbArr2HsvArr(leds []RGB) []HSV {
	hsvArr := make([]HSV, len(leds))
	for i, _ := range leds {
		hsvArr[i] = leds[i].ToHsv()
	}
	return hsvArr
}

func HsvArr2RgbArr(leds []HSV) []RGB {
	rgbArr := make([]RGB, len(leds))
	for i, _ := range leds {
		rgbArr[i] = leds[i].ToRGB()
	}
	return rgbArr
}
