package webserver

import (
	"encoding/json"
	"ledean/display"
	"net/http"

	"github.com/gorilla/mux"
)

func MakeGetLedHandler(display *display.Display) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(display.GetLedsJson())
	}
}

func MakeLedHandler(display *display.Display) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parameter := mux.Vars(r)["parameter"]
		msg := []byte{}
		var err error
		if parameter == "count" {
			msg, err = json.Marshal(display.GetLedCount())
			if err != nil {
				msg = []byte{}
			}
		} else if parameter == "rows" {
			msg, err = json.Marshal(display.GetLedRows())
			if err != nil {
				msg = []byte{}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}
