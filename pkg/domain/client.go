package domain

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/lean-enjoyers/catchat/pkg/utils"
)

type Client struct {
	// Client's web socket connection.
	conn *websocket.Conn

	// Buffer for messages to be delivered to the client.
	send chan []byte

	// The hub the client belongs in.
	hub *Hub

	// Client username
	userID string
}

const (
	writeWait = 10 * time.Second
)

func MakeClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) SelectHub(hub *Hub) {
	c.hub = hub
}

func (c *Client) Connect() {
	c.hub.RegisterClient(c)
}

func (c *Client) closeWebsocketConn() {
	c.conn.Close()
}

// Unregister self from the hub and close websocket connection.
func (c *Client) Disconnect() {
	c.hub.UnregisterClient(c)
	c.closeWebsocketConn()
}

// client messages -> hub
func (c *Client) SendLoop() {
	defer c.Disconnect()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		message = utils.TrimByte(message)
		c.hub.broadcast <- message
	}
}

// hub messages -> client
func (c *Client) ReceiveLoop() {
	defer c.Disconnect()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)

			if err != nil {
				return
			}

			w.Write(message)

			// Get queued messages and write.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(utils.NEW_LINE_BYTE)
				w.Write(<-c.send)
			}

			// Flush message to the network.
			if err := w.Close(); err != nil {
				return
			}

		}
	}
}
