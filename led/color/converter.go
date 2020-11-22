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

func RgbArrClear(leds []RGB) {
	blackRgb := RGB{R: 0, G: 0, B: 0}
	for i, _ := range leds {
		leds[i] = blackRgb
	}
}

func HsvArrClear(leds []HSV) {
	blackHsv := HSV{H: 0.0, S: 0.0, V: 0.0}
	for i, _ := range leds {
		leds[i] = blackHsv
	}
}
