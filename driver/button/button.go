package button

import (
	"ledean/dbdriver"
	"ledean/driver/pin"
	"ledean/json"
	"ledean/log"
	"ledean/websocket"
	"time"

	"errors"
)

const DEBOUNCE_NS = 50 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/

type Button struct {
	dbdriver           *dbdriver.DbDriver
	gpio               string
	longPressMs        int
	pressDoubleTimeout int
	isLocked           bool `json:"isLocked"`
	pin                *pin.Pin
	cbPressLong        []func()
	cbPressSingle      []func()
	cbPressDouble      []func()
	hub                *websocket.Hub
	pCmdButtonChannel  *chan websocket.CmdButton
}

func NewButton(dbdriver *dbdriver.DbDriver, gpio string, longPressMs int, pressDoubleTimeout int, hub *websocket.Hub) *Button {

	if gpio == "" {
		return nil
	}
	self := Button{
		dbdriver:           dbdriver,
		longPressMs:        longPressMs,
		pressDoubleTimeout: pressDoubleTimeout,
		isLocked:           false,
		pin:                pin.NewPin(gpio),
		cbPressLong:        make([]func(), 0, 4),
		cbPressSingle:      make([]func(), 0, 4),
		cbPressDouble:      make([]func(), 0, 4),
		hub:                hub,
		pCmdButtonChannel:  hub.GetCmdButtonChannel(),
	}

	err := dbdriver.Read("button", "isLocked", &self.isLocked)
	if err != nil {
		log.Info(err)
	}

	if hub != nil {
		go self.socketHandler()
		self.hub.AppendInitClientCb(self.initClientCb)
	}

	self.AddCbPressSingle(func() { log.Info("PRESS_SINGLE") })
	self.AddCbPressDouble(func() { log.Info("PRESS_DOUBLE") })
	self.AddCbPressLong(func() { log.Info("PRESS_LONG") })

	go self.listen()

	return &self
}

func (self *Button) initClientCb(client *websocket.Client) {
	err, cmd := self.GetCmdButtonIsLocked()
	if err != nil {
		log.Info(err)
		return
	}
	client.SendCmd(cmd)
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
		case "toggleLock":
			self.ToggleLock()
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
	if !self.isLocked {
		trigger(self.cbPressSingle)
	}
}
func (self *Button) PressDouble() {
	if !self.isLocked {
		trigger(self.cbPressDouble)
	}
}
func (self *Button) PressLong() {
	if !self.isLocked {
		trigger(self.cbPressLong)
	}
}
func (self *Button) ToggleLock() {
	self.isLocked = !self.isLocked
	self.dbdriver.Write("button", "isLocked", self.isLocked)
	self.BroadcastButtonIsLocked()
}

func (self *Button) BroadcastButtonIsLocked() {
	err, cmd := self.GetCmdButtonIsLocked()
	if err != nil {
		log.Info(err)
		return
	}
	self.hub.Boradcast(cmd)
}

func (self *Button) GetCmdButtonIsLocked() (error, websocket.Cmd) {
	if self.hub == nil {
		return errors.New("no hub"), websocket.Cmd{}
	}
	isLocked := "unlocked"
	if self.isLocked {
		isLocked = "locked"
	}
	cmdButtonJSON, err := json.Marshal(websocket.CmdButton{Action: isLocked})
	if err != nil {
		return err, websocket.Cmd{}
	} else {
		return nil, websocket.Cmd{
			Command:   websocket.CmdButtonId,
			Parameter: cmdButtonJSON}
	}
}
