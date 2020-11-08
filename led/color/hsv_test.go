package color_test

import (
	"LEDean/led/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHsvToRgb(t *testing.T) {
	expectedRgb := color.RGB{R: 255, G: 0, B: 0}
	hsv := color.HSV{H: 0.0, S: 1.0, V: 1.0}
	rgb := hsv.ToRGB()
	assert.Equal(t, expectedRgb, rgb)

	expectedRgb = color.RGB{R: 0, G: 255, B: 0}
	hsv = color.HSV{H: 120.0, S: 1.0, V: 1.0}
	rgb = hsv.ToRGB()
	assert.Equal(t, expectedRgb, rgb)

	expectedRgb = color.RGB{R: 0, G: 0, B: 255}
	hsv = color.HSV{H: 240.0, S: 1.0, V: 1.0}
	rgb = hsv.ToRGB()
	assert.Equal(t, expectedRgb, rgb)

	expectedRgb = color.RGB{R: 190, G: 20, B: 160}
	hsv = color.HSV{H: 311.0, S: 0.89, V: 0.75}
	rgb = hsv.ToRGB()
	assert.InDelta(t, expectedRgb.R, rgb.R, 1)
	assert.InDelta(t, expectedRgb.G, rgb.G, 1)
	assert.InDelta(t, expectedRgb.B, rgb.B, 1)

}
