package domain

import (
	"fmt"
)

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
func MakeHub() *Hub {
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

// Removes the client from the hub.
func (h *Hub) deleteClient(client *Client) {
	delete(h.clients, client)
	close(client.send)
}

func (h *Hub) broadcastToClient(payload []byte, targetUserID string) {
	for client := range h.clients {
		if client.userID == targetUserID {
			select {
			case client.send <- payload:
			default:
				h.deleteClient(client)
			}
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.deleteClient(client)
			}
		// Message received from some client.
		case message := <-h.broadcast:
			message = []byte(fmt.Sprintf("<div hx-swap-oob='beforeend:#chat_body'><p>%s</p></div>", message))
			for client := range h.clients {
				select {
				// Forward the message to all other clients.
				case client.send <- []byte(message):

				// Failed sending to the client, terminate the client.
				default:
					h.deleteClient(client)
				}
			}
		}
	}
}
