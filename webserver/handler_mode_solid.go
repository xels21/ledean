package webserver

import (
	"LEDean/led"
	"LEDean/led/mode"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func MakeGetModeSolidHandler(ledController *led.LedController) http.HandlerFunc {
	mode, err := ledController.GetModeRef((mode.ModeSolid).GetFriendlyName(mode.ModeSolid{}))
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		msg := (*mode).GetParameterJson()
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}

func MakeModeSolidHandler(ledController *led.LedController) http.HandlerFunc {
	modeSolid, err := ledController.GetModeRef((mode.ModeSolid).GetFriendlyName(mode.ModeSolid{}))
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
		var modeSolidParameter mode.ModeSolidParameter
		err = json.Unmarshal(b, &modeSolidParameter)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		(*modeSolid).SetParameter(modeSolidParameter)
		ledController.Restart()

		msg := []byte{}

		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}

func MakeGetModeSolidLimitsHandler(ledController *led.LedController) http.HandlerFunc {
	mode, err := ledController.GetModeRef((mode.ModeSolid).GetFriendlyName(mode.ModeSolid{}))
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		msg := (*mode).GetLimitsJson()
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}
