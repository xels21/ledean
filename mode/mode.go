// Package mode handels all the different led modes
package mode

import (
	"ledean/display"
	"time"

	"github.com/sdomino/scribble"
	log "github.com/sirupsen/logrus"
)

type Mode interface {
	Activate()
	Deactivate()
	Randomize()
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
	dbDriver    *scribble.Driver
	display     *display.Display
	renderType  RenderType
	name        string
	calcDisplay func()
	cExit       chan bool
}

func NewModeSuper(dbDriver *scribble.Driver, display *display.Display, name string, renderType RenderType, calcDisplay func()) *ModeSuper {
	self := ModeSuper{
		dbDriver:    dbDriver,
		display:     display,
		name:        name,
		renderType:  renderType,
		calcDisplay: calcDisplay,
	}

	if self.renderType == RenderTypeDynamic {
		self.cExit = make(chan bool, 1)
	}

	return &self
}

func (self *ModeSuper) GetName() string {
	return self.name
}

func (self *ModeSuper) Activate() {
	switch self.renderType {
	case RenderTypeStatic:
		self.calcDisplay()
		self.display.Render()
	case RenderTypeDynamic:
		ticker := time.NewTicker(RefreshIntervalNs)
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
	default:
		log.Debugf("unknown render type")
	}
}

func (self *ModeSuper) Deactivate() {
	if self.renderType == RenderTypeDynamic {
		self.cExit <- true
	}
}
