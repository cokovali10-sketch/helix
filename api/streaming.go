// Package api provides additional API handlers
// including WebSocket, SSE, and streaming handlers
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	connections map[string]*WebSocketConnection
	mu          sync.RWMutex
	broadcast   chan []byte
}

// WebSocketConnection represents a WebSocket connection
type WebSocketConnection struct {
	ID       string
	UserID   string
	TenantID string
	Send     chan []byte
	Receive  chan []byte
	Closed   bool
	mu       sync.Mutex
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		connections: make(map[string]*WebSocketConnection),
		broadcast:   make(chan []byte, 1000),
	}
}

// HandleConnection handles a new WebSocket connection
func (wh *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	// In production, would upgrade to WebSocket
	conn := &WebSocketConnection{
		ID:       fmt.Sprintf("ws-%d", time.Now().UnixNano()),
		UserID:   r.Header.Get("X-User-ID"),
		TenantID: r.Header.Get("X-Tenant-ID"),
		Send:     make(chan []byte, 256),
		Receive:  make(chan []byte, 256),
	}

	wh.mu.Lock()
	wh.connections[conn.ID] = conn
	wh.mu.Unlock()

	// Handle messages
	go wh.handleMessages(conn)

	// Send welcome message
	welcome := map[string]interface{}{
		"type":    "welcome",
		"connId":  conn.ID,
		"message": "Connected to Helix WebSocket",
	}
	data, _ := json.Marshal(welcome)
	conn.Send <- data

	// Keep connection open
	select {
	case <-r.Context().Done():
		wh.closeConnection(conn.ID)
	}
}

// handleMessages handles incoming messages
func (wh *WebSocketHandler) handleMessages(conn *WebSocketConnection) {
	for {
		select {
		case msg := <-conn.Receive:
			wh.processMessage(conn, msg)
		case <-time.After(30 * time.Second):
			// Send ping
			ping := map[string]string{"type": "ping"}
			data, _ := json.Marshal(ping)
			conn.Send <- data
		}
	}
}

// processMessage processes an incoming message
func (wh *WebSocketHandler) processMessage(conn *WebSocketConnection, msg []byte) {
	var message map[string]interface{}
	if err := json.Unmarshal(msg, &message); err != nil {
		return
	}

	msgType, _ := message["type"].(string)

	switch msgType {
	case "subscribe":
		// Subscribe to a channel
		channel, _ := message["channel"].(string)
		wh.subscribe(conn, channel)
	case "unsubscribe":
		// Unsubscribe from a channel
		channel, _ := message["channel"].(string)
		wh.unsubscribe(conn, channel)
	case "query":
		// Execute a query
		query, _ := message["query"].(string)
		wh.executeQuery(conn, query)
	case "ping":
		// Respond to ping
		pong := map[string]string{"type": "pong"}
		data, _ := json.Marshal(pong)
		conn.Send <- data
	}
}

// subscribe subscribes a connection to a channel
func (wh *WebSocketHandler) subscribe(conn *WebSocketConnection, channel string) {
	// In production, would manage subscriptions
	response := map[string]interface{}{
		"type":    "subscribed",
		"channel": channel,
	}
	data, _ := json.Marshal(response)
	conn.Send <- data
}

// unsubscribe unsubscribes a connection from a channel
func (wh *WebSocketHandler) unsubscribe(conn *WebSocketConnection, channel string) {
	response := map[string]interface{}{
		"type":    "unsubscribed",
		"channel": channel,
	}
	data, _ := json.Marshal(response)
	conn.Send <- data
}

// executeQuery executes a query and sends results
func (wh *WebSocketHandler) executeQuery(conn *WebSocketConnection, query string) {
	// In production, would execute query
	response := map[string]interface{}{
		"type":    "query_result",
		"query":   query,
		"results": []interface{}{},
	}
	data, _ := json.Marshal(response)
	conn.Send <- data
}

// closeConnection closes a connection
func (wh *WebSocketHandler) closeConnection(id string) {
	wh.mu.Lock()
	defer wh.mu.Unlock()

	if conn, ok := wh.connections[id]; ok {
		conn.mu.Lock()
		conn.Closed = true
		close(conn.Send)
		conn.mu.Unlock()
		delete(wh.connections, id)
	}
}

// Broadcast broadcasts a message to all connections
func (wh *WebSocketHandler) Broadcast(msg []byte) {
	wh.mu.RLock()
	defer wh.mu.RUnlock()

	for _, conn := range wh.connections {
		select {
		case conn.Send <- msg:
		default:
		}
	}
}

// BroadcastToTenant broadcasts to a specific tenant
func (wh *WebSocketHandler) BroadcastToTenant(tenantID string, msg []byte) {
	wh.mu.RLock()
	defer wh.mu.RUnlock()

	for _, conn := range wh.connections {
		if conn.TenantID == tenantID {
			select {
			case conn.Send <- msg:
			default:
			}
		}
	}
}

