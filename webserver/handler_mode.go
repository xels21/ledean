package webserver

import (
	"encoding/json"
	"ledean/mode"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func MakeModeGetHandler(modeController *mode.ModeController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg, err := json.Marshal(modeController.GetIndex())
		if err != nil {
			msg = []byte{}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}

func MakeModeHandler(modeController *mode.ModeController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		modeStr := mux.Vars(r)["mode"]
		msg := []byte{}
		var err error
		if modeStr == "resolver" {
			msg, err = json.Marshal(modeController.GetModeResolver())
			if err != nil {
				msg = []byte{}
			}
		} else if modeStr == "randomize" {
			modeController.Randomize()
		} else {
			mode, err := strconv.Atoi(modeStr)
			if err != nil || mode < 0 || mode > int(modeController.GetLength()) {
				log.Info("Wrong mode: " + modeStr)
			} else {
				modeController.SwitchIndex(uint8(mode))

				msg, err = json.Marshal(modeController.GetIndex())
				if err != nil {
					msg = []byte{}
				}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}
