package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	id    int
	conn  websocket.Conn
	valid bool
}

type Hub struct {
	clients []Client
}

func (h *Hub) addClient(client Client) {
	h.clients = append(h.clients, client)
}

func (h *Hub) setValid(c int, b bool) {
	h.clients[c].valid = b
}

func (h *Hub) writeByteToClient(c int, s []byte) error {
	return h.clients[c].conn.WriteMessage(websocket.TextMessage, s)
}

func (h *Hub) writeStringToClient(c int, s string) error {
	return h.clients[c].conn.WriteMessage(websocket.TextMessage, []byte(s))
}
