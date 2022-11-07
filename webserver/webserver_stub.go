//go:build tinygo
// +build tinygo

package webserver

import (
	"ledean/display"
	"ledean/driver/button"
	"ledean/mode"
)

func Start(addr string, port int, path2Frontend string, display *display.Display, modeController *mode.ModeController, button *button.Button) {
}
