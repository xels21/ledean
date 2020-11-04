package ledean

import "flag"

type Parameter struct {
	gpio_button        string
	gpio_spi_led       string
	longPressMs        int64
	doublePressTimeout int64
	// PathToDB       string
	// Pw             string
	// PathToFrontEnd string
}

func GetParameter() Parameter {
	gpio_button := flag.String("gpio_button", "GPIO17", "gpio_pin for the button")
	gpio_spi_led := flag.String("gpio_spi_led", "GPIO19", "gpio_pin for the spi MOSI led")
	longPressMs := flag.Int64("longPressMs", 1200, "Time for the button long press")
	doublePressTimeout := flag.Int64("doublePressTimeout", 350, "Time between single and double press")
	// pathToDB := flag.String("db", "../db", "Path to database folder")
	// pw := flag.String("pw", "", "Generate pw for user management")
	// pathToFrontEnd := flag.String("frontend", "../FrontRbc", "Path to frontend folder.\n"+
	// 	"If you don't want the Webserver to start, write the parameter but keep the value empty")
	// flag.Parse()
	return Parameter{
		gpio_button:        *gpio_button,
		gpio_spi_led:       *gpio_spi_led,
		longPressMs:        *longPressMs,
		doublePressTimeout: *doublePressTimeout,
		// 	PathToDB:       *pathToDB,
		// 	Pw:             *pw,
		// 	PathToFrontEnd: *pathToFrontEnd,
	}
}
