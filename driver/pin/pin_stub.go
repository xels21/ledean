//go:build !tinygo && !pi
// +build !tinygo,!pi

package pin

type Pin struct {
}

func NewPin(gpio string) *Pin {
	self := Pin{}
	return &self
}

func (self *Pin) WaitForEdge(timeout time.Duration) bool {
	return false
}

func (self *Pin) Read() bool {
	return false
}
