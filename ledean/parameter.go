package ledean

import (
	"flag"
	"regexp"
	"strconv"
	"strings"

	"ledean/log"
)

type Parameter struct {
	GpioButton         string `json:"gpioButton"`
	GpioLedData        string `json:"gpioLedData"`
	IsShowMode         bool   `json:"isShowMode"`
	PressLongMs        int    `json:"pressLongMs"`
	PressDoubleTimeout int    `json:"pressDoubleTimeout"`
	LedCount           int    `json:"ledCount"`
	LedRows            int    `json:"ledRows"`
	DirectStart        bool   `json:"directStart"`
	LogLevel           string `json:"logLevel"`
	Path2Frontend      string `json:"path2Frontend"`
	Address            string `json:"address"`
	Port               int    `json:"port"`
	Fps                int    `json:"fps"`
	Path2DB            string `json:"path2Db"`
	ReverseRows        string `json:"reverseRows"`
	NoGui              bool   `json:"noGui"`
	LedOrder           string `json:"ledOrder"`
	LedDevice          string `json:"ledDevice"`

	// PathToDB       string
	// Pw             string
	// PathToFrontEnd string
}

func GetParameter() *Parameter {
	var parm Parameter
	flag.StringVar(&parm.GpioButton, "gpio_button", "", "gpio_pin for the button")
	flag.StringVar(&parm.GpioLedData, "gpio_led_data", "",
		`For uC:
	Define a Gpio pin
For RPi (using SPI):
		Info for spi communication.
	Leave it empty for following defaults [RPi 4 - SPI0.0]:
		- CLK : SPI0_CLK (GPIO11)
		- MOSI: SPI0_MOSI(GPIO10)
		- MISO: SPI0_MISO(GPIO9)
		- CS  : SPI0_CS0 (GPIO8)
	HINT: SPI0 has to be enabled in raspi-config first
	'''sudo raspi-config nonint do_spi 0'''
	`)
	flag.BoolVar(&parm.IsShowMode, "show_mode", false, "Wheter software is used for picture (POI) controll")
	flag.BoolVar(&parm.NoGui, "no_gui", false, "SW should not provide any gui (including website + websockets)")
	flag.IntVar(&parm.PressLongMs, "long_press_ms", 1200, "Time for the button long press")
	flag.IntVar(&parm.PressDoubleTimeout, "double_press_timeout", 350, "Time between single and double press")
	flag.IntVar(&parm.LedCount, "led_count", 0, "Amount of leds")
	flag.IntVar(&parm.LedRows, "led_rows", 1, "Amount of led rows")
	flag.BoolVar(&parm.DirectStart, "direct_start", false, "Wheter the display of the LEDs should be activated on startup")
	flag.StringVar(&parm.LogLevel, "log_level", "info", `log level. possibile: 
	- panic
	- fatal
	- error
	- warn or warning
	- info
	- debug
	`)
	flag.StringVar(&parm.Path2Frontend, "path2frontend", "", "path to static frontend. Keep it empty to dont serve static files")
	flag.StringVar(&parm.Address, "address", "0.0.0.0", "Local adress. Set it to '' to make the interface globally adressable")
	flag.IntVar(&parm.Port, "port", 2211, "Port for webserver")
	flag.IntVar(&parm.Fps, "fps", 40, "Display refresh rate. Should be between 1 and 200")
	flag.StringVar(&parm.Path2DB, "path2db", "db", "Path to DB folder (folder with json files)")
	flag.StringVar(&parm.ReverseRows, "reverse_rows", "0", "defines, which rows should be reversed (e.g. if second row is reversed: 0,1,0,0")
	flag.StringVar(&parm.LedOrder, "led_order", "RGB", "order of LED for SPI. e.g.: BGR|BRG|GRB|RGB")
	flag.StringVar(&parm.LedDevice, "led_device", "WS2812", "LED device protocol type: WS2812 | APA102")
	// pathToDB := flag.String("db", "../db", "Path to database folder")
	// pw := flag.String("pw", "", "Generate pw for user management")
	// pathToFrontEnd := flag.String("frontend", "../FrontRbc", "Path to frontend folder.\n"+
	// 	"If you don't want the Webserver to start, write the parameter but keep the value empty")
	flag.Parse()
	return &parm
}

func (self *Parameter) Check() {
	if self.LedCount <= 0 {
		log.Fatal("Error in parameter 'led_count'\n  - At least one led has to be connected")
	}
	if self.LedCount%self.LedRows != 0 {
		log.Fatal("Error in parameter 'led_count' and 'led_rows'\n  - Amount of led have to be equal to each row (e.g. led_count:20, led_rows:2, => 10 leds per row")
	}

	const TINYGO_SUPPORTS_REGEX = false
	// const TINYGO_SUPPORTS_REGEX = true
	if TINYGO_SUPPORTS_REGEX {
		if !regexp.MustCompile("^[01](,[01]){" + strconv.Itoa(self.LedRows-1) + "}$").MatchString(self.ReverseRows) {
			log.Fatal("Reverse Rows are set in a wrong way")
		}
	} else {
		commaCnt := strings.Count(self.ReverseRows, ",")
		numberCnt := strings.Count(self.ReverseRows, "0")
		numberCnt += strings.Count(self.ReverseRows, "1")
		if !(commaCnt == numberCnt-1 && numberCnt == self.LedRows) {
			log.Panic("Reverse Rows are set in a wrong way")
		}
	}
}
