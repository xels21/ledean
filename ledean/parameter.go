package ledean

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

type Parameter struct {
	GpioButton         string `json:"gpioButton"`
	SpiInfo            string `json:"spiInfo"`
	LongPressMs        int64  `json:"longPressMs"`
	DoublePressTimeout int64  `json:"doublePressTimeout"`
	LedCount           int64  `json:"ledCount"`
	LogLevel           string `json:"logLevel"`
	// PathToDB       string
	// Pw             string
	// PathToFrontEnd string
}

func GetParameter() *Parameter {
	var parm Parameter
	flag.StringVar(&parm.GpioButton, "gpio_button", "GPIO17", "gpio_pin for the button")
	flag.StringVar(&parm.SpiInfo, "spi_info", "",
		`Info for spi communication.
	Leave it empty for following defaults [RPi 4 - SPI0.0]:
		- CLK : SPI0_CLK (GPIO11)
		- MOSI: SPI0_MOSI(GPIO10)
		- MISO: SPI0_MISO(GPIO9)
		- CS  : SPI0_CS0 (GPIO8)
	HINT: SPI0 has to be enabled in raspi-config first
	'''sudo raspi-config nonint do_spi 0'''
	`)
	flag.Int64Var(&parm.LongPressMs, "long_press_ms", 1200, "Time for the button long press")
	flag.Int64Var(&parm.DoublePressTimeout, "double_press_timeout", 350, "Time between single and double press")
	flag.Int64Var(&parm.LedCount, "led_count", 0, "Amount of leds")
	flag.StringVar(&parm.LogLevel, "log_level", "info", `log level. possibile: 
	- panic
	- fatal
	- error
	- warn or warning
	- info
	- debug
	`)
	// pathToDB := flag.String("db", "../db", "Path to database folder")
	// pw := flag.String("pw", "", "Generate pw for user management")
	// pathToFrontEnd := flag.String("frontend", "../FrontRbc", "Path to frontend folder.\n"+
	// 	"If you don't want the Webserver to start, write the parameter but keep the value empty")
	flag.Parse()
	return &parm
}

func (self *Parameter) Check() {
	if self.LedCount <= 0 {
		log.Panic("Error in parameter 'led_count'\n  - At least one led has to be connected")
	}

}
