package ledean

import (
	"ledean/dbdriver"
	"ledean/display"
	"ledean/driver"
	"ledean/driver/button"
	"ledean/mode"
	"ledean/webserver"
	"ledean/websocket"
	"time"

	"ledean/log"
)

type LEDeanInstance struct {
	dbdriver       *dbdriver.DbDriver
	display        *display.Display
	modeController *mode.ModeController
	button         *button.Button
	hub            *websocket.Hub
}

func (self *LEDeanInstance) GetModeController() *mode.ModeController {
	return self.modeController
}

func Run(parm *Parameter) *LEDeanInstance {
	var self LEDeanInstance
	var err error
	parm.Check()
	log.SetLogger(parm.LogLevel)

	self.dbdriver, err = dbdriver.NewDbDriver(parm.Path2DB)
	if err != nil {
		log.Fatal("Error while trying to make a new DB: ", err)
	}

	driver.Init()

	if !parm.NoGui {
		self.hub = websocket.NewHub()
	}

	self.button = button.NewButton(self.dbdriver, parm.GpioButton, parm.PressLongMs, parm.PressDoubleTimeout, self.hub)

	self.display = display.NewDisplay(parm.LedCount, parm.LedRows, parm.GpioLedData, parm.ReverseRows, self.hub)
	if !parm.IsPictureMode {
		self.modeController = mode.NewModeController(self.dbdriver, self.display, self.button, self.hub)
	}

	self.button.AddCbPressSingle(func() { log.Info("PRESS_SINGLE") })
	self.button.AddCbPressDouble(func() { log.Info("PRESS_DOUBLE") })
	self.button.AddCbPressLong(func() { log.Info("PRESS_LONG") })

	if !parm.NoGui {
		go self.hub.Run()
		go webserver.Start(parm.Address, parm.Port, parm.Path2Frontend, self.modeController, self.button, self.hub)
	}

	if parm.DirectStart {
		if !parm.IsPictureMode {
			self.modeController.Start()
		}
		// self.modeController.NextMode()
	}

	const TEST_MODE = false
	if TEST_MODE {
		for {
			time.Sleep(3 * time.Second)
			self.modeController.NextMode()
		}
	}

	return &self
}

func RunForever() {
	log.Info("Running forever ...")
	for {
		select {}
	}
}
