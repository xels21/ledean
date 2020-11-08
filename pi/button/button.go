package button

type PiButton struct {
	gpio               string
	longPressMs        int64
	doublePressTimeout int64
	cbLongPress        []func()
	cbSinglePress      []func()
	cbDoublePress      []func()
}

func NewPiButton(gpio string, longPressMs int64, doublePressTimeout int64) *PiButton {

	self := PiButton{
		gpio:               gpio,
		longPressMs:        longPressMs,
		doublePressTimeout: doublePressTimeout,
		cbLongPress:        make([]func(), 0, 4),
		cbSinglePress:      make([]func(), 0, 4),
		cbDoublePress:      make([]func(), 0, 4),
	}

	return &self
}

func (self *PiButton) AddCbSinglePress(cb func()) {
	self.cbSinglePress = append(self.cbSinglePress, cb)
}
func (self *PiButton) AddCbDoublePress(cb func()) {
	self.cbDoublePress = append(self.cbDoublePress, cb)
}
func (self *PiButton) AddCbLongPress(cb func()) {
	self.cbLongPress = append(self.cbLongPress, cb)
}
