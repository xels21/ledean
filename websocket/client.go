//go:build !tinygo
// +build !tinygo

package websocket

import (
	"ledean/json"
	"ledean/log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 2048
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan Cmd
}

func (self *Client) SendCmd(cmd Cmd) {
	self.send <- cmd
}

// func (self *Client) SendCmd(cmd Cmd) {
// 	self.send <- cmd
// }

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.closeConnection()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		cmd := Cmd{}
		// messageType, p, err := c.conn.ReadMessage()
		// log.Info(messageType)
		// log.Info(string(p))
		// log.Info(err)
		err := c.conn.ReadJSON(&cmd)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseAbnormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseInternalServerErr,
				websocket.CloseInvalidFramePayloadData,
				websocket.CloseMandatoryExtension,
				websocket.CloseMessage,
				websocket.CloseMessageTooBig,
				websocket.CloseNoStatusReceived,
				websocket.CloseNormalClosure,
				websocket.ClosePolicyViolation,
				websocket.CloseProtocolError,
				websocket.CloseServiceRestart,
				websocket.CloseTLSHandshake,
				websocket.CloseTryAgainLater,
				websocket.CloseUnsupportedData) ||
				err.Error() == "websocket: close 1001 (going away)" ||
				strings.Index(err.Error(), "i/o timeout") >= 0 {
				log.Debug("Error: ", err)
				return
			}
			log.Debug("Error reading json message from client: ", err)
			continue
		}

		c.handleCommand(&cmd)

	}
}

func (self *Client) handleCommand(cmd *Cmd) {
	switch cmd.Command {
	case CmdButtonId:
		var cmdButton CmdButton
		err := json.Unmarshal(cmd.Parameter, &cmdButton)
		if err != nil {
			log.Debug("Could not parse button parm mgs: ", string(cmd.Parameter))
			return
		}
		self.hub.CmdButtonChannel <- cmdButton
	case CmdModeActionId:
		log.Debug("got command: ", cmd.Command)
		var cmdModeAction CmdModeAction
		err := json.Unmarshal(cmd.Parameter, &cmdModeAction)
		if err != nil {
			log.Debug("Could not parse mode action parm mgs: ", string(cmd.Parameter))
			return
		}
		self.hub.CmdModeActionChannel <- cmdModeAction

	default:
		log.Debug("Unknown command: ", cmd.Command, " from client")
	}
}

func (c *Client) closeConnection() {
	log.Info("Close connection")
	c.hub.unregister <- c
	c.conn.Close()
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.closeConnection()
	}()
	for {
		select {
		case send, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				log.Trace("The hub closed the channel.")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(send)
			if err != nil {
				c.closeConnection()
				return
				// if err.Error() == "websocket: close sent" {
				// log.Debug("clean up client connection")
				// c.closeConnection()
				// return
				// }
				// log.Info("Couldn't send: ", send, " ; ", err)
				// continue
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
