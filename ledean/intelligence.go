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
	dbdriver       dbdriver.DbDriver
	display        display.Display
	modeController mode.ModeController
	button         button.Button
	hub            websocket.Hub
}

var leDeanInstance LEDeanInstance

func GetModeControllerStatic() *mode.ModeController {
	return &leDeanInstance.modeController
}

func (self *LEDeanInstance) GetModeController() *mode.ModeController {
	return &self.modeController
}

func Run(parm *Parameter) *LEDeanInstance {
	var btnPtr *button.Button
	var err error
	parm.Check()
	log.SetLogger(parm.LogLevel)

	leDeanInstance.dbdriver, err = dbdriver.NewDbDriver(parm.Path2DB)
	if err != nil {
		log.Fatal("Error while trying to make a new DB: ", err)
	}

	driver.Init()

	if !parm.NoGui {
		leDeanInstance.hub = websocket.NewHub()
	}

	leDeanInstance.button, err = button.NewButton(&leDeanInstance.dbdriver, parm.GpioButton, parm.PressLongMs, parm.PressDoubleTimeout, &leDeanInstance.hub)
	if err != nil {
		btnPtr = &leDeanInstance.button
	} else {
		btnPtr = nil
	}

	leDeanInstance.display = display.NewDisplay(parm.LedCount, parm.LedRows, parm.GpioLedData, parm.ReverseRows, parm.Fps, color.OrderStr2int(parm.LedOrder), display.LedDeviceStr2int(parm.LedDevice), &leDeanInstance.hub)

	leDeanInstance.modeController = mode.NewModeController(&leDeanInstance.dbdriver, &leDeanInstance.display, btnPtr, &leDeanInstance.hub, parm.IsPictureMode)
	leDeanInstance.modeController.PostCreate()

	if !parm.NoGui {
		go leDeanInstance.hub.Run()
		go webserver.Start(parm.Address, parm.Port, parm.Path2Frontend, &leDeanInstance.hub)
	}

	if parm.DirectStart { //&& !parm.IsPictureMode {
		leDeanInstance.modeController.Start()
		// leDeanInstance.modeController.NextMode()
	}

	const TEST_MODE = false
	if TEST_MODE {
		for {
			time.Sleep(3 * time.Second)
			leDeanInstance.modeController.NextMode()
		}
	}

	return &leDeanInstance
}

func RunForever() {
	log.Info("Running forever ...")
	for {
		select {}
	}
}
