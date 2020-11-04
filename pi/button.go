package pi

type PiButton struct {
	gpio        string
	longPressMs int64
}

func NewPiButton(gpio string, longPressMs int64) *PiButton {

	self := PiButton{
		gpio:        gpio,
		longPressMs: longPressMs,
	}

	self.register()

	return &self
}
