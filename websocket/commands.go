package websocket

import (
	"encoding/json"
	"ledean/color"
)

const (
	CmdLedsId          string = "leds"
	CmdLedsParameterId string = "ledsParameter"
	CmdButtonId        string = "button"
	Cmd2sModeId        string = "mode"
)

type Cmd struct {
	Command    string          `json:"cmd"`
	Parameters json.RawMessage `json:"parm"`
}

type CmdLeds struct {
	Leds []color.RGB `json:"leds"`
}

type CmdLedsParameter struct {
	Rows  int `json:"rows"`
	Count int `json:"count"`
}

type CmdButton struct {
	Action string `json:"action"`
}
