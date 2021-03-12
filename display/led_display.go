package display

import (
	"encoding/json"
	"ledean/color"
	"ledean/pi/ws28x"
	"strings"
)

type Display struct {
	led_count            int
	led_rows             int
	leds_per_row         int
	singleRowRGB         []color.RGB
	reversedSingleRowRGB []color.RGB
	reverse_rows         []bool
	piWs28xConnector     *ws28x.PiWs28xConnector
	leds                 []color.RGB
	buffer               []byte
	// active           bool
	// modeController *mode.ModeController //[]mode.Mode
}

func NewDisplay(led_count int, led_rows int, spiInfo string, reverse_rows_raw string) *Display {
	self := Display{
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

	self.piWs28xConnector = ws28x.NewPiWs28xConnector(spiInfo)
	self.piWs28xConnector.Connect(led_count)

	// self.registerEvents()
	self.Clear()
	// go self.listen()

	return &self
}

func (self *Display) GetRowLedCount() int {
	return self.leds_per_row
}

// func (self *Display) GetLength() uint8 {
// 	return self.modeController.GetLength()
// }

// func (self *Display) GetIndex() uint8 {
// 	return self.modeController.GetIndex()
// }

// func (self *Display) GetModeRef(friendlyName string) (*mode.Mode, error) {
// 	return self.modeController.GetModeRef(friendlyName)
// }

// func (self *Display) GetModeResolver() []string {
// 	return self.modeController.GetModeResolver()
// }

// func (self *Display) IsActive() bool {
// 	return self.active
// }

func (self *Display) GetLeds() []color.RGB {
	return self.leds
}
func (self *Display) GetLedCount() int {
	return len(self.leds)
}
func (self *Display) GetLedRows() int {
	return self.led_rows
}

func (self *Display) GetLedsJson() []byte {
	msg, err := json.Marshal(self.leds)
	if err != nil {
		msg = []byte(err.Error())
	}
	return msg
}

func reverseRgb(fromRGBs []color.RGB, toRGBs []color.RGB) {
	for i, j := 0, len(fromRGBs)-1; i < j; i, j = i+1, j-1 {
		toRGBs[i], toRGBs[j] = fromRGBs[j], fromRGBs[i]
	}
}

func (self *Display) applySingleRow() {
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
func (self *Display) ApplySingleRowRGB(singleRow []color.RGB) {
	self.singleRowRGB = singleRow
	self.applySingleRow()
}
func (self *Display) ApplySingleRowHSV(singleRow []color.HSV) {
	for i := 0; i < len(singleRow); i++ {
		self.singleRowRGB[i] = singleRow[i].ToRGB()
	}
	self.applySingleRow()
}

func (self *Display) Render() {
	self.leds2Buffer()
	self.piWs28xConnector.Write(self.buffer)
}

func (self *Display) leds2Buffer() {
	self.buffer = make([]byte, 0, 9*len(self.leds))
	for _, led := range self.leds {
		self.buffer = append(self.buffer, led.ToSpi()...)
	}
}

func (self *Display) Clear() {
	self.AllSolid(color.RGB{R: 0, G: 0, B: 0})
}

func (self *Display) AllSolid(rgb color.RGB) {
	for i, _ := range self.leds {
		self.leds[i] = rgb
	}
}
