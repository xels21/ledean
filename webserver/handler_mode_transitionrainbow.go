package webserver

import (
	"encoding/json"
	"io/ioutil"
	"ledean/mode"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func MakeGetModeTransitionRainbowHandler(modeController *mode.ModeController) http.HandlerFunc {
	mode, err := modeController.GetModeRef((mode.ModeTransitionRainbow).GetFriendlyName(mode.ModeTransitionRainbow{}))
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		msg := (*mode).GetParameterJson()
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}

func MakeModeTransitionRainbowHandler(modeController *mode.ModeController) http.HandlerFunc {
	modeSolid, err := modeController.GetModeRef((mode.ModeTransitionRainbow).GetFriendlyName(mode.ModeTransitionRainbow{}))
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
		modeController.Restart()

		msg := []byte{}

		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}

func MakeGetModeTransitionRainbowLimitsHandler(modeController *mode.ModeController) http.HandlerFunc {
	mode, err := modeController.GetModeRef((mode.ModeTransitionRainbow).GetFriendlyName(mode.ModeTransitionRainbow{}))
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		msg := (*mode).GetLimitsJson()
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}
