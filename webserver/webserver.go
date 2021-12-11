package webserver

import (
	"ledean/display"
	"ledean/mode"
	"ledean/pi/button"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
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

	for _, mode := range modeController.GetModes() {
		log.Debug("URL registering mode: " + mode.GetName())
		router.HandleFunc("/"+mode.GetName(), MakeGetModeParameterHandler(mode)).Methods("GET")
		router.HandleFunc("/"+mode.GetName()+"/limits", MakeGetModeLimitsHandler(mode)).Methods("GET")
		router.HandleFunc("/"+mode.GetName(), MakePostModeParameterHandler(modeController, mode)).Methods("POST")
	}

	if path2Frontend != "" {
		router.PathPrefix("/").Handler(http.FileServer(http.Dir(path2Frontend)))
	}

	log.Fatal(http.ListenAndServe(addr+":"+strconv.Itoa(int(port)), c.Handler(router)))
	// log.Fatal(http.ListenAndServe(addr+":"+strconv.Itoa(int(port)), router))
}
