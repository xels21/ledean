package color_test

import (
	"ledean/color"
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

// func TestChannel2Spi(t *testing.T) {
// 	toTest := color.Channel2Spi(0)
// 	assert.ElementsMatch(t, toTest, []uint8{uint8(36), uint8(146), uint8(73)})
// 	toTest = color.Channel2Spi(1)
// 	assert.ElementsMatch(t, toTest, []uint8{uint8(36), uint8(146), uint8(75)})
// 	toTest = color.Channel2Spi(128)
// 	assert.ElementsMatch(t, toTest, []uint8{uint8(100), uint8(146), uint8(73)})
// 	toTest = color.Channel2Spi(255)
// 	assert.ElementsMatch(t, toTest, []uint8{uint8(109), uint8(182), uint8(219)})
// }

func TestToSpi(t *testing.T) {
	cR := uint8(50)
	cG := uint8(100)
	cB := uint8(150)
	c := color.RGB{R: cR, G: cG, B: cB}
	b := c.ToSpi(color.SPI_ORDER_RGB)
	assert.ElementsMatch(t, b, []uint8{cR, cG, cB})
	b = c.ToSpi(color.SPI_ORDER_RBG)
	assert.ElementsMatch(t, b, []uint8{cR, cB, cG})
	b = c.ToSpi(color.SPI_ORDER_GRB)
	assert.ElementsMatch(t, b, []uint8{cG, cR, cB})
	b = c.ToSpi(color.SPI_ORDER_GBR)
	assert.ElementsMatch(t, b, []uint8{cG, cB, cR})
	b = c.ToSpi(color.SPI_ORDER_BRG)
	assert.ElementsMatch(t, b, []uint8{cB, cR, cG})
	b = c.ToSpi(color.SPI_ORDER_BGR)
	assert.ElementsMatch(t, b, []uint8{cB, cG, cR})

	// 	assert.ElementsMatch(t, b, []uint8{uint8(36), uint8(146), uint8(73), uint8(36), uint8(146), uint8(73), uint8(36), uint8(146), uint8(73)})
	// 	c = color.RGB{R: 255, G: 0, B: 0}
	// 	b = c.ToSpi()
	// 	assert.ElementsMatch(t, b, []uint8{uint8(109), uint8(182), uint8(219), uint8(36), uint8(146), uint8(73), uint8(36), uint8(146), uint8(73)})
}

func TestRgbAdd(t *testing.T) {
	c1 := color.RGB{R: 10, G: 20, B: 30}
	c2 := color.RGB{R: 40, G: 50, B: 255}
	expected := color.RGB{R: 50, G: 70, B: 255}
	c1.Add(c2)
	assert.Equal(t, c1, expected)
}

func TestRgbSub(t *testing.T) {
	c1 := color.RGB{R: 15, G: 20, B: 30}
	c2 := color.RGB{R: 5, G: 50, B: 15}
	expected := color.RGB{R: 10, G: 0, B: 15}
	c1.Sub(c2)
	assert.Equal(t, c1, expected)
}
