//go:build !tinygo
// +build !tinygo

package webserver

import (
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

func Start(addr string, port int, path2Frontend string, modeController *mode.ModeController, button *button.Button, hub *websocket.Hub) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowedOrigins: []string{"http://127.0.0.1*", "http://localhost*"},
		// AllowedMethods: []string{"GET", "PUT", "DELETE"},
		// Debug: true,
	})
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/exit", MakeExitHandler())

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(w, r)
	})

	//Needs to be last registration
	if path2Frontend != "" {
		router.PathPrefix("/").Handler(http.FileServer(http.Dir(path2Frontend)))
	}

	log.Fatal(http.ListenAndServe(addr+":"+strconv.Itoa(int(port)), c.Handler(router)))
}

func MakeExitHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Exit api was called. Shutting down LEDean")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte{})
		os.Exit(0)
	}
}

// HandleFileServer -
func HandleFileServer(path2FrontEnd string) { //}, settingsWebserver *schema.SettingsWebserver) {
	fs := http.FileServer(http.Dir(path2FrontEnd))
	http.Handle("/", fs)
}
