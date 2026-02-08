//go:build tinygo
// +build tinygo

package websocket

import (
	"ledean/log"
)

type Client struct{}

func (self *Client) SendCmd(cmd Cmd) {
	log.Error("Not possible with tinygo yet: client_stub - SendCmd")
}