// GetConnectionCount returns the number of connections
func (wh *WebSocketHandler) GetConnectionCount() int {
	wh.mu.RLock()
	defer wh.mu.RUnlock()
	return len(wh.connections)
}

// GetConnections returns all connections
func (wh *WebSocketHandler) GetConnections() []*WebSocketConnection {
	wh.mu.RLock()
	defer wh.mu.RUnlock()

	conns := make([]*WebSocketConnection, 0, len(wh.connections))
	for _, conn := range wh.connections {
		conns = append(conns, conn)
	}
	return conns
}

// SSEHandler handles Server-Sent Events
type SSEHandler struct {
	clients map[string]*SSEClient
	mu      sync.RWMutex
}

// SSEClient represents an SSE client
type SSEClient struct {
	ID       string
	UserID   string
	TenantID string
	Channel  chan string
	Closed   bool
	mu       sync.Mutex
}

// NewSSEHandler creates a new SSE handler
func NewSSEHandler() *SSEHandler {
	return &SSEHandler{
		clients: make(map[string]*SSEClient),
	}
}

// HandleConnection handles a new SSE connection
func (sh *SSEHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	client := &SSEClient{
		ID:       fmt.Sprintf("sse-%d", time.Now().UnixNano()),
		UserID:   r.Header.Get("X-User-ID"),
		TenantID: r.Header.Get("X-Tenant-ID"),
		Channel:  make(chan string, 256),
	}

	sh.mu.Lock()
	sh.clients[client.ID] = client
	sh.mu.Unlock()

	defer func() {
		sh.mu.Lock()
		delete(sh.clients, client.ID)
		sh.mu.Unlock()
	}()

	// Send initial event
	fmt.Fprintf(w, "event: connected\ndata: {\"id\":\"%s\"}\n\n", client.ID)
	w.(http.Flusher).Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case msg := <-client.Channel:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			w.(http.Flusher).Flush()
		case <-time.After(30 * time.Second):
			fmt.Fprintf(w, ": heartbeat\n\n")
			w.(http.Flusher).Flush()
		}
	}
}

// Send sends a message to a client
func (sh *SSEHandler) Send(clientID string, msg string) {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	if client, ok := sh.clients[clientID]; ok {
		select {
		case client.Channel <- msg:
		default:
		}
	}
}

// Broadcast broadcasts a message to all clients
func (sh *SSEHandler) Broadcast(msg string) {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	for _, client := range sh.clients {
		select {
		case client.Channel <- msg:
		default:
		}
	}
}

// BroadcastToTenant broadcasts to a specific tenant
func (sh *SSEHandler) BroadcastToTenant(tenantID string, msg string) {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	for _, client := range sh.clients {
		if client.TenantID == tenantID {
			select {
			case client.Channel <- msg:
			default:
			}
		}
	}
}

// GetClientCount returns the number of clients
func (sh *SSEHandler) GetClientCount() int {
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	return len(sh.clients)
}

// StreamingHandler handles streaming responses
type StreamingHandler struct {
	bufferSize int
}

// NewStreamingHandler creates a new streaming handler
func NewStreamingHandler(bufferSize int) *StreamingHandler {
	return &StreamingHandler{bufferSize: bufferSize}
}

// StreamJSON streams JSON data
func (sh *StreamingHandler) StreamJSON(w http.ResponseWriter, r *http.Request, data <-chan interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")

	encoder := json.NewEncoder(w)
	flusher := w.(http.Flusher)

	for item := range data {
		if err := encoder.Encode(item); err != nil {
			return
		}
		flusher.Flush()
	}
}

// StreamCSV streams CSV data
func (sh *StreamingHandler) StreamCSV(w http.ResponseWriter, r *http.Request, data <-chan []string) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Transfer-Encoding", "chunked")

	flusher := w.(http.Flusher)

	for row := range data {
		for i, col := range row {
			if i > 0 {
				w.Write([]byte(","))
			}
			w.Write([]byte(col))
		}
		w.Write([]byte("\n"))
		flusher.Flush()
	}
}

// StreamText streams text data
func (sh *StreamingHandler) StreamText(w http.ResponseWriter, r *http.Request, data <-chan string) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Transfer-Encoding", "chunked")

	flusher := w.(http.Flusher)

	for text := range data {
		w.Write([]byte(text))
		flusher.Flush()
	}
}

// LongPollingHandler handles long polling
type LongPollingHandler struct {
	timeout time.Duration
}

// NewLongPollingHandler creates a new long polling handler
func NewLongPollingHandler(timeout time.Duration) *LongPollingHandler {
	return &LongPollingHandler{timeout: timeout}
}

// Wait waits for an event or timeout
func (lph *LongPollingHandler) Wait(w http.ResponseWriter, r *http.Request, eventChan <-chan interface{}) {
	select {
	case event := <-eventChan:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(event)
	case <-time.After(lph.timeout):
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "timeout"})
	case <-r.Context().Done():
		return
	}
}