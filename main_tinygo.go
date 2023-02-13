//go:build tinygo
// +build tinygo

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
	// os.Args = []string{"tinygo_stub", "-gpio_led_data=23", "-gpio_button=22", "-led_count=28", "-direct_start"}
	//Arduino Nano
	// os.Args = []string{"tinygo_stub", "-gpio_led_data=2", "-gpio_button=22", "-led_count=50", "-direct_start"}
	parm := ledean.GetParameter()

	fmt.Print("Starting with:\n", os.Args, "\n\n")

	ledean.Run(parm)

	ledean.RunForever()
}
