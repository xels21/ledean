package ledean

import (
	"LEDean/led"
	"LEDean/pi/button"
	pi "LEDean/pi/general"
	"LEDean/pi/ws28x"

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
	pi_button := button.NewPiButton(parm.GpioButton, parm.LongPressMs, parm.DoublePressTimeout)
	pi_button.Register()

	pi_ws28xConnector := ws28x.NewPiWs28xConnector(parm.SpiInfo)
	pi_ws28xConnector.Connect()

	ledController := led.NewLedController(parm.LedCount, pi_ws28xConnector, pi_button)

	pi_button.CbSinglePress = append(pi_button.CbSinglePress, func() { log.Info("PRESS_SINGLE") })
	pi_button.CbDoublePress = append(pi_button.CbDoublePress, func() { log.Info("PRESS_DOUBLE") })
	pi_button.CbLongPress = append(pi_button.CbLongPress, func() { log.Info("PRESS_LONG") })

	return LEDeanInstance{
		pi_button:         pi_button,
		pi_ws28xConnector: pi_ws28xConnector,
		ledController:     ledController}
}

func RunForever() {
	log.Info("Running forever ...")
	for {
	}
}
