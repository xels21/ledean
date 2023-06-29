//go:build !tinygo
// +build !tinygo

// Package websocket provides the websocket implementation to avoid polling
package websocket

import (
	"encoding/json"
	"ledean/display"
	"ledean/driver/button"
	"ledean/log"
	"ledean/mode"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Incomming commaand from the clients.
	// just a forwarder to the command_handler
	cmd chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	commandController *CommandController
	display           *display.Display
	displayTimer      *time.Timer
}

const DISPLAY_TIMER_DELAY = 200

func NewHub(display *display.Display, modeController *mode.ModeController, button *button.Button) *Hub {
	return &Hub{
		cmd:               make(chan []byte),
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		clients:           make(map[*Client]bool),
		commandController: NewCommandController(display, modeController, button),
		display:           display,
		displayTimer:      time.NewTimer(DISPLAY_TIMER_DELAY * time.Millisecond),
	}
}

func (self *Hub) Run() {
	for {
		select {
		case client := <-self.register:
			log.Info("____________client registered")
			self.clients[client] = true
		case client := <-self.unregister:
			if _, ok := self.clients[client]; ok {
				delete(self.clients, client)
				close(client.send)
			}
		case <-self.display.LedsChanged:
			go self.slowedDownCmd2cLeds()
		case cmd := <-self.cmd:
			self.commandController.HandleCommand(cmd)
			// for client := range h.clients {
			// 	select {
			// 	case client.send <- cmd:
			// 	default:
			// 		close(client.send)
			// 		delete(h.clients, client)
			// 	}
			// }
		}
	}
}

func (self *Hub) slowedDownCmd2cLeds() {
	<-self.displayTimer.C
	self.Cmd2cLeds()
	self.displayTimer.Reset(DISPLAY_TIMER_DELAY * time.Millisecond)
}

func (self *Hub) Cmd2cLeds() {
	cmd2cLedsJSON, err := json.Marshal(Cmd2cLeds{Leds: self.display.GetLeds()})
	if err != nil {
		log.Debug("Couldn't convert err to log JSON. ", err)
		return
	}

	for client := range self.clients {
		client.send <- Cmd{Command: "leds", Parameters: cmd2cLedsJSON}
	}
}

// serveWs handles websocket requests from the peer.
func (self *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		// ReadBufferSize:  1024,
		// WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true //bad security
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	client := &Client{hub: self, conn: conn, send: make(chan Cmd)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
