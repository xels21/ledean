package color

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
