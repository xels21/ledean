//go:build pi
// +build pi

package pin

type Pin struct {
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
}

func NewPin(gpio string) *Pin {
	pin := gpioreg.ByName(gpio)
	if pin == nil {
		log.Fatal("Failed to find: ", gpio)
	}
	self := Pin{pin: pin}
	// var p PinIn
	if err := self.pin.In(gpio.PullDown, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}

	return &self
}

func (self *Pin) WaitForEdge(timeout time.Duration) bool {
	return self.pin.WaitForEdge(timeout)
}

func (self *Pin) Read() bool {
	return self.pin.Read()
}
