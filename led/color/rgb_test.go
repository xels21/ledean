package color_test

import (
	"LEDean/led/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRgbToHsv(t *testing.T) {
	rgb := color.RGB{R: 255, G: 0, B: 0}
	expectedHsv := color.HSV{H: 0.0, S: 1.0, V: 1.0}
	hsv := rgb.ToHsv()
	assert.Equal(t, expectedHsv, hsv)

	rgb = color.RGB{R: 0, G: 255, B: 0}
	expectedHsv = color.HSV{H: 120.0, S: 1.0, V: 1.0}
	hsv = rgb.ToHsv()
	assert.Equal(t, expectedHsv, hsv)

	rgb = color.RGB{R: 0, G: 0, B: 255}
	expectedHsv = color.HSV{H: 240.0, S: 1.0, V: 1.0}
	hsv = rgb.ToHsv()
	assert.Equal(t, expectedHsv, hsv)

	rgb = color.RGB{R: 190, G: 20, B: 160}
	expectedHsv = color.HSV{H: 311.0, S: 0.89, V: 0.75}
	hsv = rgb.ToHsv()
	assert.InDelta(t, expectedHsv.H, hsv.H, 1.0)
	assert.InDelta(t, expectedHsv.S, hsv.S, 0.05)
	assert.InDelta(t, expectedHsv.V, hsv.V, 0.05)
}

func TestChannel2Spi(t *testing.T) {
	toTest := color.Channel2Spi(0)
	// reflect.DeepEqual(toTest, [3]uint8{strconv.ParseInt("00100100", 2, 8), strconv.ParseInt("10010010", 2, 8), strconv.ParseInt("01001001", 2, 8)})
	assert.ElementsMatch(t, toTest, []uint8{uint8(36), uint8(146), uint8(73)})
	toTest = color.Channel2Spi(1)
	assert.ElementsMatch(t, toTest, []uint8{uint8(36), uint8(146), uint8(75)})
	toTest = color.Channel2Spi(128)
	assert.ElementsMatch(t, toTest, []uint8{uint8(100), uint8(146), uint8(73)})
	toTest = color.Channel2Spi(255)
	assert.ElementsMatch(t, toTest, []uint8{uint8(109), uint8(182), uint8(219)})
	// strconv.ParseInt("1001", 2, 64)
	// assert.(0, )
}

func TestToSpi(t *testing.T) {
	c := color.RGB{R: 0, G: 0, B: 0}
	b := c.ToSpi()
	assert.ElementsMatch(t, b, []uint8{uint8(36), uint8(146), uint8(73), uint8(36), uint8(146), uint8(73), uint8(36), uint8(146), uint8(73)})
	c = color.RGB{R: 255, G: 0, B: 0}
	b = c.ToSpi()
	assert.ElementsMatch(t, b, []uint8{uint8(109), uint8(182), uint8(219), uint8(36), uint8(146), uint8(73), uint8(36), uint8(146), uint8(73)})

}
