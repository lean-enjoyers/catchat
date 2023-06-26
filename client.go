package main

import (
	"bytes"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
)

var (
	newLineByte = []byte{'\n'}
	spaceByte   = []byte{' '}
)

type Client struct {
	// Client's web socket connection.
	conn *websocket.Conn

	// Buffer for messages to be delivered to the client.
	send chan []byte

	// The hub the client belongs in.
	hub *Hub
}

func makeClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}
}
func (c *Client) connect() {
	c.hub.RegisterClient(c)
}

// Unregister self from the hub and close websocket connection.
func (c *Client) disconnect() {
	c.hub.UnregisterClient(c)
	c.conn.Close()
}

// client messages -> hub
func (c *Client) sendLoop() {
	defer c.disconnect()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newLineByte, spaceByte, -1))
		c.hub.broadcast <- message
	}
}

// hub messages -> client
func (c *Client) receiveLoop() {
	defer c.disconnect()

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
				w.Write(newLineByte)
				w.Write(<-c.send)
			}

			// Flush message to the network.
			if err := w.Close(); err != nil {
				return
			}

		}
	}
}

type Hub struct {
	// Registered Clients.
	clients map[*Client]bool

	// Register channel
	register chan *Client

	// Unregister channel
	unregister chan *Client

	// Messages to send to all clients.
	broadcast chan []byte
}

// Creates a new empty hub
func makeHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

// Registers the client by sending the client into the register channel.
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// Unregisters the client by sending the client into the unregister channel
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// Adds a client to the hub.
func (h *Hub) addClient(client *Client) {
	h.clients[client] = true
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)

				// Close the client's send channel since it's no longer in use.
				close(client.send)
			}
		// Message received from some client.
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				// Forward the message to all other clients.
				case client.send <- message:

				// Failed sending to the client, terminate the client.
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
		}

	}
}
