package webserver

import (
	"LEDean/led"
	"LEDean/pi/button"
	"net/http"
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
