package webserver

import (
	"LEDean/led"
	"net/http"
)

func MakeGetLedHandler(ledController *led.LedController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(ledController.GetLedsJson())
	}
}
