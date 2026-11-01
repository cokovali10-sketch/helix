// Package api provides websocket hub
package api

import (
	"encoding/json"
	"sync"
)

// Hub manages websocket clients
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

// Client represents a websocket client
type Client struct {
	hub  *Hub
	send chan []byte
	id   string
}

// NewHub creates hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Register client
func (h *Hub) Register(c *Client) {
	h.register <- c
}

// Unregister client
func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
}

// Broadcast sends message
func (h *Hub) Broadcast(msg interface{}) {
	data, _ := json.Marshal(msg)
	h.broadcast <- data
}

// Run starts hub
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = true
			h.mu.Unlock()
		case c := <-h.unregister:
			h.mu.Lock()
			delete(h.clients, c)
			close(c.send)
			h.mu.Unlock()
		case msg := <-h.broadcast:
			h.mu.Lock()
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					close(c.send)
					delete(h.clients, c)
				}
			}
			h.mu.Unlock()
		}
	}
}

// ClientCount returns count
func (h *Hub) ClientCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.clients)
}
