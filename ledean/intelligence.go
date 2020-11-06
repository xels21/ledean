package ledean

import (
	"LEDean/pi"

	log "github.com/sirupsen/logrus"
)

type LEDeanInstance struct {
	pi_button *pi.PiButton
}

func Run(parameter Parameter) LEDeanInstance {
	pi.Init()
	pi_button := pi.NewPiButton(parameter.gpio_button, parameter.longPressMs, parameter.doublePressTimeout)

	pi_button.CbSinglePress = append(pi_button.CbSinglePress, func() { log.Info("PRESS_SINGLE") })
	pi_button.CbDoublePress = append(pi_button.CbDoublePress, func() { log.Info("PRESS_DOUBLE") })
	pi_button.CbLongPress = append(pi_button.CbLongPress, func() { log.Info("PRESS_Long") })

	return LEDeanInstance{pi_button: pi_button}
}

func RunForever() {
	log.Info("Running forever ...")
	for {
	}
}
