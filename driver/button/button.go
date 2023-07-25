package button

import (
	"ledean/driver/pin"
	"ledean/log"
	"ledean/websocket"
	"time"
)

const DEBOUNCE_NS = 50 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/

type Button struct {
	gpio               string
	longPressMs        int
	pressDoubleTimeout int
	pin                *pin.Pin
	cbPressLong        []func()
	cbPressSingle      []func()
	cbPressDouble      []func()
	hub                *websocket.Hub
	pCmdButtonChannel  *chan websocket.CmdButton
}

func NewButton(gpio string, longPressMs int, pressDoubleTimeout int, hub *websocket.Hub) *Button {

	self := Button{
		longPressMs:        longPressMs,
		pressDoubleTimeout: pressDoubleTimeout,
		pin:                pin.NewPin(gpio),
		cbPressLong:        make([]func(), 0, 4),
		cbPressSingle:      make([]func(), 0, 4),
		cbPressDouble:      make([]func(), 0, 4),
		hub:                hub,
		pCmdButtonChannel:  hub.GetCmdButtonChannel(),
	}

	if hub != nil {
		go self.socketHandler()
	}

	go self.listen()

	return &self
}

func (self *Button) socketHandler() {
	for {
		cmdButton := <-*self.pCmdButtonChannel
		switch cmdButton.Action {
		case "single":
			self.PressSingle()
		case "double":
			self.PressDouble()
		case "long":
			self.PressLong()
		default:
			log.Info("Unknown button action: ", cmdButton.Action)
		}
	}
}

func (self *Button) listen() {
	var passedTimeNs, risingNs, fallingNs int64
	var inTime bool
	lastActionNs := int64(0)
	for {
		// Initial press can be triggered anywhere in time
		// button press should set the pin to High
		// Caused by debouncing behaviour, action has a timeout of {{DEBOUNCE_NS}}
		self.pin.WaitForEdge(-1)
		log.Trace("1_", self.pin.Read())

		if self.pin.Read() != true {
			continue
		}

		risingNs = time.Now().UnixNano()
		if risingNs < lastActionNs+int64(DEBOUNCE_NS) {
			continue
		}

		//labeled breakout because of complexity
	press_detection:
		for {
			//longpress timeout -> when no Edge is coming, user is still pressing
			inTime = self.pin.WaitForEdge(time.Millisecond * time.Duration(self.longPressMs))
			log.Trace("2_", self.pin.Read())

			if !inTime {
				self.PressLong()
				break press_detection
			}
			if self.pin.Read() != false {
				continue
			}

			fallingNs = time.Now().UnixNano()
			passedTimeNs = fallingNs - risingNs
			if passedTimeNs < DEBOUNCE_NS { // Debouncing
				risingNs = fallingNs
				continue
			}

			for {
				inTime = self.pin.WaitForEdge(time.Millisecond * time.Duration(self.pressDoubleTimeout))
				log.Trace("3_", self.pin.Read())

				if !inTime {
					self.PressSingle()
					break press_detection
				}
				if self.pin.Read() != true {
					continue
				}

				risingNs = time.Now().UnixNano()
				passedTimeNs = risingNs - fallingNs
				if passedTimeNs < DEBOUNCE_NS { // Debouncing
					fallingNs = risingNs
					continue
				}
				self.PressDouble()
				for {
					self.pin.WaitForEdge(-1)
					log.Trace("4_", self.pin.Read())
					if self.pin.Read() != false {
						continue
					}

					fallingNs = time.Now().UnixNano()
					passedTimeNs = fallingNs - risingNs
					if passedTimeNs < DEBOUNCE_NS { // Debouncing
						risingNs = fallingNs
						continue
					}

					break press_detection
				}
			}
		}
		lastActionNs = time.Now().UnixNano()
	}
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
