package websocket

import (
	"ledean/color"
	"ledean/json"
)

const (
	CmdLedsId          string = "leds"
	CmdLedsParameterId string = "ledsParameter"
	CmdButtonId        string = "button"
	CmdModeId          string = "mode"
	CmdModeLimitsId    string = "modeLimits"
	CmdModeResolverId  string = "modeResolver"
	CmdModeActionId    string = "action"
)

type Cmd struct {
	Command string `json:"cmd"`
	// Parameter any    `json:"parm"`
	Parameter json.RawMessage `json:"parm"` //must be json
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

type CmdMode struct {
	Id        string          `json:"id"`
	Parameter json.RawMessage `json:"parm"`
}

type CmdModeLimits struct {
	Id string `json:"id"`
	// Limits any    `json:"limits"`
	Limits json.RawMessage `json:"limits"`
}

type CmdModeResolver struct {
	Modes []string `json:"modes"`
}

const (
	CmdModeActionRandomize string = "randomize"
)

type CmdModeAction struct {
	Action string `json:"action"`
}
