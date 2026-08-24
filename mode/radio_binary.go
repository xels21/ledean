package mode

import (
	"errors"
	"math"

	"ledean/color"
)

const radioParameterVersion = 1

func (self *ModeController) GetCurrentModeBinary() ([]byte, error) {
	mode := self.modes[self.modesIndex]
	data := []byte{radioParameterVersion, modeRadioID(mode.GetName())}
	switch p := mode.GetParameter().(type) {
	case *ModeSolidParameter:
		data = append(data, p.RGB.R, p.RGB.G, p.RGB.B)
		data = appendFloat64(data, p.Brightness)
	case *ModeSolidRainbowParameter:
		data = appendFloat64(data, p.Brightness)
		data = appendU32(data, p.RoundTimeMs)
		data = appendHSV(data, p.Hsv)
	case *ModeTransitionRainbowParameter:
		data = appendFloat64(data, p.Brightness)
		data = appendFloat64(data, p.Spectrum)
		data = appendU32(data, p.RoundTimeMs)
		if p.Reverse {
			data = append(data, 1)
		} else {
			data = append(data, 0)
		}
	case *ModeRunningLedParameter:
		data = appendFloat64(data, p.Brightness)
		data = appendFloat64(data, p.RoundTimeMs)
		data = appendFloat64(data, p.HueFrom)
		data = appendFloat64(data, p.HueTo)
		data = appendFloat64(data, p.FadePct)
		data = append(data, byte(runningLedStyleID(p.Style)))
	case *ModeEmitterParameter:
		data = append(data, p.EmitCount, byte(emitStyleID(p.EmitStyle)))
		data = appendFloat64(data, p.MinBrightness)
		data = appendFloat64(data, p.MaxBrightness)
		data = appendU32(data, p.MinEmitLifetimeMs)
		data = appendU32(data, p.MaxEmitLifetimeMs)
		data = appendFloat64(data, p.WaveSpeedFac)
		data = appendFloat64(data, p.WaveWidthFac)
	case *ModeGradientParameter:
		data = appendFloat64(data, p.Brightness)
		data = appendU32(data, p.Count)
		data = appendU32(data, p.RoundTimeMs)
	case *ModeSpectrumParameter:
		data = appendFloat64(data, p.HueFrom720)
		data = appendFloat64(data, p.HueTo720)
		data = appendFloat64(data, p.Brightness)
		for _, position := range p.Positions {
			data = appendFloat64(data, position.FacFrom)
			data = appendFloat64(data, position.FacTo)
			data = appendU32(data, position.FacRoundTimeMs)
			data = appendFloat64(data, position.OffFrom)
			data = appendFloat64(data, position.OffTo)
			data = appendU32(data, position.OffRoundTimeMs)
		}
	case *ModePictureParameter:
		data = appendU32(data, p.PictureColumnUs)
		data = appendU32(data, p.PictureChangeIntervallMs)
		data = appendFloat64(data, p.Brightness)
		data = append(data, byte(p.PicIndex))
	default:
		return nil, errors.New("unsupported mode parameter")
	}
	return data, nil
}

