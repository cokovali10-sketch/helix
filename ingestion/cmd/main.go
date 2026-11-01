// Package ingestion implements the event ingestion gateway
// with tenancy, rate-limiting, and NATS/Kafka integration
package main

import (
	"fmt"
	"sync"
	"time"
)

// Gateway manages event ingestion with tenancy and rate limiting
type Gateway struct {
	events     chan Event
	tenants    map[string]int // tenant -> event count
	rateLimits map[string]int // tenant -> max events/second
	mu         sync.RWMutex
	quit       chan struct{}
	wg         sync.WaitGroup
}

// Event represents an observability event
type Event struct {
	Timestamp time.Time
	Source    string
	Metrics   map[string]float64
	Tags      map[string]string
	Tenant    string
}

// NewGateway creates a new ingestion gateway
func NewGateway() *Gateway {
	gw := &Gateway{
		events:     make(chan Event, 10000),
		tenants:    make(map[string]int),
		rateLimits: make(map[string]int),
		quit:       make(chan struct{}),
	}
	
	// Default rate limits
	gw.rateLimits["default"] = 10000 // 10k events/sec default
	
	return gw
}

// RegisterTenant registers a new tenant with its rate limit
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
		return false // Rate limit exceeded
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
			// Forward to storage engine
			fmt.Printf("Ingesting event from %s (tenant: %s)\n", event.Source, event.Tenant)
			// TODO: forward to storage engine
		}
	}
}

func main() {
	gateway := NewGateway()
	gateway.Start()
	
	// Register tenants
	gateway.RegisterTenant("soc-team-a", 5000)
	gateway.RegisterTenant("soc-team-b", 3000)
	
	// Simulate ingesting events
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			event := Event{
				Timestamp: time.Now(),
				Source:    fmt.Sprintf("agent-%d", i),
				Metrics:   map[string]float64{"temp": 22.5},
				Tags:      map[string]string{"host": "soc-01"},
				Tenant:    "soc-team-a",
			}
			success := gateway.Ingest(event)
			if success {
				fmt.Printf("Event %d ingested successfully\n", i)
			} else {
				fmt.Printf("Event %d dropped (rate limit)\n", i)
			}
		}(i)
	}
	wg.Wait()
	
	gateway.Stop()
	fmt.Println("Ingestion gateway complete")
}