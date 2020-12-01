package webserver

import (
	"encoding/json"
	"io/ioutil"
	"ledean/mode"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func MakeGetModeSolidRainbowHandler(modeController *mode.ModeController) http.HandlerFunc {
	mode, err := modeController.GetModeRef((mode.ModeSolidRainbow).GetFriendlyName(mode.ModeSolidRainbow{}))
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		msg := (*mode).GetParameterJson()
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}

func MakeModeSolidRainbowHandler(modeController *mode.ModeController) http.HandlerFunc {
	modeSolid, err := modeController.GetModeRef((mode.ModeSolidRainbow).GetFriendlyName(mode.ModeSolidRainbow{}))
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
		modeController.Restart()

		msg := []byte{}

		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}

func MakeGetModeSolidRainbowLimitsHandler(modeController *mode.ModeController) http.HandlerFunc {
	mode, err := modeController.GetModeRef((mode.ModeSolidRainbow).GetFriendlyName(mode.ModeSolidRainbow{}))
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		msg := (*mode).GetLimitsJson()
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}
