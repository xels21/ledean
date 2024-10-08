//go:build !tinygo
// +build !tinygo

// Package websocket provides the websocket implementation to avoid polling
package websocket

import (
	"ledean/log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	clientsLock sync.Mutex

	// Incomming commaand from the clients.
	// just a forwarder to the command_handler
	cmd chan Cmd

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	initClientCbs []func(*Client)

	cmdButtonChannel     chan CmdButton
	cmdModeActionChannel chan CmdModeAction
	cmdModeChannel       chan CmdMode
	// Cmd2sMode        chan Cmd2sMode
}

func NewHub() *Hub {
	return &Hub{
		cmd:                  make(chan Cmd),
		register:             make(chan *Client),
		unregister:           make(chan *Client),
		clients:              make(map[*Client]bool),
		cmdButtonChannel:     make(chan CmdButton),
		cmdModeActionChannel: make(chan CmdModeAction),
		cmdModeChannel:       make(chan CmdMode),
		// initClientCbs: make([]func(*Client), 16),
	}
}

func (self *Hub) GetCmdButtonChannel() *chan CmdButton {
	return &self.cmdButtonChannel
}
func (self *Hub) GetCmdModeActionChannel() *chan CmdModeAction {
	return &self.cmdModeActionChannel
}
func (self *Hub) GetCmdModeChannel() *chan CmdMode {
	return &self.cmdModeChannel
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
			self.clientsLock.Lock()
			log.Info("Registered client")
			self.clients[client] = true
			self.clientInit(client)
			self.clientsLock.Unlock()
		case client := <-self.unregister:
			if _, ok := self.clients[client]; ok {
				self.clientsLock.Lock()
				log.Info("Unregistered client")
				delete(self.clients, client)
				self.clientsLock.Unlock()
				client.conn.Close()
				close(client.send)
			}

		case cmd := <-self.cmd:
			self.handleCommand(cmd)

			// case <-self.display.LedsChanged:
			// 	go self.delayedCmdLeds()
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
	// case CmdButtonId:
	// 	var cmdButton CmdButton
	// 	err := json.Unmarshal(cmd.Parameter, &cmdButton)
	// 	if err != nil {
	// 		return
	// 	}
	// 	self.CmdButtonChannel <- cmdButton
	// case CmdActionId:
	// 	var cmdAction CmdAction
	// 	err := json.Unmarshal(cmd.Parameter, &cmdAction)
	// 	if err != nil {
	// 		return
	// 	}
	// 	self.CmdActionChannel <- cmdAction
	case "":
		log.Trace("Empty message. can be ignored")
	default:
		log.Info("unknown command: ", cmd.Command)
	}
}

func (self *Hub) Boradcast(cmd Cmd) {
	self.clientsLock.Lock()
	for client := range self.clients {
		client.send <- cmd
	}
	self.clientsLock.Unlock()
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
	client := NewClient(self, conn)
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// func (self *Hub) CmdLeds() {
// 	cmd2cLedsJSON, err := json.Marshal(CmdLeds{Leds: self.display.GetLeds()})
// 	if err != nil {
// 		log.Debug("Couldn't convert err to log JSON. ", err)
// 		return
// 	}

// 	for client := range self.clients {
// 		client.send <- Cmd{Command: "leds", Parameter: cmd2cLedsJSON}
// 	}
// }
