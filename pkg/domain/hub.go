package domain

import (
	"fmt"

	"github.com/lean-enjoyers/catchat/pkg/parser"
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

func (h *Hub) broadcastMessage(message string) {
	payload := []byte(fmt.Sprintf("<div hx-swap-oob='beforeend:#chat_body'><p>%s</p></div>", message))
	for client := range h.clients {
		select {
		// Forward the message to all other clients.
		case client.send <- payload:

		// Failed sending to the client, terminate the client.
		default:
			h.deleteClient(client)
		}
	}
}

func (h *Hub) handleCommand(command string) {
	program := parser.NewParser(command).Parse()

	if program.GetCommand() == "say" {
		msg, ok := program.GetFlag("message")

		if ok {
			h.broadcastMessage(msg)
		}

		msg1, ok1 := program.GetFlag("m")

		if ok1 {
			h.broadcastMessage(msg1)
		}

		// Neither specified
		if !(ok || ok1) {
			h.broadcastMessage("Say Error: No message.")
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
		// Message received from client.
		case message := <-h.broadcast:
			h.handleMessage(string(message))
		}
	}
}

func (h *Hub) handleMessage(message string) {
	if len(message) > 0 && message[0] == '/' {
		h.handleCommand(message[1:])
	} else {
		h.broadcastMessage(message)
	}
}
