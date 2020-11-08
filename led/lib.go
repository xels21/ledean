package led

import "LEDean/led/color"

func (self *LedController) Clear() {
	self.AllSolid(color.RGB{R: 0, G: 0, B: 0})
}

func (self *LedController) AllSolid(rgb color.RGB) {
	for i, _ := range self.leds {
		self.leds[i] = rgb
	}
}
