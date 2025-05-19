// Package mode handels all the different led modes
package mode

import (
	"ledean/dbdriver"
	"ledean/display"
	"math/rand"
	"time"

	"ledean/log"
)

type Mode interface {
	Activate()
	Deactivate()
	Randomize()
	RandomizePreset()
	TrySetParameter(b []byte) error
	GetParameter() interface{}
	GetLimits() interface{}
	GetName() string
}

type RenderType uint8

const (
	RenderTypeStatic  RenderType = iota
	RenderTypeDynamic RenderType = iota
)

type ModeSuper struct {
	dbdriver         *dbdriver.DbDriver
	display          *display.Display
	renderType       RenderType
	name             string
	calcDisplay      func()
	calcDisplayDelta func(deltaTimeNs int64)
	cExit            chan bool
	lastUpdateTime   time.Time
	rand             *rand.Rand
}

func NewModeSuper(dbdriver *dbdriver.DbDriver, display *display.Display, name string, renderType RenderType, calcDisplay func(), calcDisplayDelta func(deltaTimeNs int64), isRandDeterministic bool) *ModeSuper {
	self := ModeSuper{
		dbdriver:         dbdriver,
		display:          display,
		name:             name,
		renderType:       renderType,
		calcDisplay:      calcDisplay,
		calcDisplayDelta: calcDisplayDelta,
	}
	if isRandDeterministic {
		self.rand = rand.New(rand.NewSource(0))
	} else {
		self.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	if self.renderType == RenderTypeDynamic {
		self.cExit = make(chan bool, 1)
	}

	return &self
}

func (self *ModeSuper) GetName() string {
	return self.name
}

func (self *ModeSuper) GetDbDriver() *dbdriver.DbDriver {
	return self.dbdriver
}

func (self *ModeSuper) GetDisplay() *display.Display {
	return self.display
}

func (self *ModeSuper) Activate() {
	switch self.renderType {
	case RenderTypeStatic:
		self.calcDisplay()
		self.display.Render()
		self.display.ForceLedsChanged() //needed for delayed display render
	case RenderTypeDynamic:
		if self.display.GetFps() == 0 {
			go func() {
				self.lastUpdateTime = time.Now()
				var now time.Time
				var deltaTimeNs int64
				for {
					now = time.Now()
					deltaTimeNs = min(100 /*ms*/ *1000 /*us*/ *1000 /*ns*/, max(1, now.Sub(self.lastUpdateTime).Nanoseconds()))
					// log.Debugf("deltaTimeNs: %d", deltaTimeNs)
					self.lastUpdateTime = now
					select {
					case <-self.cExit:
						return
					default:
						self.calcDisplayDelta(deltaTimeNs)
						self.display.Render()
					}
				}
			}()
		} else {
			ticker := time.NewTicker(self.display.GetRefreshIntervalNs())
			go func() {
				for {
					select {
					case <-self.cExit:
						ticker.Stop()
						return
					case <-ticker.C:
						self.calcDisplay()
						self.display.Render()
					}
				}
			}()
		}
	default:
		log.Debugf("unknown render type")
	}
}

func (self *ModeSuper) Deactivate() {
	if self.renderType == RenderTypeDynamic {
		self.cExit <- true
	}
}
