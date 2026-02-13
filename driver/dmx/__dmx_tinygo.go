//go:build false
// +build false

package dmx

import (
	"ledean/log"
	"machine"
	"time"
)

const (
	dmxBaud      = 250000
	startCode    = 0x00
	maxDmxRead   = 512 // max channel index to accept from wire
	maxListeners = 16
)

// dmxListener is a compact per-channel callback entry.
type dmxListener struct {
	chn      uint16
	prevVal  byte
	callback func(value byte)
}

type Dmx struct {
	rxPin         machine.Pin
	dmx           *machine.UART
	listeners     [maxListeners]dmxListener
	listenerCount uint8
	maxChn        uint16 // highest registered channel (for early exit)
}

func NewDmx() *Dmx {
	self := Dmx{
		rxPin: machine.GPIO20,
		dmx:   machine.UART1,
	}
	log.Debug("DMX adapter created (UART config deferred to Run)")
	return &self
}

func (self *Dmx) AddChnListener(chn int, callback func(value byte)) {
	if chn < 0 || chn >= maxDmxRead {
		log.Warningf("DMX: invalid channel %d (valid range 0-%d)", chn, maxDmxRead-1)
		return
	}
	if callback == nil {
		log.Warning("DMX: nil listener callback ignored")
		return
	}
	if self.listenerCount >= maxListeners {
		log.Warning("DMX: max listeners reached, ignoring")
		return
	}
	self.listeners[self.listenerCount] = dmxListener{
		chn:      uint16(chn),
		callback: callback,
	}
	self.listenerCount++
	if uint16(chn) > self.maxChn {
		self.maxChn = uint16(chn)
	}
}

func (self *Dmx) Run() {
	// Configure GPIO and UART here (not in NewDmx) so that UART RX interrupts
	// don't fire during the allocation-heavy init phase.
	self.rxPin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	self.dmx.Configure(machine.UARTConfig{
		BaudRate: dmxBaud,
		RX:       self.rxPin,
		TX:       machine.NoPin,
	})

	// Disable ALL UART interrupts. At 250 kbaud the RX-FIFO-full interrupt
	// fires every ~44 µs, which collides with GC/alloc in other goroutines
	// and causes "heap alloc in interrupt" panics on the ESP32-C3.
	// We poll the hardware FIFO directly instead.
	self.dmx.Bus.INT_ENA.Set(0)

	// self.drainUART()
	log.Debug("DMX UART configured (polling mode). Listening...")
	var b byte
	var err error
	for {
		// STEP 1: Detect DMX break via GPIO pin polling.
		// CRITICAL: drain UART buffer in every wait loop to prevent
		// the ring buffer from overflowing (which causes "blocked inside interrupt").
		_, ok := self.detectBreak()
		if !ok {
			continue
		}
		log.Debug("DMX break detected")

		// STEP 2: Break detected! We're in MAB (Mark After Break).
		// Flush UART - the break generated garbage bytes.
		self.drainUART()
		// self.dmx.Bus.INT_ENA.Set(1)

		// STEP 3: Read start code + channel data via UART hardware.
		// First byte after MAB is the start code.
		// sc, ok := self.readByteUART(2 * time.Millisecond)
		// if !ok {
		// 	continue
		// }
		time.Sleep(10 * time.Millisecond)
		b, err = self.dmx.ReadByte() // discard start code for now, since some fixtures don't send it
		if err != nil {
			log.Debug("DMX: failed to read start code")
			continue
		}
		if b != startCode {
			log.Debugf("DMX: invalid start code 0x%02X", b)
			continue
		}

		// STEP 4: Read channel data via UART and notify listeners inline.
		// DMX512 channels are 1-based: the first data byte after the start
		// code is channel 1. Read up to maxChn channels, then stop.
		chn := uint16(0)
		for chn <= self.maxChn {
			b, err = self.dmx.ReadByte() // discard start code for now, since some fixtures don't send it
			if err != nil {
				break // inter-byte timeout = end of this frame's data
			}
			self.notifyListeners(chn, b)
			chn++
		}
		self.drainUART()

	}
}

// notifyListeners checks all registered listeners for the given channel
// and fires callbacks only when the value has changed.
func (self *Dmx) notifyListeners(chn uint16, value byte) {
	for i := uint8(0); i < self.listenerCount; i++ {
		if self.listeners[i].chn == chn && self.listeners[i].prevVal != value {
			self.listeners[i].prevVal = value
			self.listeners[i].callback(value)
		}
	}
}

// drainUART reads and discards all bytes in the UART hardware FIFO.
func (self *Dmx) drainUART() {
	self.dmx.Buffer.Clear()
}

// detectBreak detects a DMX break signal: GPIO20 LOW for 80µs-1ms.
// Drains UART buffer in every polling loop to prevent interrupt panic.
func (self *Dmx) detectBreak() (time.Duration, bool) {
	// Wait for line HIGH (idle state)
	for self.rxPin.Get() {
		time.Sleep(10 * time.Microsecond) // avoid busy loop, but keep it tight enough to detect short break
	}

	startOfBreak := time.Now()
	// Wait for line to go LOW (break start)
	for !self.rxPin.Get() {
		time.Sleep(10 * time.Microsecond) // avoid busy loop, but keep it tight enough to detect short break
	}

	dur := time.Since(startOfBreak)
	// Valid DMX break: 88µs min (use 80µs margin), 1ms max
	if dur < 80*time.Microsecond || dur > time.Millisecond {
		return dur, false
	}
	return dur, true
}
