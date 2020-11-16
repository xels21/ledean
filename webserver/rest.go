package webserver

import (
	"LEDean/led"
	"LEDean/pi/button"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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

func MakeModeHandler(ledController *led.LedController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mode := mux.Vars(r)["mode"]
		msg := []byte{}
		var err error
		if mode == "resolver" {
			msg, err = json.Marshal(ledController.GetModeResolver())
			if err != nil {
				msg = []byte{}
			}
		} else {
			if mode != "" {
				//set mode
			}
			msg, err = json.Marshal(ledController.GetModeIndex())
			if err != nil {
				msg = []byte{}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}
