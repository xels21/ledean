package button

type PiButton struct {
	gpio               string
	longPressMs        int
	pressDoubleTimeout int
	cbPressLong        []func()
	cbPressSingle      []func()
	cbPressDouble      []func()
}

func NewPiButton(gpio string, longPressMs int, pressDoubleTimeout int) *PiButton {

	self := PiButton{
		gpio:               gpio,
		longPressMs:        longPressMs,
		pressDoubleTimeout: pressDoubleTimeout,
		cbPressLong:        make([]func(), 0, 4),
		cbPressSingle:      make([]func(), 0, 4),
		cbPressDouble:      make([]func(), 0, 4),
	}

	return &self
}

func (self *PiButton) AddCbPressSingle(cb func()) {
	self.cbPressSingle = append(self.cbPressSingle, cb)
}
func (self *PiButton) AddCbPressDouble(cb func()) {
	self.cbPressDouble = append(self.cbPressDouble, cb)
}
func (self *PiButton) AddCbPressLong(cb func()) {
	self.cbPressLong = append(self.cbPressLong, cb)
}

func trigger(cbs []func()) {
	for _, cb := range cbs {
		cb()
	}
}

func (self *PiButton) PressSingle() {
	trigger(self.cbPressSingle)
}
func (self *PiButton) PressDouble() {
	trigger(self.cbPressDouble)
}
func (self *PiButton) PressLong() {
	trigger(self.cbPressLong)
}
