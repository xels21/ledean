//go:build tinygo
// +build tinygo

package webserver

import (
	"ledean/log"
	"ledean/websocket"
)

func Start(addr string, port int, path2Frontend string, hub *websocket.Hub) {
	log.Error("Not possible with tinygo yet")
}
