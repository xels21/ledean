package webserver

import (
	"LEDean/led"
	"LEDean/led/mode"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func MakeGetModeTransitionRainbowHandler(ledController *led.LedController) http.HandlerFunc {
	mode, err := ledController.GetModeRef((mode.ModeTransitionRainbow).GetFriendlyName(mode.ModeTransitionRainbow{}))
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		msg := (*mode).GetParameterJson()
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}

func MakeModeTransitionRainbowHandler(ledController *led.LedController) http.HandlerFunc {
	modeSolid, err := ledController.GetModeRef((mode.ModeTransitionRainbow).GetFriendlyName(mode.ModeTransitionRainbow{}))
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
		var modeTransitionRainbowParameter mode.ModeTransitionRainbowParameter
		err = json.Unmarshal(b, &modeTransitionRainbowParameter)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		(*modeSolid).SetParameter(modeTransitionRainbowParameter)
		ledController.Restart()

		msg := []byte{}

		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}

func MakeGetModeTransitionRainbowLimitsHandler(ledController *led.LedController) http.HandlerFunc {
	mode, err := ledController.GetModeRef((mode.ModeTransitionRainbow).GetFriendlyName(mode.ModeTransitionRainbow{}))
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		msg := (*mode).GetLimitsJson()
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}
