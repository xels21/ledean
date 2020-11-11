package webserver

import (
	"LEDean/led"
	"encoding/json"
	"net/http"
)

func MakeGetLedHandler(ledController *led.LedController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		msg, err := json.Marshal(ledController.GetLeds())
		if err != nil {
			msg = []byte(err.Error())
		}
		w.Write(msg)
	}
}
