package mode

// import (
// 	"ledean/color"
// 	"encoding/json"
// 	"math"
// 	"math/rand"
// 	"time"

// 	"github.com/sdomino/scribble"
// 	log "github.com/sirupsen/logrus"
// )

// const (
// 	RaindropStyleDrop  RaindropStyle = "drop"
// 	RaindropStylePulse RaindropStyle = "pulse"
// )

// type RaindropStyle string

// type ModeRaindrop struct {
// 	dbDriver         *scribble.Driver
// 	leds             []color.RGB
// 	led_rows         int
// 	led_per_row      int
// 	cUpdate          *chan bool
// 	parameter        ModeRaindropParameter
// 	limits           ModeRaindropLimits
// 	cExit            chan bool
// 	ledsSingleRowRGB []color.RGB
// 	// activatedLeds     [][]float64
// 	// position          []float64
// 	// preogressStepSize []float64
// 	// hueFrom           []float64
// 	// hueTo             []float64
// }
// type Raindrop struct {
// 	Brightness float64       `json:"brightness"`
// 	LifespanMs float64       `json:"roundTimeMs"`
// 	HueFrom    float64       `json:"hueFrom"`
// 	HueTo      float64       `json:"hueTo"`
// 	Impact     float64       `json:"impact"`
// 	Style      RaindropStyle `json:"style"`
// }
// type ModeRaindropParameter struct {
// 	Brightness  float64       `json:"brightness"`
// 	RoundTimeMs float64       `json:"roundTimeMs"`
// 	HueFrom     float64       `json:"hueFrom"`
// 	HueTo       float64       `json:"hueTo"`
// 	FadePct     float64       `json:"fadePct"`
// 	Style       RaindropStyle `json:"style"`
// }
// type ModeRaindropLimits struct {
// 	MinRoundTimeMs uint32  `json:"minRoundTimeMs"`
// 	MaxRoundTimeMs uint32  `json:"maxRoundTimeMs"`
// 	MinBrightness  float64 `json:"minBrightness"`
// 	MaxBrightness  float64 `json:"maxBrightness"`
// 	MinFadePct     float64 `json:"minFadePct"`
// 	MaxFadePct     float64 `json:"maxFadePct"`
// }

// func NewModeRaindrop(dbDriver *scribble.Driver, cUpdate *chan bool, leds []color.RGB, led_rows int) *ModeRaindrop {
// 	self := ModeRaindrop{
// 		dbDriver:    dbDriver,
// 		leds:        leds,
// 		led_rows:    led_rows,
// 		led_per_row: len(leds) / led_rows,
// 		cUpdate:     cUpdate,
// 		limits: ModeRaindropLimits{
// 			MinRoundTimeMs: 1000,  //1s
// 			MaxRoundTimeMs: 30000, //30s
// 			MinBrightness:  0.3,
// 			MaxBrightness:  1.0,
// 			MinFadePct:     0.0,
// 			MaxFadePct:     1.0,
// 		},
// 		positionDeg:      0.0,
// 		ledsSingleRowRGB: make([]color.RGB, len(leds)/led_rows),
// 		activatedLeds:    make([]float64, len(leds)/led_rows),
// 		cExit:            make(chan bool, 1),
// 	}

// 	err := dbDriver.Read(self.GetFriendlyName(), "parameter", &self.parameter)
// 	if err != nil {
// 		self.Randomize()
// 	} else {
// 		self.postSetParameter()
// 	}

// 	return &self
// }

// func (ModeRaindrop) GetFriendlyName() string {
// 	return "ModeRaindrop"
// }

// func (self *ModeRaindrop) GetParameterJson() []byte {
// 	msg, _ := json.Marshal(self.parameter)
// 	return msg
// }

// func (self *ModeRaindrop) GetLimitsJson() []byte {
// 	msg, _ := json.Marshal(self.limits)
// 	return msg
// }

// func (self *ModeRaindrop) SetParameter(parm interface{}) {
// 	switch parm.(type) {
// 	case ModeRaindropParameter:
// 		self.parameter = parm.(ModeRaindropParameter)
// 		self.dbDriver.Write(self.GetFriendlyName(), "parameter", self.parameter)
// 		self.postSetParameter()
// 	}
// }

// func (self *ModeRaindrop) postSetParameter() {
// 	self.hueDistance = math.Abs(self.parameter.HueFrom - self.parameter.HueTo)
// 	self.hueDistanceFct = 1.0
// 	if self.parameter.HueFrom > self.parameter.HueTo {
// 		self.hueDistanceFct = -1.0
// 	}
// 	self.positionDegStepSize = 360.0 / (float64(self.parameter.RoundTimeMs) / 1000.0 /*s*/) * (float64(REFRESH_INTERVAL_NS) / 1000.0 /*s*/ / 1000.0 /*ms*/ / 1000.0 /*us*/)
// 	self.darkenStepSize = (1 / self.parameter.FadePct) / (float64(self.parameter.RoundTimeMs) / 1000.0 /*s*/) * (float64(REFRESH_INTERVAL_NS) / 1000.0 /*s*/ / 1000.0 /*ms*/ / 1000.0 /*us*/)
// 	self.lightenStepSize = 2 * self.parameter.Brightness * (float64(len(self.leds)) / float64(self.led_rows)) / (float64(self.parameter.RoundTimeMs) / 1000.0 /*s*/) * (float64(REFRESH_INTERVAL_NS) / 1000.0 /*s*/ / 1000.0 /*ms*/ / 1000.0 /*us*/)
// }