func (self *ModeController) ApplyModeBinary(data []byte) error {
	if len(data) < 2 || data[0] != radioParameterVersion {
		return errors.New("invalid mode packet")
	}
	r := radioReader{data: data[2:]}
	switch data[1] {
	case 1:
		self.modeSolid.SetParameter(ModeSolidParameter{RGB: color.RGB{R: r.byte(), G: r.byte(), B: r.byte()}, Brightness: r.float64()})
		self.SwitchIndexFriendlyName("ModeSolid")
	case 2:
		self.modeSolidRainbow.SetParameter(ModeSolidRainbowParameter{Brightness: r.float64(), RoundTimeMs: r.u32(), Hsv: r.hsv()})
		self.SwitchIndexFriendlyName("ModeSolidRainbow")
	case 3:
		self.modeTransitionRainbow.SetParameter(ModeTransitionRainbowParameter{Brightness: r.float64(), Spectrum: r.float64(), RoundTimeMs: r.u32(), Reverse: r.byte() != 0})
		self.SwitchIndexFriendlyName("ModeTransitionRainbow")
	case 4:
		self.modeRunningLed.SetParameter(ModeRunningLedParameter{Brightness: r.float64(), RoundTimeMs: r.float64(), HueFrom: r.float64(), HueTo: r.float64(), FadePct: r.float64(), Style: runningLedStyle(r.byte())})
		self.SwitchIndexFriendlyName("ModeRunningLed")
	case 5:
		self.modeEmitter.SetParameter(ModeEmitterParameter{EmitCount: r.byte(), EmitStyle: emitStyle(r.byte()), MinBrightness: r.float64(), MaxBrightness: r.float64(), MinEmitLifetimeMs: r.u32(), MaxEmitLifetimeMs: r.u32(), WaveSpeedFac: r.float64(), WaveWidthFac: r.float64()})
		self.SwitchIndexFriendlyName("ModeEmitter")
	case 6:
		self.modeGradient.SetParameter(ModeGradientParameter{Brightness: r.float64(), Count: r.u32(), RoundTimeMs: r.u32()})
		self.SwitchIndexFriendlyName("ModeGradient")
	case 7:
		p := ModeSpectrumParameter{HueFrom720: r.float64(), HueTo720: r.float64(), Brightness: r.float64()}
		for i := range p.Positions {
			p.Positions[i] = ModeSpectrumParameterPosition{FacFrom: r.float64(), FacTo: r.float64(), FacRoundTimeMs: r.u32(), OffFrom: r.float64(), OffTo: r.float64(), OffRoundTimeMs: r.u32()}
		}
		self.modeSpectrum.SetParameter(p)
		self.SwitchIndexFriendlyName("ModeSpectrum")
	case 8:
		self.modePicture.SetParameter(ModePictureParameter{PictureColumnUs: r.u32(), PictureChangeIntervallMs: r.u32(), Brightness: r.float64(), PicIndex: int8(r.byte())})
		self.SwitchIndexFriendlyName("ModePicture")
	default:
		return errors.New("unknown mode")
	}
	if r.err != nil || r.pos != len(r.data) {
		return errors.New("invalid mode parameter")
	}
	return nil
}

func modeRadioID(name string) byte {
	switch name {
	case "ModeSolid":
		return 1
	case "ModeSolidRainbow":
		return 2
	case "ModeTransitionRainbow":
		return 3
	case "ModeRunningLed":
		return 4
	case "ModeEmitter":
		return 5
	case "ModeGradient":
		return 6
	case "ModeSpectrum":
		return 7
	case "ModePicture":
		return 8
	}
	return 0
}
func runningLedStyleID(style RunningLedStyle) byte {
	if style == RunningLedStyleTrigonometric {
		return 1
	}
	return 0
}
func runningLedStyle(value byte) RunningLedStyle {
	if value == 1 {
		return RunningLedStyleTrigonometric
	}
	return RunningLedStyleLinear
}
func emitStyleID(style EmitStyle) byte {
	if style == EmitStyleDrop {
		return 1
	}
	return 0
}
func emitStyle(value byte) EmitStyle {
	if value == 1 {
		return EmitStyleDrop
	}
	return EmitStylePulse
}
func appendU32(data []byte, value uint32) []byte {
	return append(data, byte(value), byte(value>>8), byte(value>>16), byte(value>>24))
}
func appendFloat64(data []byte, value float64) []byte {
	bits := math.Float64bits(value)
	for i := 0; i < 8; i++ {
		data = append(data, byte(bits>>(8*i)))
	}
	return data
}
func appendHSV(data []byte, value color.HSV) []byte {
	data = appendFloat64(data, value.H)
	data = appendFloat64(data, value.S)
	return appendFloat64(data, value.V)
}

type radioReader struct {
	data []byte
	pos  int
	err  error
}

func (r *radioReader) take(size int) []byte {
	if r.err != nil || r.pos+size > len(r.data) {
		r.err = errors.New("short mode parameter")
		return nil
	}
	value := r.data[r.pos : r.pos+size]
	r.pos += size
	return value
}
func (r *radioReader) byte() byte {
	value := r.take(1)
	if value == nil {
		return 0
	}
	return value[0]
}
func (r *radioReader) u32() uint32 {
	value := r.take(4)
	if value == nil {
		return 0
	}
	return uint32(value[0]) | uint32(value[1])<<8 | uint32(value[2])<<16 | uint32(value[3])<<24
}
func (r *radioReader) float64() float64 {
	value := r.take(8)
	if value == nil {
		return 0
	}
	var bits uint64
	for i := 0; i < 8; i++ {
		bits |= uint64(value[i]) << (8 * i)
	}
	return math.Float64frombits(bits)
}
func (r *radioReader) hsv() color.HSV {
	return color.HSV{H: r.float64(), S: r.float64(), V: r.float64()}
}
