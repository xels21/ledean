//go:build tinygo
// +build tinygo

package dmx

import (
	"ledean/log"
	"machine"
	"time"
)

const (
	dmxBaud     = 250000
	dmxChannels = 512 // DMX512 standard - will work also with less
	startCode   = 0x00
)

type Dmx struct {
	// 	pin machine.Pin
	frame        [dmxChannels]byte
	prevFrame    [dmxChannels]byte
	chnListeners [dmxChannels][]func(value byte)
	rxPin        machine.Pin
	dmx          *machine.UART
}

func NewDmx() *Dmx {
	self := Dmx{
		frame:     [dmxChannels]byte{},
		prevFrame: [dmxChannels]byte{},
		rxPin:     machine.GPIO20,
		dmx:       machine.UART1,
	}
	log.Debug("Initializing DMX adapter on GPIO20 (UART1 RX)")

	// Configure GPIO20 as input first (for break detection)
	self.rxPin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// Configure UART1 on the same pin for frame data
	self.dmx.Configure(machine.UARTConfig{
		BaudRate: dmxBaud,
		RX:       self.rxPin,
		TX:       machine.NoPin,
	})

	// Drain any buffered bytes from UART startup
	self.drainUART()

	log.Debug("UART configured. Waiting for DMX...")
	return &self
}

func (self *Dmx) GetFrame() [dmxChannels]byte {
	return self.frame
}

func (self *Dmx) AddChnListener(chn int, callback func(value byte)) {
	if chn < 0 || chn >= dmxChannels {
		log.Warningf("DMX: invalid channel %d (valid range 0-%d)", chn, dmxChannels-1)
		return
	}
	if callback == nil {
		log.Warning("DMX: nil listener callback ignored")
		return
	}
	idx := chn
	self.chnListeners[idx] = append(self.chnListeners[idx], callback)
}

func (self *Dmx) Run() {
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
			log.Debugf("Invalid DMX start code: 0x%02X", sc)
			continue
		}

		// STEP 4: Read channel data via UART
		received := 0
		for received < dmxChannels {
			b, ok := self.readByteUART(1 * time.Millisecond)
			if !ok {
				break // inter-byte timeout = end of this frame's data
			}
			self.frame[received] = b
			received++
		}
		log.Debugf("Received DMX frame: %d channels", received)

		self.checkChanges(received)
	}
}

// drainUART reads and discards all bytes currently in the UART buffer.
// Must be called frequently during GPIO polling loops to prevent overflow.
func (self *Dmx) drainUART() {
	for self.dmx.Buffered() > 0 {
		self.dmx.ReadByte()
	}
}

// detectBreak detects a DMX break signal: GPIO20 LOW for 80µs-1ms.
// Drains UART buffer in every polling loop to prevent interrupt panic.
func (self *Dmx) detectBreak() (time.Duration, bool) {
	// Wait for line HIGH (idle state)
	deadline := time.Now().Add(200 * time.Millisecond)
	for !self.rxPin.Get() {
		self.drainUART()
		if time.Now().After(deadline) {
			return 0, false
		}
	}

	// Wait for line to go LOW (break start)
	deadline = time.Now().Add(200 * time.Millisecond)
	for self.rxPin.Get() {
		self.drainUART()
		if time.Now().After(deadline) {
			return 0, false
		}
	}

	// Measure LOW duration (break length)
	start := time.Now()
	for !self.rxPin.Get() {
		if time.Since(start) > time.Millisecond {
			break
		}
	}
	dur := time.Since(start)

	// Valid DMX break: 88µs min (use 80µs margin), 1ms max
	if dur < 80*time.Microsecond || dur > time.Millisecond {
		return dur, false
	}
	return dur, true
}

// readByteUART reads one byte from UART with a timeout.
func (self *Dmx) readByteUART(timeout time.Duration) (byte, bool) {
	deadline := time.Now().Add(timeout)
	for self.dmx.Buffered() == 0 {
		if time.Now().After(deadline) {
			return 0, false
		}
	}
	b, err := self.dmx.ReadByte()
	if err != nil {
		return 0, false
	}
	return b, true
}

func (self *Dmx) checkChanges(length int) {
	changedChannels := []int{}

	for i := 0; i < length; i++ {
		if self.frame[i] != self.prevFrame[i] {
			changedChannels = append(changedChannels, i+1) // DMX channels are 1-indexed
			if len(self.chnListeners[i]) > 0 {
				value := self.frame[i]
				for _, callback := range self.chnListeners[i] {
					callback(value)
				}
			}
			self.prevFrame[i] = self.frame[i]
		}
	}

	for _, chn := range changedChannels {
		for _, cb := range self.chnListeners[chn] {
			cb(self.frame[chn])
		}
	}
}
