package webserver

import (
	"ledean/pi/button"
	"net/http"
)

func MakePressSingleHandler(piButton *button.PiButton) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		piButton.PressSingle()
		w.Write([]byte{})
	}
}
func MakePressDoubleHandler(piButton *button.PiButton) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		piButton.PressDouble()
		w.Write([]byte{})
	}
}

func MakePressLongHandler(piButton *button.PiButton) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		piButton.PressLong()
		w.Write([]byte{})
	}
}
