package main

import (
	"encoding/json"
	"fmt"
	"ledean/ledean"

	log "github.com/sirupsen/logrus"
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
