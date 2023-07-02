//go:build tinygo
// +build tinygo

package webserver

import (
	"ledean/driver/button"
	"ledean/mode"
	"ledean/websocket"
)

func Start(addr string, port int, path2Frontend string, modeController *mode.ModeController, button *button.Button, hub *websocket.Hub) {
}
