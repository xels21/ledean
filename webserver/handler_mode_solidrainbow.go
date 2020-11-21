package webserver

import (
	"LEDean/led"
	"LEDean/led/mode"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func MakeGetModeSolidRainbowHandler(ledController *led.LedController) http.HandlerFunc {
	mode, err := ledController.GetModeRef((mode.ModeSolidRainbow).GetFriendlyName(mode.ModeSolidRainbow{}))
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		msg := (*mode).GetParameterJson()
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}

func MakeModeSolidRainbowHandler(ledController *led.LedController) http.HandlerFunc {
	modeSolid, err := ledController.GetModeRef((mode.ModeSolidRainbow).GetFriendlyName(mode.ModeSolidRainbow{}))
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {

		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Unmarshal
		var modeSolidRainbowParameter mode.ModeSolidRainbowParameter
		err = json.Unmarshal(b, &modeSolidRainbowParameter)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		(*modeSolid).SetParameter(modeSolidRainbowParameter)
		ledController.Restart()

		msg := []byte{}

		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}

func MakeGetModeSolidRainbowLimitsHandler(ledController *led.LedController) http.HandlerFunc {
	mode, err := ledController.GetModeRef((mode.ModeSolidRainbow).GetFriendlyName(mode.ModeSolidRainbow{}))
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		msg := (*mode).GetLimitsJson()
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}
