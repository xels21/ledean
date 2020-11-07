package led_test

import (
	"LEDean/led"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannel2Spi(t *testing.T) {
	toTest := led.Channel2Spi(0)
	// reflect.DeepEqual(toTest, [3]uint8{strconv.ParseInt("00100100", 2, 8), strconv.ParseInt("10010010", 2, 8), strconv.ParseInt("01001001", 2, 8)})
	assert.ElementsMatch(t, toTest, []uint8{uint8(36), uint8(146), uint8(73)})
	toTest = led.Channel2Spi(1)
	assert.ElementsMatch(t, toTest, []uint8{uint8(36), uint8(146), uint8(75)})
	toTest = led.Channel2Spi(128)
	assert.ElementsMatch(t, toTest, []uint8{uint8(100), uint8(146), uint8(73)})
	toTest = led.Channel2Spi(255)
	assert.ElementsMatch(t, toTest, []uint8{uint8(109), uint8(182), uint8(219)})
	// strconv.ParseInt("1001", 2, 64)
	// assert.(0, )
}

func TestToSpi(t *testing.T) {
	c := led.ColorRGB{R: 0, G: 0, B: 0}
	b := c.ToSpi()
	assert.ElementsMatch(t, b, []uint8{uint8(36), uint8(146), uint8(73), uint8(36), uint8(146), uint8(73), uint8(36), uint8(146), uint8(73)})
	c = led.ColorRGB{R: 255, G: 0, B: 0}
	b = c.ToSpi()
	assert.ElementsMatch(t, b, []uint8{uint8(109), uint8(182), uint8(219), uint8(36), uint8(146), uint8(73), uint8(36), uint8(146), uint8(73)})

}
