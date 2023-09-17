package domain

import (
	"fmt"

	command "github.com/lean-enjoyers/catchat/pkg/command/base"
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
func (h *Hub) AddClient(client *Client) {
	h.clients[client] = true
}

// Removes the client from the hub.
func (h *Hub) DeleteClient(client *Client) {
	delete(h.clients, client)
	close(client.send)
}

func (h *Hub) BroadcastToClient(payload []byte, targetUserID string) {
	for client := range h.clients {
		if client.userID == targetUserID {
			select {
			case client.send <- payload:
			default:
				h.DeleteClient(client)
			}
		}
	}
}

func (h *Hub) BroadcastMessage(message string) {
	payload := []byte(fmt.Sprintf("<div hx-swap-oob='beforeend:#chat_body'><p>%s</p></div>", message))
	for client := range h.clients {
		select {
		// Forward the message to all other clients.
		case client.send <- payload:

		// Failed sending to the client, terminate the client.
		default:
			h.DeleteClient(client)
		}
	}
}

func (h *Hub) HandleCommand(cmd string) {
	args := command.GetArgs(cmd)
	commandName := args.GetCommand()

	if len(commandName) > 0 {
		p := command.Commands.Get(commandName)
		if p != nil {
			p.Execute(args, h)
		} else {
			h.BroadcastMessage("Error: command not found")
		}
	}

}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.AddClient(client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.DeleteClient(client)
			}
		// Message received from client.
		case message := <-h.broadcast:
			h.HandleMessage(string(message))
		}
	}
}

func (h *Hub) HandleMessage(message string) {
	if len(message) > 0 && message[0] == '/' {
		h.HandleCommand(message[1:])
	} else {
		h.BroadcastMessage(message)
	}
}
