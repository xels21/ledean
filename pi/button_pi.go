// +build linux

package pi

import (
	"log"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

func (self *PiButton) register() {
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
	var risingMs, fallingMs int64

	lastState := p.Read()
	for p.WaitForEdge(-1) {
		currentState := p.Read()
		if lastState == currentState {
			continue
		}
		nowMs := time.Now().UnixNano() / 1000 /*us*/ / 1000 /*ns*/
		switch currentState {
		case gpio.Low:
			fallingMs = nowMs
		case gpio.High:
			risingMs = nowMs
		}

		if currentState == gpio.Low {
			passedTimeMs := fallingMs - risingMs
			if passedTimeMs < 50 /*ms*/ {
				continue
			}
			if passedTimeMs >= self.longPressMs {
				log.Println("longPress")
			} else {
				log.Println("shortPress")
			}
		}

		// log.Printf("%s went %s\n", p, currentState)
		lastState = currentState
	}
	self.register()
}
