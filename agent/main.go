// Package agent collects and transmits observability events
package agent

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Event represents a single observability event
type Event struct {
	Timestamp time.Time      `json:"timestamp"`
	Source    string         `json:"source"`
	Metrics   map[string]float64 `json:"metrics"`
	Tags      map[string]string  `json:"tags"`
	Tenant    string         `json:"tenant"`
}

// Collector handles event collection with backpressure
type Collector struct {
	events chan Event
	quit   chan struct{}
	wg     sync.WaitGroup
	maxRate int
}

// NewCollector creates a new event collector
func NewCollector(maxRate int) *Collector {
	return &Collector{
		events:  make(chan Event, 10000),
		quit:    make(chan struct{}),
		maxRate: maxRate,
	}
}

// Start begins collecting events
func (c *Collector) Start(ctx context.Context) {
	c.wg.Add(1)
	go c.processLoop(ctx)
}

// Stop gracefully shuts down the collector
func (c *Collector) Stop() {
	close(c.quit)
	c.wg.Wait()
	close(c.events)
}

// processLoop continuously processes events
func (c *Collector) processLoop(ctx context.Context) {
	defer c.wg.Done()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.quit:
			return
		case <-ticker.C:
		case event, ok := <-c.events:
			if !ok {
				return
			}
			fmt.Printf("Processing event from %s\n", event.Source)
		}
	}
}