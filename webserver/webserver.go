package webserver

import (
	"LEDean/led"
	"LEDean/led/mode"
	"LEDean/pi/button"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Start(addr string, port int, path2Frontend string, ledController *led.LedController, piButton *button.PiButton) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowedOrigins: []string{"http://127.0.0.1*", "http://localhost*"},
		// AllowedMethods: []string{"GET", "PUT", "DELETE"},
		// Debug: true,
	})
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/leds", MakeGetLedHandler(ledController)).Methods("GET")
	router.HandleFunc("/leds/{parameter}", MakeLedHandler(ledController)).Methods("GET")
	router.HandleFunc("/leds_rows", MakeGetLedHandler(ledController)).Methods("GET")
	router.HandleFunc("/press_single", MakePressSingleHandler(ledController, piButton))
	router.HandleFunc("/press_double", MakePressDoubleHandler(ledController, piButton))
	router.HandleFunc("/press_long", MakePressLongHandler(ledController, piButton))
	router.HandleFunc("/mode/", MakeModeGetHandler(ledController))
	router.HandleFunc("/mode/{mode}", MakeModeHandler(ledController))

	router.HandleFunc("/"+(mode.ModeSolid).GetFriendlyName(mode.ModeSolid{}), MakeGetModeSolidHandler(ledController)).Methods("GET")
	router.HandleFunc("/"+(mode.ModeSolid).GetFriendlyName(mode.ModeSolid{}), MakeModeSolidHandler(ledController)).Methods("POST")
	router.HandleFunc("/"+(mode.ModeSolid).GetFriendlyName(mode.ModeSolid{})+"/limits", MakeGetModeSolidLimitsHandler(ledController)).Methods("GET")

	router.HandleFunc("/"+(mode.ModeSolidRainbow).GetFriendlyName(mode.ModeSolidRainbow{}), MakeGetModeSolidRainbowHandler(ledController)).Methods("GET")
	router.HandleFunc("/"+(mode.ModeSolidRainbow).GetFriendlyName(mode.ModeSolidRainbow{}), MakeModeSolidRainbowHandler(ledController)).Methods("POST")
	router.HandleFunc("/"+(mode.ModeSolidRainbow).GetFriendlyName(mode.ModeSolidRainbow{})+"/limits", MakeGetModeSolidRainbowLimitsHandler(ledController)).Methods("GET")

	router.HandleFunc("/"+(mode.ModeTransitionRainbow).GetFriendlyName(mode.ModeTransitionRainbow{}), MakeGetModeTransitionRainbowHandler(ledController)).Methods("GET")
	router.HandleFunc("/"+(mode.ModeTransitionRainbow).GetFriendlyName(mode.ModeTransitionRainbow{}), MakeModeTransitionRainbowHandler(ledController)).Methods("POST")
	router.HandleFunc("/"+(mode.ModeTransitionRainbow).GetFriendlyName(mode.ModeTransitionRainbow{})+"/limits", MakeGetModeTransitionRainbowLimitsHandler(ledController)).Methods("GET")

	router.HandleFunc("/"+(mode.ModeRunningLed).GetFriendlyName(mode.ModeRunningLed{}), MakeGetModeRunningLedHandler(ledController)).Methods("GET")
	router.HandleFunc("/"+(mode.ModeRunningLed).GetFriendlyName(mode.ModeRunningLed{}), MakeModeRunningLedHandler(ledController)).Methods("POST")
	router.HandleFunc("/"+(mode.ModeRunningLed).GetFriendlyName(mode.ModeRunningLed{})+"/limits", MakeGetModeRunningLedLimitsHandler(ledController)).Methods("GET")

	if path2Frontend != "" {
		router.PathPrefix("/").Handler(http.FileServer(http.Dir(path2Frontend)))
	}

	log.Fatal(http.ListenAndServe(addr+":"+strconv.Itoa(int(port)), c.Handler(router)))
	// log.Fatal(http.ListenAndServe(addr+":"+strconv.Itoa(int(port)), router))
}
