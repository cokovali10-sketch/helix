// Package ingestion implements the event ingestion gateway
package ingestion

import (
	"fmt"
	"sync"
	"time"
)

// Event represents an observability event
type Event struct {
	Timestamp time.Time      `json:"timestamp"`
	Source    string         `json:"source"`
	Metrics   map[string]float64 `json:"metrics"`
	Tags      map[string]string  `json:"tags"`
	Tenant    string         `json:"tenant"`
}

// Gateway manages event ingestion with tenancy and rate limiting
type Gateway struct {
	events     chan Event
	tenants    map[string]int
	rateLimits map[string]int
	mu         sync.RWMutex
	quit       chan struct{}
	wg         sync.WaitGroup
}

// NewGateway creates a new ingestion gateway
func NewGateway() *Gateway {
	gw := &Gateway{
		events:     make(chan Event, 10000),
		tenants:    make(map[string]int),
		rateLimits: make(map[string]int),
		quit:       make(chan struct{}),
	}
	gw.rateLimits["default"] = 10000
	return gw
}

// RegisterTenant registers a new tenant
func (g *Gateway) RegisterTenant(tenant string, rateLimit int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.tenants[tenant] = 0
	g.rateLimits[tenant] = rateLimit
}

// Ingest receives an event from the agent
func (g *Gateway) Ingest(event Event) bool {
	g.mu.RLock()
	tenantLimit, exists := g.rateLimits[event.Tenant]
	if !exists {
		tenantLimit = g.rateLimits["default"]
	}
	g.mu.RUnlock()

	g.mu.Lock()
	defer g.mu.Unlock()

	if g.tenants[event.Tenant] >= tenantLimit {
		return false
	}
	g.tenants[event.Tenant]++

	g.events <- event
	return true
}

// Start begins the gateway processing loop
func (g *Gateway) Start() {
	g.wg.Add(1)
	go g.processLoop()
}

// Stop gracefully shuts down the gateway
func (g *Gateway) Stop() {
	close(g.quit)
	g.wg.Wait()
	close(g.events)
}

// processLoop processes incoming events
func (g *Gateway) processLoop() {
	defer g.wg.Done()
	for {
		select {
		case <-g.quit:
			return
		case event := <-g.events:
			fmt.Printf("Ingesting event from %s (tenant: %s)\n", event.Source, event.Tenant)
		}
	}
}