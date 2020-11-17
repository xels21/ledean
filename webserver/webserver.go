package webserver

import (
	"LEDean/led"
	"LEDean/pi/button"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Start(addr string, port int64, path2Frontend string, ledController *led.LedController, piButton *button.PiButton) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowedOrigins: []string{"http://127.0.0.1*", "http://localhost*"},
		// AllowedMethods: []string{"GET", "PUT", "DELETE"},
		// Debug: true,
	})
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/leds", MakeGetLedHandler(ledController)).Methods("GET")
	router.HandleFunc("/press_single", MakePressSingleHandler(ledController, piButton))
	router.HandleFunc("/press_double", MakePressDoubleHandler(ledController, piButton))
	router.HandleFunc("/press_long", MakePressLongHandler(ledController, piButton))
	router.HandleFunc("/mode/", MakeModeGetHandler(ledController))
	router.HandleFunc("/mode/{mode}", MakeModeHandler(ledController))

	if path2Frontend != "" {
		router.PathPrefix("/").Handler(http.FileServer(http.Dir(path2Frontend)))
	}

	log.Fatal(http.ListenAndServe(addr+":"+strconv.Itoa(int(port)), c.Handler(router)))
	// log.Fatal(http.ListenAndServe(addr+":"+strconv.Itoa(int(port)), router))
}
