package webserver

import (
	"encoding/json"
	"io/ioutil"
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

func MakeGetModeParameterHandler(pMode mode.Mode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg, err := json.Marshal(pMode.GetParameter())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}

func MakeGetModeLimitsHandler(pMode mode.Mode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg, err := json.Marshal(pMode.GetLimits())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}

func MakePostModeParameterHandler(modeController *mode.ModeController, pMode mode.Mode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		err = pMode.TrySetParameter(b)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		modeController.Restart()

		msg := []byte{}
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg)
	}
}
