//go:build has_pins
// +build has_pins

package driver

import "periph.io/x/periph/host"

func Init() {
	host.Init()
}
