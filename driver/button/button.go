package button

type Button struct {
	gpio               string
	longPressMs        int
	pressDoubleTimeout int
	cbPressLong        []func()
	cbPressSingle      []func()
	cbPressDouble      []func()
}

func NewButton(gpio string, longPressMs int, pressDoubleTimeout int) *Button {

	self := Button{
		gpio:               gpio,
		longPressMs:        longPressMs,
		pressDoubleTimeout: pressDoubleTimeout,
		cbPressLong:        make([]func(), 0, 4),
		cbPressSingle:      make([]func(), 0, 4),
		cbPressDouble:      make([]func(), 0, 4),
	}

	return &self
}

func (self *Button) AddCbPressSingle(cb func()) {
	self.cbPressSingle = append(self.cbPressSingle, cb)
}
func (self *Button) AddCbPressDouble(cb func()) {
	self.cbPressDouble = append(self.cbPressDouble, cb)
}
func (self *Button) AddCbPressLong(cb func()) {
	self.cbPressLong = append(self.cbPressLong, cb)
}

func trigger(cbs []func()) {
	for _, cb := range cbs {
		cb()
	}
}

func (self *Button) PressSingle() {
	trigger(self.cbPressSingle)
}
func (self *Button) PressDouble() {
	trigger(self.cbPressDouble)
}
func (self *Button) PressLong() {
	trigger(self.cbPressLong)
}
