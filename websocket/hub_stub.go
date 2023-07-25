//go:build tinygo
// +build tinygo

// Package websocket provides the websocket implementation to avoid polling
package websocket

import (
	"ledean/log"
)

type Hub struct {
}

func NewHub() *Hub {
	log.Error("Not possible with tinygo yet")
	return nil
}

func (self *Hub) GetCmdButtonChannel() *chan CmdButton {
	log.Error("Not possible with tinygo yet")
	return nil
}
func (self *Hub) GetCmdModeActionChannel() *chan CmdModeAction {
	log.Error("Not possible with tinygo yet")
	return nil
}
func (self *Hub) GetCmdModeChannel() *chan CmdMode {
	log.Error("Not possible with tinygo yet")
	return nil
}

func (self *Hub) AppendInitClientCb(cb func(*Client)) {
	log.Error("Not possible with tinygo yet")
}

func (self *Hub) Run() {
	log.Error("Not possible with tinygo yet")
}

func (self *Hub) Boradcast(cmd Cmd) {
	log.Error("Not possible with tinygo yet")
}
