//go:build tinygo
// +build tinygo

package websocket

type Client struct{}

func (self *Client) SendCmd(cmd Cmd) {}
