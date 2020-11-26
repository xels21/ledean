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

func TestHsvAdd(t *testing.T) {
	c1 := color.HSV{H: 0.0, S: 1.0, V: 0.5}
	c2 := color.HSV{H: 0.0, S: 1.0, V: 0.2}
	expected := color.HSV{H: 0.0, S: 1.0, V: 0.7}
	c1.Add(c2)
	assert.InDelta(t, expected.H, c1.H, 1.0)
	assert.InDelta(t, expected.S, c1.S, 0.05)
	assert.InDelta(t, expected.V, c1.V, 0.05)
}

func TestHsvSub(t *testing.T) {
	c1 := color.HSV{H: 0.0, S: 1.0, V: 0.5}
	c2 := color.HSV{H: 0.0, S: 1.0, V: 0.2}
	expected := color.HSV{H: 0.0, S: 1.0, V: 0.3}
	c1.Sub(c2)
	assert.InDelta(t, expected.H, c1.H, 1.0)
	assert.InDelta(t, expected.S, c1.S, 0.05)
	assert.InDelta(t, expected.V, c1.V, 0.05)
}
