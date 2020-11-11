package mode

import "LEDean/led/color"

func GetAllModes(leds []color.RGB, cUpdate *chan bool) []Mode {
	modes := make([]Mode, 0, 16)
	modes = append(modes, NewModeRainBowSolid(leds, cUpdate))

	return modes
}
