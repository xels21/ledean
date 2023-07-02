//go:build tinygo
// +build tinygo

package websocket

type Hub struct{}

func NewHub() *Hub {
	return &Hub{}
}

func (self *Hub) Run() {}

func (self *Hub) AppendInitClientCb(cb func(*Client)) {}

func (self *Hub) Boradcast(cmd Cmd) {}
