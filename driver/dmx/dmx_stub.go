//go:build !tinygo
// +build !tinygo

package dmx

type Dmx struct {
}

func NewDmx() *Dmx {
	return nil
}

func (self *Dmx) AddChnListener(chn int, callback func(value byte)) {}

func (self *Dmx) Run() {}
