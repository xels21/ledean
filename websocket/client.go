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

	pCmdButtonChannel     *chan CmdButton
	pCmdModeActionChannel *chan CmdModeAction
	pCmdModeChannel       *chan CmdMode
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:                   hub,
		conn:                  conn,
		send:                  make(chan Cmd),
		pCmdButtonChannel:     hub.GetCmdButtonChannel(),
		pCmdModeActionChannel: hub.GetCmdModeActionChannel(),
		pCmdModeChannel:       hub.GetCmdModeChannel(),
	}
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
			if websocket.IsCloseError(err,
				// websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
			// websocket.CloseProtocolError,
			// websocket.CloseUnsupportedData,
			// websocket.CloseNoStatusReceived,
			// websocket.CloseAbnormalClosure,
			// websocket.CloseInvalidFramePayloadData,
			// websocket.ClosePolicyViolation,
			// websocket.CloseMessageTooBig,
			// websocket.CloseMandatoryExtension,
			// websocket.CloseInternalServerErr,
			// websocket.CloseServiceRestart,
			// websocket.CloseTryAgainLater,
			// websocket.CloseTLSHandshake,
			) || strings.Contains(err.Error(), "i/o timeout") {
				// Connection closed due to client, should be unregistered -> return
				log.Debug("Connection closed with: '", err, "'")
				return
			}

			log.Error("Error reading json message from client: ", err)
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
		*self.pCmdButtonChannel <- cmdButton
	case CmdModeActionId:
		var cmdModeAction CmdModeAction
		err := json.Unmarshal(cmd.Parameter, &cmdModeAction)
		if err != nil {
			log.Debug("Could not parse mode action parm mgs: ", string(cmd.Parameter))
			return
		}
		*self.pCmdModeActionChannel <- cmdModeAction
	case CmdModeId: //for parameter
		var cmdMode CmdMode
		err := json.Unmarshal(cmd.Parameter, &cmdMode)
		if err != nil {
			log.Debug("Could not parse mode parm mgs: ", string(cmd.Parameter))
			return
		}
		*self.pCmdModeChannel <- cmdMode
	case "":
		log.Trace("Empty message. can be ignored")
	default:
		log.Debug("Unknown command: '", cmd.Command, "' from client")
	}
}

func (c *Client) closeConnection() {
	c.hub.unregister <- c
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
				// c.closeConnection() //will be done by defer
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
