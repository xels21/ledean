package websocket

import (
	"encoding/json"
	"ledean/color"
)

const (
	Cmd2cLedsId      string = "leds"
	Cmd2cLedsRowsId  string = "ledsRows"
	Cmd2cLedsCountId string = "ledsCount"
)

type Cmd struct {
	Command    string          `json:"cmd"`
	Parameters json.RawMessage `json:"parm"`
}

type Cmd2cLeds struct {
	Leds []color.RGB `json:"leds"`
}

type Cmd2cLedsRows struct {
	Rows int `json:"rows"`
}
type Cmd2cLedsCount struct {
	Count int `json:"count"`
}
