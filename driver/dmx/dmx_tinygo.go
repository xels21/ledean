//go:build tinygo
// +build tinygo

package dmx

import (
	"ledean/log"
	"machine"
	"runtime"
	"time"
)

const (
	dmxBaud    = 250000
	startCode  = 0x00
	maxDmxRead = 512 // max channel index to accept from wire
	// pollLoopsPerUs is a rough busy-wait calibration for ESP32-C3.
	// It avoids time.Now/After in interrupt-sensitive paths.
	pollLoopsPerUs = 8
	// maxListeners is the maximum number of per-channel callbacks.
	// 16 covers the 8 currently used and leaves room for growth.
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

	self.drainUART()
	log.Debug("DMX UART configured (polling mode). Listening...")

	for {
		// STEP 1: Detect DMX break via GPIO pin polling.
		// CRITICAL: drain UART buffer in every wait loop to prevent
		// the ring buffer from overflowing (which causes "blocked inside interrupt").
		_, ok := self.detectBreak()
		if !ok {
			continue
		}

		// STEP 2: Break detected! We're in MAB (Mark After Break).
		// Flush UART - the break generated garbage bytes.
		self.drainUART()

		// STEP 3: Read start code + channel data via UART hardware.
		// First byte after MAB is the start code.
		sc, ok := self.readByteUART(2 * time.Millisecond)
		if !ok {
			continue
		}
		if sc != startCode {
			// Avoid formatting/logging in tight loop to prevent allocations.
			continue
		}

		// STEP 4: Read channel data via UART and notify listeners inline.
		// DMX512 channels are 1-based: the first data byte after the start
		// code is channel 1. Read up to maxChn channels, then stop.
		chn := uint16(0)
		for chn <= self.maxChn {
			b, ok := self.readByteUART(1 * time.Millisecond)
			if !ok {
				break // inter-byte timeout = end of this frame's data
			}
			self.notifyListeners(chn, b)
			chn++
		}
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
	for self.dmx.Bus.GetSTATUS_RXFIFO_CNT() > 0 {
		self.dmx.Bus.GetFIFO_RXFIFO_RD_BYTE()
	}
}

// detectBreak detects a DMX break signal: GPIO20 LOW for 80µs-1ms.
// Drains UART buffer in every polling loop to prevent interrupt panic.
func (self *Dmx) detectBreak() (time.Duration, bool) {
	// Wait for line HIGH (idle state)
	if !self.waitForPinState(true, 200*time.Millisecond) {
		return 0, false
	}

	// Wait for line to go LOW (break start)
	if !self.waitForPinState(false, 200*time.Millisecond) {
		return 0, false
	}

	// Measure LOW duration (break length) with a bounded loop.
	maxLoops := loopsForTimeout(1 * time.Millisecond)
	loops := 0
	for !self.rxPin.Get() {
		loops++
		if loops >= maxLoops {
			break
		}
	}
	if loops == 0 {
		return 0, false
	}
	durUs := loops / pollLoopsPerUs
	dur := time.Duration(durUs) * time.Microsecond

	// Valid DMX break: 88µs min (use 80µs margin), 1ms max
	if dur < 80*time.Microsecond || dur > time.Millisecond {
		return dur, false
	}
	return dur, true
}

// readByteUART reads one byte from the UART hardware FIFO with a timeout.
// Uses direct register polling (no interrupts).
func (self *Dmx) readByteUART(timeout time.Duration) (byte, bool) {
	loops := loopsForTimeout(timeout)
	yieldCounter := 0
	for self.dmx.Bus.GetSTATUS_RXFIFO_CNT() == 0 {
		loops--
		if loops <= 0 {
			return 0, false
		}
		yieldCounter++
		if yieldCounter >= 1000 {
			yieldCounter = 0
			runtime.Gosched()
		}
	}
	return byte(self.dmx.Bus.GetFIFO_RXFIFO_RD_BYTE() & 0xFF), true
}

func loopsForTimeout(timeout time.Duration) int {
	if timeout <= 0 {
		return 1
	}
	us := int(timeout / time.Microsecond)
	if us < 1 {
		us = 1
	}
	return us * pollLoopsPerUs
}

func (self *Dmx) waitForPinState(state bool, timeout time.Duration) bool {
	loops := loopsForTimeout(timeout)
	yieldCounter := 0
	for self.rxPin.Get() != state {
		self.drainUART()
		loops--
		if loops <= 0 {
			return false
		}
		yieldCounter++
		if yieldCounter >= 1000 {
			yieldCounter = 0
			runtime.Gosched()
		}
	}
	return true
}
