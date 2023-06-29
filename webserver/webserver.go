//go:build !tinygo
// +build !tinygo

package webserver

import (
	"ledean/display"
	"ledean/driver/button"
	"ledean/mode"
	"ledean/websocket"
	"net/http"
	"os"
	"strconv"

	"ledean/log"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Start(addr string, port int, path2Frontend string, display *display.Display, modeController *mode.ModeController, button *button.Button) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowedOrigins: []string{"http://127.0.0.1*", "http://localhost*"},
		// AllowedMethods: []string{"GET", "PUT", "DELETE"},
		// Debug: true,
	})
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/leds", MakeGetLedHandler(display)).Methods("GET")
	router.HandleFunc("/leds/{parameter}", MakeLedHandler(display)).Methods("GET")

	router.HandleFunc("/press_single", MakePressSingleHandler(button))
	router.HandleFunc("/press_double", MakePressDoubleHandler(button))
	router.HandleFunc("/press_long", MakePressLongHandler(button))

	router.HandleFunc("/exit", MakeExitHandler())

	if modeController == nil {
		router.HandleFunc("/mode/", MakeModePictureHandler())
	} else {
		router.HandleFunc("/mode/", MakeModeGetHandler(modeController))
		router.HandleFunc("/mode/{mode}", MakeModeHandler(modeController))

		for _, mode := range modeController.GetModes() {
			log.Debug("URL registering mode: " + mode.GetName())
			router.HandleFunc("/"+mode.GetName(), MakeGetModeParameterHandler(mode)).Methods("GET")
			router.HandleFunc("/"+mode.GetName()+"/limits", MakeGetModeLimitsHandler(mode)).Methods("GET")
			router.HandleFunc("/"+mode.GetName(), MakePostModeParameterHandler(modeController, mode)).Methods("POST")
		}
	}

	hub := websocket.NewHub(display, modeController, button)
	go hub.Run()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(w, r)
	})

	//Needs to be last registration
	if path2Frontend != "" {
		router.PathPrefix("/").Handler(http.FileServer(http.Dir(path2Frontend)))
	}

	log.Fatal(http.ListenAndServe(addr+":"+strconv.Itoa(int(port)), c.Handler(router)))
	// log.Fatal(http.ListenAndServe(addr+":"+strconv.Itoa(int(port)), router))
}

func MakeModePictureHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("-1"))
	}
}

func MakeExitHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Exit api was called. Shutting down LEDean")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte{})
		os.Exit(0)
	}
}
