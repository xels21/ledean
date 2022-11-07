//go:build !tinygo
// +build !tinygo

package webserver

import (
	"ledean/driver/button"
	"net/http"
)

func MakePressSingleHandler(button *button.Button) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		button.PressSingle()
		w.Write([]byte{})
	}
}
func MakePressDoubleHandler(button *button.Button) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		button.PressDouble()
		w.Write([]byte{})
	}
}

func MakePressLongHandler(button *button.Button) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		button.PressLong()
		w.Write([]byte{})
	}
}
