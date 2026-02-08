//go:build tinygo
// +build tinygo

package webserver

import (
	"ledean/log"
	"ledean/mode"
	"ledean/websocket"
)

func Start(addr string, port int, path2Frontend string, modeController *mode.ModeController, hub *websocket.Hub) {
	log.Error("Not possible with tinygo yet: webserver Start")
}
