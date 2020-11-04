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
	pi_button := pi.NewPiButton(parameter.gpio_button, parameter.longPressMs)

	return LEDeanInstance{pi_button: pi_button}
}

func RunForever() {
	log.Info("Running forever ...")
	for {
	}
}
