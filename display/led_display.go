package display

import (
	"ledean/color"
	"ledean/json"
	"ledean/log"
	"ledean/websocket"
	"strings"
	"time"
)

const DISPLAY_TIMER_DELAY = 50

const (
	LED_DEVICE_WS2812 = iota

	LED_DEVICE_APA102
)

func LedDeviceStr2int(LedDevice string) int {
	switch LedDevice {

	case "WS2812":
		return LED_DEVICE_WS2812
	case "APA102":
		return LED_DEVICE_APA102

	default:
		return LED_DEVICE_WS2812 //error
	}
}

type DisplayBase struct {
	led_count            int
	led_rows             int
	leds_per_row         int
	order                int
	singleRowRGB         []color.RGB
	reversedSingleRowRGB []color.RGB
	reverse_rows         []bool
	leds                 []color.RGB
	buffer               []byte
	displayTimer         *time.Timer
	hub                  *websocket.Hub
	fps                  int //not used in display, but its central point for the modes
	refreshIntervalNs    time.Duration
}

func NewDisplayBase(led_count int, led_rows int, reverse_rows_raw string, fps int, order int, hub *websocket.Hub) *DisplayBase {
	self := DisplayBase{
		led_count:            led_count,
		led_rows:             led_rows,
		leds_per_row:         led_count / led_rows,
		order:                order,
		singleRowRGB:         make([]color.RGB, led_count/led_rows),
		reversedSingleRowRGB: make([]color.RGB, led_count/led_rows),
		reverse_rows:         make([]bool, led_rows),
		leds:                 make([]color.RGB, led_count),
		buffer:               make([]byte, 9*led_count),
		displayTimer:         time.NewTimer(DISPLAY_TIMER_DELAY * time.Millisecond),
		hub:                  hub,
		fps:                  fps,
		// active:    false,
	}
	if fps > 0 { //FPS==0 means highest refresh rate (delta time)
		self.refreshIntervalNs = time.Duration(1000 /*ms*/ *1000 /*us*/ *1000 /*ns*/ /fps) * time.Nanosecond
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
	if self.hub != nil {
		self.hub.AppendInitClientCb(self.initClientCb)
	}

	return &self
}

func (self *DisplayBase) GetFps() int {
	return self.fps
}

func (self *DisplayBase) GetRefreshIntervalNs() time.Duration {
	return self.refreshIntervalNs
}

func (self *DisplayBase) GetLedsJson() []byte {
	msg, err := json.Marshal(self.leds)
	if err != nil {
		msg = []byte(err.Error())
	}
	return msg
}

func (self *DisplayBase) GetRowLedCount() int {
	return self.leds_per_row
}

// func (self *DisplayBase) GetLeds() []color.RGB {
// 	return self.leds
// }
// func (self *DisplayBase) GetLedCount() int {
// 	return len(self.leds)
// }
// func (self *DisplayBase) GetLedRows() int {
// 	return self.led_rows
// }

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
	self.ledsChanged()
}

func (self *DisplayBase) initClientCb(client *websocket.Client) {
	var cmd2cLedsJSON, cmdLedsParameterJSON []byte
	var err error

	cmdLedsParameterJSON, err = json.Marshal(websocket.CmdLedsParameter{Rows: self.led_rows, Count: self.led_count})
	if err == nil {
		client.SendCmd(websocket.Cmd{Command: websocket.CmdLedsParameterId, Parameter: cmdLedsParameterJSON})
	}

	cmd2cLedsJSON, err = json.Marshal(websocket.CmdLeds{Leds: self.leds})
	if err == nil {
		client.SendCmd(websocket.Cmd{Command: websocket.CmdLedsId, Parameter: cmd2cLedsJSON})
	}
}

func (self *DisplayBase) ledsChanged() {
	select { //non blocking channels
	case <-self.displayTimer.C:
		self.displayTimer.Reset(DISPLAY_TIMER_DELAY * time.Millisecond)
		self.ForceLedsChanged()
	default:
		log.Trace("Leds got updated, but timer is still running")
	}
}

func (self *DisplayBase) ForceLedsChanged() {
	if self.hub != nil {
		cmd2cLedsJSON, err := json.Marshal(websocket.CmdLeds{Leds: self.leds})
		if err == nil {
			self.hub.Boradcast(websocket.Cmd{Command: websocket.CmdLedsId, Parameter: cmd2cLedsJSON})
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
		self.buffer = append(self.buffer, led.ToSpi(self.order)...)
	}
}

func (self *DisplayBase) _leds2Buffer() {
	needed := len(self.leds) * 3
	if cap(self.buffer) < needed {
		self.buffer = make([]byte, needed)
	} else {
		self.buffer = self.buffer[:needed]
	}
	idx := 0
	switch self.order {
	case color.SPI_ORDER_BGR:
		for i := range self.leds {
			led := self.leds[i]
			self.buffer[idx] = led.B
			self.buffer[idx+1] = led.G
			self.buffer[idx+2] = led.R
			idx += 3
		}
	case color.SPI_ORDER_BRG:
		for i := range self.leds {
			led := self.leds[i]
			self.buffer[idx] = led.B
			self.buffer[idx+1] = led.R
			self.buffer[idx+2] = led.G
			idx += 3
		}
	case color.SPI_ORDER_GRB:
		for i := range self.leds {
			led := self.leds[i]
			self.buffer[idx] = led.G
			self.buffer[idx+1] = led.R
			self.buffer[idx+2] = led.B
			idx += 3
		}
	case color.SPI_ORDER_GBR:
		for i := range self.leds {
			led := self.leds[i]
			self.buffer[idx] = led.G
			self.buffer[idx+1] = led.B
			self.buffer[idx+2] = led.R
			idx += 3
		}
	case color.SPI_ORDER_RBG:
		for i := range self.leds {
			led := self.leds[i]
			self.buffer[idx] = led.R
			self.buffer[idx+1] = led.B
			self.buffer[idx+2] = led.G
			idx += 3
		}
	case color.SPI_ORDER_RGB:
		fallthrough
	default:
		for i := range self.leds {
			led := self.leds[i]
			self.buffer[idx] = led.R
			self.buffer[idx+1] = led.G
			self.buffer[idx+2] = led.B
			idx += 3
		}
	}
}

func (self *DisplayBase) Clear() {
	self.AllSolid(color.RGB{R: 0, G: 0, B: 0})
	self.ForceLedsChanged() //due to start stop
}

func (self *DisplayBase) AllSolid(rgb color.RGB) {
	for i, _ := range self.singleRowRGB {
		self.singleRowRGB[i] = rgb
	}
	self.applySingleRow()
}
