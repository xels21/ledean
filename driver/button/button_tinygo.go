//go:build tinygo
// +build tinygo

package button

import (
	"machine"
	"strconv"
)

func (self *Button) Register() {
	gpio, _ := strconv.Atoi(self.gpio)
	pin := machine.Pin(gpio)
	pin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	// pin.SetInterrupt(machine.PinToggle, func(p machine.Pin) {
	// 	self.PressSingle()
	// })
	// go self.listen(p)

}

// func (self *Button) listen(p gpio.PinIO) {
// for {
// }
// }