// func (self *ModeRaindrop) Activate() {
// 	log.Debugf("start "+self.GetFriendlyName()+" with:\n %s\n", self.GetParameterJson())
// 	ticker := time.NewTicker(REFRESH_INTERVAL_NS)

// 	go func() {
// 		for {
// 			select {
// 			case <-self.cExit:
// 				ticker.Stop()
// 				return
// 			case <-ticker.C:
// 				self.renderLoop()
// 			}
// 		}
// 	}()
// }

// func (self *ModeRaindrop) darkenLeds() {
// 	for i := 0; i < len(self.activatedLeds); i++ {
// 		if self.activatedLeds[i] != 0.0 {
// 			if self.activatedLeds[i] <= self.darkenStepSize {
// 				self.activatedLeds[i] = 0.0
// 			} else {
// 				self.activatedLeds[i] -= self.darkenStepSize
// 			}
// 		}
// 	}
// }

// func (self *ModeRaindrop) renderLoop() {

// 	self.positionDeg += self.positionDegStepSize
// 	if self.positionDeg >= 360.0 {
// 		self.positionDeg -= 360.0
// 	}

// 	self.darkenLeds()

// 	var activeLedIdx int

// 	switch self.parameter.Style {
// 	case RaindropStyleLinear:
// 		position := self.positionDeg / 180.0
// 		if position > 1.0 {
// 			position = 2.0 - position
// 		}
// 		activeLedIdx = int(position * float64(len(self.activatedLeds)))

// 		break
// 	case RaindropStyleTrigonometric:
// 		activeLedIdx = int(((math.Cos((self.positionDeg+180.0)*math.Pi/180.0) + 1.0) / 2) * float64(len(self.activatedLeds)))
// 		break
// 	}

// 	if self.activatedLeds[activeLedIdx] != 0.0 {
// 		self.activatedLeds[activeLedIdx] += self.darkenStepSize
// 	}
// 	self.activatedLeds[activeLedIdx] += self.lightenStepSize
// 	if self.activatedLeds[activeLedIdx] > 1.0 {
// 		self.activatedLeds[activeLedIdx] = 1.0
// 	}

// 	c := color.HSV{H: 0.0, S: 1.0, V: 0.0}
// 	for i, activatedLed := range self.activatedLeds {
// 		// if i == activeLedIdx {
// 		// c.H = self.parameter.HueTo
// 		// } else {
// 		c.H = self.parameter.HueFrom + (self.hueDistanceFct * self.hueDistance * activatedLed)
// 		// }
// 		c.V = activatedLed
// 		self.ledsSingleRowRGB[i] = c.ToRGB()
// 	}

// 	for r := 0; r < self.led_rows; r++ {
// 		for ri := 0; ri < len(self.leds)/self.led_rows; ri++ {
// 			i := r*self.led_per_row + ri
// 			self.leds[i] = self.ledsSingleRowRGB[ri]
// 		}
// 	}

// 	*self.cUpdate <- true
// }

// func (self *ModeRaindrop) AddRaindrop(position float64, speed float64) {

// }

// func (self *ModeRaindrop) Deactivate() {
// 	self.cExit <- true
// }

// func (self *ModeRaindrop) Randomize() {
// 	rand.Seed(time.Now().UnixNano())

// 	self.SetParameter(ModeRaindropParameter{
// 		Brightness:  rand.Float64()*(self.limits.MaxBrightness-self.limits.MinBrightness) + self.limits.MinBrightness,
// 		FadePct:     rand.Float64()*(self.limits.MaxFadePct-self.limits.MinFadePct) + self.limits.MinFadePct,
// 		HueFrom:     rand.Float64() * 360.0,
// 		HueTo:       rand.Float64() * 360.0,
// 		RoundTimeMs: rand.Float64()*float64(self.limits.MaxRoundTimeMs-self.limits.MinRoundTimeMs) + float64(self.limits.MinRoundTimeMs),
// 		Style:       self.getRandomStyle(),
// 	})
// }

// func (self *ModeRaindrop) getRandomStyle() RaindropStyle {
// 	styleSwitch := rand.Uint32() % 2
// 	var style RaindropStyle
// 	switch styleSwitch {
// 	case 0:
// 		style = RaindropStyleLinear
// 	case 1:
// 		style = RaindropStyleTrigonometric
// 	}
// 	return style
// }
