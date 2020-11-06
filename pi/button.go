package pi

type PiButton struct {
	gpio               string
	longPressMs        int64
	doublePressTimeout int64
	CbLongPress        []func()
	CbSinglePress      []func()
	CbDoublePress      []func()
}

func NewPiButton(gpio string, longPressMs int64, doublePressTimeout int64) *PiButton {

	self := PiButton{
		gpio:               gpio,
		longPressMs:        longPressMs,
		doublePressTimeout: doublePressTimeout,
		CbLongPress:        make([]func(), 0, 4),
		CbSinglePress:      make([]func(), 0, 4),
		CbDoublePress:      make([]func(), 0, 4),
	}

	self.register()

	return &self
}
