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
	// cliParameter := "-gpio_led_data=23 -gpio_button=24 -led_count=50 -direct_start"
	os.Args = []string{"tinygo_stub", "-gpio_led_data=23", "-gpio_button=22", "-led_count=50", "-direct_start"}
	parm := ledean.GetParameter()
	// parm := Parameter{}

	fmt.Print("Starting with:\n", os.Args, "\n\n")

	ledean.Run(parm)

	ledean.RunForever()
}
