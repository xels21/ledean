package webserver

import (
	"LEDean/led"
	"LEDean/led/mode"
	"LEDean/pi/button"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

func MakeGetLedHandler(ledController *led.LedController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(ledController.GetLedsJson())
	}
}

func MakePressSingleHandler(ledController *led.LedController, piButton *button.PiButton) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		piButton.PressSingle()
		w.Write([]byte{})
	}
}
func MakePressDoubleHandler(ledController *led.LedController, piButton *button.PiButton) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		piButton.PressDouble()
		w.Write([]byte{})
	}
}

func MakePressLongHandler(ledController *led.LedController, piButton *button.PiButton) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		piButton.PressLong()
		w.Write([]byte{})
	}
}

func MakeModeGetHandler(ledController *led.LedController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg, err := json.Marshal(ledController.GetModeIndex())
		if err != nil {
			msg = []byte{}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}

func MakeModeHandler(ledController *led.LedController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		modeStr := mux.Vars(r)["mode"]
		msg := []byte{}
		var err error
		if modeStr == "resolver" {
			msg, err = json.Marshal(ledController.GetModeResolver())
			if err != nil {
				msg = []byte{}
			}
		} else if modeStr == "randomize" {
			ledController.Randomize()
		} else {
			mode, err := strconv.Atoi(modeStr)
			if err != nil || mode < 0 || mode > int(ledController.GetModeLength()) {
				log.Info("Wrong mode: " + modeStr)
			} else {
				ledController.SwitchModeIndex(uint8(mode))

				msg, err = json.Marshal(ledController.GetModeIndex())
				if err != nil {
					msg = []byte{}
				}
			}
		}
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
