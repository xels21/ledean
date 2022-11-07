package display

import (
	"ledean/color"
	"ledean/json"
	"strings"
)

type DisplayBase struct {
	led_count            int
	led_rows             int
	leds_per_row         int
	singleRowRGB         []color.RGB
	reversedSingleRowRGB []color.RGB
	reverse_rows         []bool
	leds                 []color.RGB
	buffer               []byte
	// active           bool
	// modeController *mode.ModeController //[]mode.Mode
}

func (self *DisplayBase) GetLedsJson() []byte {
	msg, err := json.Marshal(self.leds)
	if err != nil {
		msg = []byte(err.Error())
	}
	return msg
}

func NewDisplayBase(led_count int, led_rows int, reverse_rows_raw string) *DisplayBase {
	self := DisplayBase{
		led_count:            led_count,
		led_rows:             led_rows,
		leds_per_row:         led_count / led_rows,
		singleRowRGB:         make([]color.RGB, led_count/led_rows),
		reversedSingleRowRGB: make([]color.RGB, led_count/led_rows),
		reverse_rows:         make([]bool, led_rows),
		leds:                 make([]color.RGB, led_count),
		buffer:               make([]byte, 9*led_count),
		// active:    false,
	}

	reverse_rows_arr := strings.Split(reverse_rows_raw, ",")
	for i := 0; i < len(reverse_rows_arr); i++ {
		if reverse_rows_arr[i] == "1" {
			self.reverse_rows[i] = true
		}
	}

	// self.registerEvents()
	self.Clear()
	// go self.listen()

	return &self
}

func (self *DisplayBase) GetRowLedCount() int {
	return self.leds_per_row
}

// func (self *DisplayBase) GetLength() uint8 {
// 	return self.modeController.GetLength()
// }

// func (self *DisplayBase) GetIndex() uint8 {
// 	return self.modeController.GetIndex()
// }

// func (self *DisplayBase) GetModeRef(friendlyName string) (*mode.Mode, error) {
// 	return self.modeController.GetModeRef(friendlyName)
// }

// func (self *DisplayBase) GetModeResolver() []string {
// 	return self.modeController.GetModeResolver()
// }

// func (self *DisplayBase) IsActive() bool {
// 	return self.active
// }

func (self *DisplayBase) GetLeds() []color.RGB {
	return self.leds
}
func (self *DisplayBase) GetLedCount() int {
	return len(self.leds)
}
func (self *DisplayBase) GetLedRows() int {
	return self.led_rows
}

func reverseRgb(fromRGBs []color.RGB, toRGBs []color.RGB) {
	for i, j := 0, len(fromRGBs)-1; i < j; i, j = i+1, j-1 {
		toRGBs[i], toRGBs[j] = fromRGBs[j], fromRGBs[i]
	}
}

func (self *DisplayBase) applySingleRow() {
	var usedRow *[]color.RGB
	reverseRgb(self.singleRowRGB, self.reversedSingleRowRGB)
	for r := 0; r < self.led_rows; r++ {
		if self.reverse_rows[r] {
			usedRow = &self.reversedSingleRowRGB
		} else {
			usedRow = &self.singleRowRGB
		}
		for ri := 0; ri < self.leds_per_row; ri++ {
			i := r*self.leds_per_row + ri
			self.leds[i] = (*usedRow)[ri]
		}
	}
}
func (self *DisplayBase) ApplySingleRowRGB(singleRow []color.RGB) {
	self.singleRowRGB = singleRow
	self.applySingleRow()
}
func (self *DisplayBase) ApplySingleRowHSV(singleRow []color.HSV) {
	for i := 0; i < len(singleRow); i++ {
		self.singleRowRGB[i] = singleRow[i].ToRGB()
	}
	self.applySingleRow()
}

func (self *DisplayBase) leds2Buffer() {
	self.buffer = make([]byte, 0, 9*len(self.leds))
	for _, led := range self.leds {
		self.buffer = append(self.buffer, led.ToSpi()...)
	}
}

func (self *DisplayBase) Clear() {
	self.AllSolid(color.RGB{R: 0, G: 0, B: 0})
}

func (self *DisplayBase) AllSolid(rgb color.RGB) {
	for i, _ := range self.leds {
		self.leds[i] = rgb
	}
}
