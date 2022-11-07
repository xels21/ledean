//go:build !tinygo
// +build !tinygo

package main

import (
	"fmt"
	"ledean/json"
	"ledean/ledean"

	"ledean/log"
)

func main() {
	fmt.Print(ledean.GetStartScreen())
	parm := ledean.GetParameter()
	jsonParm, err := json.MarshalIndent(parm, "", "\t")
	if err != nil {
		log.Panic("Could not convert params to JSON")
	}

	fmt.Print("Starting with:\n", string(jsonParm), "\n\n")

	ledean.Run(parm)
	ledean.RunForever()
}
