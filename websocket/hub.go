//go:build !tinygo
// +build !tinygo

// Package websocket provides the websocket implementation to avoid polling
package websocket

import (
	"ledean/log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Incomming commaand from the clients.
	// just a forwarder to the command_handler
	cmd chan Cmd

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	initClientCbs []func(*Client)
}

func NewHub() *Hub {
	return &Hub{
		cmd:        make(chan Cmd),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		// initClientCbs: make([]func(*Client), 16),
	}
}

func (self *Hub) AppendInitClientCb(cb func(*Client)) {
	self.initClientCbs = append(self.initClientCbs, cb)
}

func (self *Hub) clientInit(client *Client) {
	for _, cb := range self.initClientCbs {
		cb(client)
	}
}

func (self *Hub) Run() {
	for {
		select {
		case client := <-self.register:
			self.clients[client] = true
			self.clientInit(client)
		case client := <-self.unregister:
			if _, ok := self.clients[client]; ok {
				delete(self.clients, client)
				close(client.send)
			}

		case cmd := <-self.cmd:
			self.handleCommand(cmd)

			// case <-self.display.LedsChanged:
			// 	go self.delayedCmd2cLeds()
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

func (self *Hub) handleCommand(cmd Cmd) {
	switch cmd.Command {
	// case "leds":
	// go self.delayedCmd2cLeds(cmd)
	default:
		log.Info("unknown command: ", cmd.Command)
	}
}

func (self *Hub) Boradcast(cmd Cmd) {
	for client := range self.clients {
		client.send <- cmd
	}
}

// func (self *Hub) Cmd2cLeds() {
// 	cmd2cLedsJSON, err := json.Marshal(Cmd2cLeds{Leds: self.display.GetLeds()})
// 	if err != nil {
// 		log.Debug("Couldn't convert err to log JSON. ", err)
// 		return
// 	}

// 	for client := range self.clients {
// 		client.send <- Cmd{Command: "leds", Parameters: cmd2cLedsJSON}
// 	}
// }

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
