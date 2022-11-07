//go:build !tinygo
// +build !tinygo

package webserver

import (
	"net/http"
)

// HandleFileServer -
func HandleFileServer(path2FrontEnd string) { //}, settingsWebserver *schema.SettingsWebserver) {
	fs := http.FileServer(http.Dir(path2FrontEnd))
	http.Handle("/", fs)
}
