//go:build !tinygo || !poi
// +build !tinygo !poi

package ledean

import "ledean/mode"

func startPoiRadio(modeController *mode.ModeController) {}
