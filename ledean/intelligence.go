package ledean

import (
	"LEDean/led"
	"LEDean/pi/button"
	pi "LEDean/pi/general"
	"LEDean/pi/ws28x"
	"LEDean/webserver"

	log "github.com/sirupsen/logrus"
)

type LEDeanInstance struct {
	pi_button         *button.PiButton
	pi_ws28xConnector *ws28x.PiWs28xConnector
	ledController     *led.LedController
}

func Run(parm *Parameter) LEDeanInstance {
	parm.Check()

	SetLogger(parm.LogLevel)

	pi.Init()
	pi_button := button.NewPiButton(parm.GpioButton, parm.PressLongMs, parm.PressDoubleTimeout)
	pi_button.Register()

	pi_ws28xConnector := ws28x.NewPiWs28xConnector(parm.SpiInfo)
	pi_ws28xConnector.Connect()

	ledController := led.NewLedController(parm.LedCount, pi_ws28xConnector, pi_button)

	pi_button.AddCbPressSingle(func() { log.Info("PRESS_SINGLE") })
	pi_button.AddCbPressDouble(func() { log.Info("PRESS_DOUBLE") })
	pi_button.AddCbPressLong(func() { log.Info("PRESS_LONG") })

	go webserver.Start(parm.Address, parm.Port, parm.Path2Frontend, ledController, pi_button)

	ledController.Start()

	return LEDeanInstance{
		pi_button:         pi_button,
		pi_ws28xConnector: pi_ws28xConnector,
		ledController:     ledController}
}

func RunForever() {
	log.Info("Running forever ...")
	for {
		select {}
	}
}
