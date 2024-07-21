//go:build tinygo && poi
// +build tinygo,poi

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
	os.Args = []string{"tinygo_stub", "-gpio_led_data=26", "-led_count=50", "-direct_start", "-no_gui"}
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
	modeGradient := modeController.GetModeGradient()
	modeGradient.SetParameter(mode.ModeGradientParameter{Brightness: 0.08, Count: 3, RoundTimeMs: 6666})
	modeController.SwitchIndexFriendlyName("ModeGradient")
}
