//go:build !tinygo
// +build !tinygo

package button

import "ledean/log"

func (self *Button) socketHandler() {
	for {
		cmdButton := <-self.hub.CmdButtonChannel
		switch cmdButton.Action {
		case "single":
			self.PressSingle()
		case "double":
			self.PressDouble()
		case "long":
			self.PressLong()
		default:
			log.Info("Unknown button action: ", cmdButton.Action)
		}
	}
}
