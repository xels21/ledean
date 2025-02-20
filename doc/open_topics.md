# Open Topics

## [CLOSED]no output with ESP32 (Chip: ESP32-D0WD-V3 (revision v3.0))

Something is wrong with the protocol timing for WS28x (NZP).
I guess the new generaton of ESP32 uC has higher frequency, but tinygo lib provides lower.

So the generated ASM code is too fast (should be ~300ms, is ~250ns)

I fixed it quick and dirty by hard overwriting the template:
`pkg/mod/tinygo.org/x/drivers@v0.26.0/ws2812/ws2812_xtensa.go`

```go
	switch machine.CPUFrequency() {
	case 160e6: // 160MHz
		device.AsmFull(`
		1: // send_bit
			s32i  {maskSet}, {portSet}, 0     // [1]  T0H and T1H start here
			nop                               // [37]
```
here I added 6 more `nop`'s (extended 37 to 43)

### solution
Update delivered in a fork
`https://github.com/xels21/tinygo-drivers`
and adjust in gomod
`replace tinygo.org/x/drivers => github.com/xels21/tinygo-drivers v0.0.0-20231211215924-957e975a7b22`
write `latest` and go mod tidy will resolve the version
`replace tinygo.org/x/drivers => github.com/xels21/tinygo-drivers latest`


latest working go version with tinygo `0.29` is `1.22.12`

some defines are not working, so i just commented them in the sourcecode
`C:\Users\00dea\scoop\apps\tinygo\current\src\machine\machine_esp32_i2c.go:38:16: undefined: SCL_PIN`