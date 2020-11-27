package webserver

import (
	"LEDean/led"
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

func MakeLedHandler(ledController *led.LedController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parameter := mux.Vars(r)["parameter"]
		msg := []byte{}
		var err error
		if parameter == "count" {
			msg, err = json.Marshal(ledController.GetLedCount())
			if err != nil {
				msg = []byte{}
			}
		} else if parameter == "rows" {
			msg, err = json.Marshal(ledController.GetLedRows())
			if err != nil {
				msg = []byte{}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}
