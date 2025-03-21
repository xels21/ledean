//go:build tinygo && poi
// +build tinygo,poi

package main

import (
	"fmt"
	"ledean/ledean"
	"os"
)

func main() {
	fmt.Print(ledean.GetStartScreen())
	//ESP32 REICHI
	// os.Args = []string{"tinygo_stub", "-gpio_led_data=23", "-gpio_button=22", "-led_count=50", "-direct_start"}
	//ESP32 BABSTADT
	// os.Args = []string{"tinygo_stub", "-gpio_led_data=0", "-led_count=50", "-direct_start", "-no_gui"}
	// os.Args = []string{"tinygo_stub", "-gpio_led_data=0", "-led_count=50", "-direct_start", "-no_gui", "-fps=250", "-led_order=GRB", "-show_mode"}
	// os.Args = []string{"tinygo_stub", "-gpio_led_data=0", "-led_count=50", "-direct_start", "-no_gui", "-fps=250", "-led_type=APA102", "-show_mode"}

	os.Args = []string{"tinygo_stub", "-gpio_led_data=0", "-led_count=58", "-direct_start", "-no_gui", "-fps=0", "-show_mode", "-led_order=GRB"}
	//Arduino Nano
	// os.Args = []string{"tinygo_stub", "-gpio_led_data=2", "-gpio_button=22", "-led_count=50", "-direct_start"}
	parm := ledean.GetParameter()

	fmt.Print("Starting with:\n", os.Args, "\n\n")

	ledean.Run(parm)
	// pLedean := ledean.Run(parm)
	// injectFavMode(pLedean)

	ledean.RunForever()
}

// func injectFavMode(pLedean *ledean.LEDeanInstance) {
// 	modeController := pLedean.GetModeController()
// 	modeGradient := modeController.GetModeGradient()
// 	modeGradient.SetParameter(mode.ModeGradientParameter{Brightness: 0.08, Count: 3, RoundTimeMs: 6666})
// 	modeController.SwitchIndexFriendlyName("ModeGradient")
// }
