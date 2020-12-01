package webserver

import (
	"ledean/display"
	"ledean/mode"
	"ledean/pi/button"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Start(addr string, port int, path2Frontend string, display *display.Display, modeController *mode.ModeController, piButton *button.PiButton) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowedOrigins: []string{"http://127.0.0.1*", "http://localhost*"},
		// AllowedMethods: []string{"GET", "PUT", "DELETE"},
		// Debug: true,
	})
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/leds", MakeGetLedHandler(display)).Methods("GET")
	router.HandleFunc("/leds/{parameter}", MakeLedHandler(display)).Methods("GET")

	router.HandleFunc("/press_single", MakePressSingleHandler(piButton))
	router.HandleFunc("/press_double", MakePressDoubleHandler(piButton))
	router.HandleFunc("/press_long", MakePressLongHandler(piButton))

	router.HandleFunc("/mode/", MakeModeGetHandler(modeController))
	router.HandleFunc("/mode/{mode}", MakeModeHandler(modeController))

	router.HandleFunc("/"+(mode.ModeSolid).GetFriendlyName(mode.ModeSolid{}), MakeGetModeSolidHandler(modeController)).Methods("GET")
	router.HandleFunc("/"+(mode.ModeSolid).GetFriendlyName(mode.ModeSolid{}), MakeModeSolidHandler(modeController)).Methods("POST")
	router.HandleFunc("/"+(mode.ModeSolid).GetFriendlyName(mode.ModeSolid{})+"/limits", MakeGetModeSolidLimitsHandler(modeController)).Methods("GET")

	router.HandleFunc("/"+(mode.ModeSolidRainbow).GetFriendlyName(mode.ModeSolidRainbow{}), MakeGetModeSolidRainbowHandler(modeController)).Methods("GET")
	router.HandleFunc("/"+(mode.ModeSolidRainbow).GetFriendlyName(mode.ModeSolidRainbow{}), MakeModeSolidRainbowHandler(modeController)).Methods("POST")
	router.HandleFunc("/"+(mode.ModeSolidRainbow).GetFriendlyName(mode.ModeSolidRainbow{})+"/limits", MakeGetModeSolidRainbowLimitsHandler(modeController)).Methods("GET")

	router.HandleFunc("/"+(mode.ModeTransitionRainbow).GetFriendlyName(mode.ModeTransitionRainbow{}), MakeGetModeTransitionRainbowHandler(modeController)).Methods("GET")
	router.HandleFunc("/"+(mode.ModeTransitionRainbow).GetFriendlyName(mode.ModeTransitionRainbow{}), MakeModeTransitionRainbowHandler(modeController)).Methods("POST")
	router.HandleFunc("/"+(mode.ModeTransitionRainbow).GetFriendlyName(mode.ModeTransitionRainbow{})+"/limits", MakeGetModeTransitionRainbowLimitsHandler(modeController)).Methods("GET")

	router.HandleFunc("/"+(mode.ModeRunningLed).GetFriendlyName(mode.ModeRunningLed{}), MakeGetModeRunningLedHandler(modeController)).Methods("GET")
	router.HandleFunc("/"+(mode.ModeRunningLed).GetFriendlyName(mode.ModeRunningLed{}), MakeModeRunningLedHandler(modeController)).Methods("POST")
	router.HandleFunc("/"+(mode.ModeRunningLed).GetFriendlyName(mode.ModeRunningLed{})+"/limits", MakeGetModeRunningLedLimitsHandler(modeController)).Methods("GET")

	if path2Frontend != "" {
		router.PathPrefix("/").Handler(http.FileServer(http.Dir(path2Frontend)))
	}

	log.Fatal(http.ListenAndServe(addr+":"+strconv.Itoa(int(port)), c.Handler(router)))
	// log.Fatal(http.ListenAndServe(addr+":"+strconv.Itoa(int(port)), router))
}
