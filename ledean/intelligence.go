package ledean

import (
	"ledean/display"
	"ledean/mode"
	"ledean/pi/button"
	pi "ledean/pi/general"
	"ledean/webserver"

	scribble "github.com/sdomino/scribble"
	log "github.com/sirupsen/logrus"
)

type LEDeanInstance struct {
	dbDriver       *scribble.Driver
	display        *display.Display
	modeController *mode.ModeController
	piButton       *button.PiButton
}

func Run(parm *Parameter) *LEDeanInstance {
	var self LEDeanInstance
	var err error
	parm.Check()
	SetLogger(parm.LogLevel)

	self.dbDriver, err = scribble.New(parm.Path2DB, nil)
	if err != nil {
		log.Panic("Error while trying to make a new DB: ", err)
	}

	pi.Init()
	self.piButton = button.NewPiButton(parm.GpioButton, parm.PressLongMs, parm.PressDoubleTimeout)
	self.piButton.Register()

	self.display = display.NewDisplay(parm.LedCount, parm.LedRows, parm.SpiInfo, parm.ReverseRows)
	self.modeController = mode.NewModeController(self.dbDriver, self.display, self.piButton)

	self.piButton.AddCbPressSingle(func() { log.Info("PRESS_SINGLE") })
	self.piButton.AddCbPressDouble(func() { log.Info("PRESS_DOUBLE") })
	self.piButton.AddCbPressLong(func() { log.Info("PRESS_LONG") })

	go webserver.Start(parm.Address, parm.Port, parm.Path2Frontend, self.display, self.modeController, self.piButton)

	if parm.DirectStart {
		self.modeController.Start()
	}

	return &self
}

func RunForever() {
	log.Info("Running forever ...")
	for {
		select {}
	}
}
