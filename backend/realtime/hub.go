package realtime

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

type Hub struct {
	Clients map[*Client]bool
	Mutex   sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[*Client]bool),
	}
}

func (h *Hub) AddClient(client *Client) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	h.Clients[client] = true
}

func (h *Hub) RemoveClient(client *Client) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	delete(h.Clients, client)
	close(client.Send)
}

func (h *Hub) Broadcast(message []byte) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	for client := range h.Clients {
		select {
		case client.Send <- message:
		default:
			// Skip clients whose send channel is full.
		}
	}
}