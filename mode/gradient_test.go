package mode

import (
	"ledean/dbdriver"
	"ledean/display"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestModeGradient(t *testing.T) {
	ledCount := 5
	dbdriver, _ := dbdriver.NewDbDriver("../tst/temp/db")
	display := display.NewDisplay(ledCount, 1, "0", "0")
	// roundTimeMs := uint32(1000)
	// // progressDegStepSize := 360 / (float64(roundTimeMs) / 1000) * (float64(RefreshIntervalNs) / 1000 / 1000 / 1000)

	modeGradient := NewModeGradient(dbdriver, display)
	modeGradientParameter := ModeGradientParameter{Brightness: 1.0, Count: 3}
	modeGradient.setParameter(modeGradientParameter)
	for i := range modeGradient.positions {
		modeGradient.positions[i].percent = 0
	}
	if modeGradient.positions[0].percent != 0 {
		t.Fatalf(`progress set not working`)
	}
	modeGradient.calcDisplay()
	if modeGradient.positions[0].percent == 0 {
		t.Fatalf(`progress not increasing`)
	}
	if modeGradient.positions[0].hue_current != modeGradient.positions[0].hue_from {
		t.Fatalf(`first position is wrong`)
	}
	if modeGradient.ledsHSV[0].H != modeGradient.positions[0].hue_current {
		t.Fatalf(`first color is wrong`)
	}
	// modeGradient.progressDegStep = 360.0 - modeGradient.progressDegStepSize

	// modeGradient.hues[0] = 80
	// modeGradient.hues[1] = 120

	// hues := modeGradient.hues //copy by value

	// modeGradient.calcDisplay()
	// modeGradient.hues[0] = 40

	// if modeGradient.progressDegStep != 0 {
	// 	t.Fatalf(`progressDegStep is false`)
	// }
	// if modeGradient.hues[1] != hues[0] {
	// 	t.Fatalf(`led shift not working`)
	// }
	// if modeGradient.ledsHSV[0].H != 80 {
	// 	t.Fatalf(`first led has wrong color`)
	// }

	// if modeGradient.ledsHSV[ledCount-1].H != 120 {
	// 	t.Fatalf(`last led wrong`)
	// }

	// for i := 0; i < FPS; i++ {
	// 	modeGradient.calcDisplay()
	// }
	// if modeGradient.ledsHSV[0].H != 40 {
	// 	t.Fatalf(`first led has wrong color`)
	// }

	// if modeGradient.ledsHSV[ledCount-1].H != 80 {
	// 	t.Fatalf(`last led wrong`)
	// }

}
