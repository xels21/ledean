package websocket

import (
	"ledean/display"
	"ledean/driver/button"
	"ledean/log"
	"ledean/mode"
)

type CommandController struct {
	display        *display.Display
	modeController *mode.ModeController
	button         *button.Button
}

func NewCommandController(display *display.Display, modeController *mode.ModeController, button *button.Button) *CommandController {
	return &CommandController{
		display:        display,
		modeController: modeController,
		button:         button}
}

func (self *CommandController) HandleCommand(cmd []byte) {
	log.Debug("got cmd: " + string(cmd))
}
