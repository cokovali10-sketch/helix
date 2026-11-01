// Package ingestion provides connector
package ingestion

import (
	"sync"
)

// Connector connects to source
type Connector struct {
	name   string
	status string
	mu     sync.Mutex
}

// NewConnector creates connector
func NewConnector(name string) *Connector {
	return &Connector{name: name, status: "disconnected"}
}

// Connect connects
func (c *Connector) Connect() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status = "connected"
}

// Disconnect disconnects
func (c *Connector) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status = "disconnected"
}

// Status returns status
func (c *Connector) Status() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.status
}

// Name returns name
func (c *Connector) Name() string {
	return c.name
}

// IsConnected returns connected
func (c *Connector) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.status == "connected"
}
