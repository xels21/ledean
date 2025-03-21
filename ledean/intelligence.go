package ledean

import (
	"ledean/color"
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

	self.display = display.NewDisplay(parm.LedCount, parm.LedRows, parm.GpioLedData, parm.ReverseRows, parm.Fps, color.OrderStr2int(parm.LedOrder), display.LedDeviceStr2int(parm.LedDevice), self.hub)

	self.modeController = mode.NewModeController(self.dbdriver, self.display, self.button, self.hub, parm.IsShowMode)

	if !parm.NoGui {
		go self.hub.Run()
		go webserver.Start(parm.Address, parm.Port, parm.Path2Frontend, self.modeController, self.hub)
	}

	if parm.DirectStart { //&& !parm.IsShowMode {
		self.modeController.Start()
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
