package pi

type PiButton struct {
	gpio               string
	longPressMs        int64
	doublePressTimeout int64
}

func NewPiButton(gpio string, longPressMs int64, doublePressTimeout int64) *PiButton {

	self := PiButton{
		gpio:               gpio,
		longPressMs:        longPressMs,
		doublePressTimeout: doublePressTimeout,
	}

	self.register()

	return &self
}
