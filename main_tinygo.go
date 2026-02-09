//go:build tinygo && !poi
// +build tinygo,!poi

package main

import (
	"fmt"
	"ledean/ledean"
	"ledean/mode"
	"os"
)

func main() {
	fmt.Print(ledean.GetStartScreen())
	//ESP32 REICHI
	// os.Args = []string{"tinygo_stub", "-gpio_led_data=23", "-gpio_button=22", "-led_count=50", "-direct_start"}
	//ESP32 BABSTADT
	os.Args = []string{"tinygo_stub", "-gpio_led_data=0", "-led_count=50", "-no_gui", "-fps=10"}
	// os.Args = []string{"tinygo_stub", "-gpio_led_data=23", "-gpio_button=22", "-led_count=50", "-direct_start", "-no_gui"}
	//Arduino Nano
	// os.Args = []string{"tinygo_stub", "-gpio_led_data=2", "-gpio_button=22", "-led_count=50", "-direct_start"}
	parm := ledean.GetParameter()

	fmt.Print("Starting with:\n", os.Args, "\n\n")

	pLedean := ledean.Run(parm)
	injectFavMode(pLedean)

	ledean.RunForever()
}

func injectFavMode(pLedean *ledean.LEDeanInstance) {
	modeController := pLedean.GetModeController()
	modeSpectrum := modeController.GetModeSpectrum()
	modeSpectrum.SetParameter(mode.ModeSpectrumParameter{
		HueFrom720: 147.5,
		HueTo720:   308.6,
		// Brightness: 0.75,
		Brightness: 0.2,
		Positions: [2]mode.ModeSpectrumParameterPosition{
			{
				FacFrom:        2.0816143975048464,
				FacTo:          4.332379301570255,
				FacRoundTimeMs: 51440,
				OffFrom:        1.1849568322507775,
				OffTo:          2.9930981722250705,
				OffRoundTimeMs: 40118,
			},
			{
				FacFrom:        7,
				FacTo:          8.2,
				FacRoundTimeMs: 36921,
				OffFrom:        0.17155930631002372,
				OffTo:          0.8138194589841541,
				OffRoundTimeMs: 49997,
			}},
	})
	modeController.SwitchIndexFriendlyName("ModeSpectrum")
	// modeGradient := modeController.GetModeGradient()
	// modeGradient.SetParameter(mode.ModeGradientParameter{Brightness: 0.08, Count: 3, RoundTimeMs: 6666})
	// modeController.SwitchIndexFriendlyName("ModeGradient")
}
