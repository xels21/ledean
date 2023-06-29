package websocket

import (
	"encoding/json"
	"ledean/color"
)

type Cmd struct {
	Command    string          `json:"cmd"`
	Parameters json.RawMessage `json:"parm"`
}

type Cmd2cLeds struct {
	Leds []color.RGB `json:"leds"`
}
