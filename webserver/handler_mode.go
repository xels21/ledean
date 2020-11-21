package webserver

import (
	"LEDean/led"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

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
