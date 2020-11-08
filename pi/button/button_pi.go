// +build linux

package button

import (
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"

	log "github.com/sirupsen/logrus"
)

const DEBOUNCE_NS = 50 /*ms*/ * 1000 /*us*/ * 1000 /*ns*/

func (self *PiButton) Register() {
	p := gpioreg.ByName(self.gpio)
	if p == nil {
		log.Fatal("Failed to find: ", self.gpio)
	}
	// var p PinIn
	if err := p.In(gpio.PullDown, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}

	go self.listen(p)
}

func (self *PiButton) listen(p gpio.PinIO) {
	var passedTimeNs, risingNs, fallingNs int64
	var inTime bool
	lastActionNs := int64(0)
	for {
		// Initial press can be triggered anywhere in time
		// button press should set the pin to High
		// Caused by debouncing behaviour, action has a timeout of {{DEBOUNCE_NS}}
		p.WaitForEdge(-1)
		log.Trace("1_", p.Read())

		if p.Read() != gpio.High {
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
			inTime = p.WaitForEdge(time.Millisecond * time.Duration(self.longPressMs))
			log.Trace("2_", p.Read())

			if !inTime {
				trigger(self.cbLongPress)
				break press_detection
			}
			if p.Read() != gpio.Low {
				continue
			}

			fallingNs = time.Now().UnixNano()
			passedTimeNs = fallingNs - risingNs
			if passedTimeNs < DEBOUNCE_NS { // Debouncing
				risingNs = fallingNs
				continue
			}

			for {
				inTime = p.WaitForEdge(time.Millisecond * time.Duration(self.doublePressTimeout))
				log.Trace("3_", p.Read())

				if !inTime {
					trigger(self.cbSinglePress)
					break press_detection
				}
				if p.Read() != gpio.High {
					continue
				}

				risingNs = time.Now().UnixNano()
				passedTimeNs = risingNs - fallingNs
				if passedTimeNs < DEBOUNCE_NS { // Debouncing
					fallingNs = risingNs
					continue
				}
				trigger(self.cbDoublePress)
				for {
					p.WaitForEdge(-1)
					log.Trace("4_", p.Read())
					if p.Read() != gpio.Low {
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
	self.Register()
}

func trigger(cbs []func()) {
	for _, cb := range cbs {
		cb()
	}
}
