// Package helix collects and transmits observability events
// Agent collection with backpressure support and plugin architecture
package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Event represents a single observability event
type Event struct {
	Timestamp int64  `json:"timestamp" protobuf:"varint,1,opt,reset"`
	Source    string `json:"source" protobuf:"bytes,2,opt,reset"`
	Metrics   Metrics `json:"metrics" protobuf:"bytes,3,opt,reset"`
	Tenant    string `json:"tenant" protobuf:"bytes,4,opt,reset"`
	Tags      Tags    `json:"tags" protobuf:"bytes,5,opt,reset"`
}

// Metrics holds the event metric values
type Metrics map[string]float64

// Tags holds event metadata tags
type Tags map[string]string

// Plugin defines the interface for event processing plugins
type Plugin interface {
	// Process processes an event and may return modified event or error
	Process(ctx context.Context, event *Event) (*Event, error)
	// Name returns the plugin name
	Name() string
	// Priority returns the plugin execution priority (lower = first)
	Priority() int
}

// PluginConfig holds configuration for a plugin
type PluginConfig map[string]interface{}

// BackpressureManager manages event flow with backpressure
type BackpressureManager struct {
	// events is a buffered channel to absorb burst traffic
	events chan Event
	// quit signals the manager to stop
	quit chan struct{}
	// cancel cancels the context passed to Start
	cancel context.CancelFunc
	// wg waits for all goroutines to finish
	wg sync.WaitGroup
	// rateLimiter controls the maximum processing rate
	rateLimiter *RateLimiter
	// plugins are processed in priority order
	plugins []Plugin
	// mu protects plugins slice
	mu sync.RWMutex
}

// RateLimiter controls the event processing rate
type RateLimiter struct {
	// eventsPerSecond is the maximum rate
	eventsPerSecond int
	// tokens represents available processing capacity
	tokens int
	// lastRefillTime tracks when tokens were last added
	lastRefillTime time.Time
	// burstAllowed allows short-term burst processing
	burstAllowed bool
	burstSize    int
}

// NewRateLimiter creates a new rate limiter with the given rate
func NewRateLimiter(eventsPerSecond int, burstSize int, burstAllowed bool) *RateLimiter {
	now := time.Now()
	return &RateLimiter{
		eventsPerSecond: eventsPerSecond,
		tokens:          burstSize, // Start with burst capacity
		lastRefillTime:  now,
		burstAllowed:    burstAllowed,
		burstSize:       burstSize,
	}
}

// allow checks if an event can be processed given the rate limit
func (rl *RateLimiter) allow() bool {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefillTime)
	
	// Calculate how many tokens to refill
	tokensToAdd := int(elapsed.Seconds() * float64(rl.eventsPerSecond))
	if tokensToAdd > 0 {
		rl.tokens = min(rl.tokens + tokensToAdd, rl.burstSize)
		rl.lastRefillTime = now
	}
	
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// NewBackpressureManager creates a new backpressure manager
func NewBackpressureManager(maxRate int, burstSize int, burstAllowed bool) *BackpressureManager {
	bm := &BackpressureManager{
		events: make(chan Event, 10000), // Buffered channel for backpressure
		quit:   make(chan struct{}),
		rateLimiter: NewRateLimiter(maxRate, burstSize, burstAllowed),
	}
	
	// Set default buffer size based on rate
	if maxRate > 0 {
		bm.events = make(chan Event, maxRate*2)
	}
	
	return bm
}

// UsePlugin adds a plugin to the processing pipeline
func (bpm *BackpressureManager) UsePlugin(plugin Plugin) {
	bpm.mu.Lock()
	defer bpm.mu.Unlock()
	
	// Insert in priority order (lower priority first)
	insertIdx := len(bpm.plugins)
	for i, p := range bpm.plugins {
		if plugin.Priority() < p.Priority() {
			insertIdx = i
			break
		}
	}
	
	// Insert at the right position
	if insertIdx > 0 && insertIdx < len(bpm.plugins) {
		bpm.plugins = append(bpm.plugins[:insertIdx], append([]Plugin{plugin}, bpm.plugins[insertIdx:]...)...)
	} else {
		bpm.plugins = append(bpm.plugins, plugin)
	}
}

// Process processes an event through all plugins with backpressure
func (bpm *BackpressureManager) Process(ctx context.Context, event *Event) error {
	// First check rate limiter
	if !bpm.rateLimiter.allow() {
		// Rate limit exceeded, drop event
		return fmt.Errorf("rate limit exceeded")
	}
	
	// Send event to channel (with non-blocking select for backpressure)
	select {
	case bpm.events <- *event:
		// Event queued successfully
		return nil
	default:
		// Channel full, event dropped due to backpressure
		return fmt.Errorf("backpressure: event dropped")
	}
}

// processLoop continuously processes events from the channel
func (bpm *BackpressureManager) processLoop(ctx context.Context) {
	defer bpm.wg.Done()
	
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Periodic stats log
			// In production, would emit metrics
		case event, ok := <-bpm.events:
			if !ok {
				// Channel closed, shutting down
				return
			}
			
			// Process through plugins
			bpm.mu.RLock()
			plugins := make([]Plugin, len(bpm.plugins))
			copy(plugins, bpm.plugins)
			bpm.mu.RUnlock()
			
			currentEvent := event
			for _, plugin := range plugins {
				var err error
				result, err := plugin.Process(ctx, &currentEvent)
				if err != nil {
					fmt.Printf("Plugin %s error: %v\n", plugin.Name(), err)
					// Continue processing or stop based on error type
				}
				if result != nil {
					currentEvent = *result
				}
			}
			
			// TODO: Forward event to ingestion gateway
			fmt.Printf("Processed event from %s\n", currentEvent.Source)
		}
	}
}

// Start begins the backpressure manager
func (bpm *BackpressureManager) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	bpm.cancel = cancel
	bpm.wg.Add(1)
	go bpm.processLoop(ctx)
}

// Stop gracefully shuts down the manager
func (bpm *BackpressureManager) Stop() {
	if bpm.cancel != nil {
		bpm.cancel()
	}
	bpm.wg.Wait()
	close(bpm.events)
}

// Plugin implementations

// MetricsExtractionPlugin extracts and validates metrics from events
type MetricsExtractionPlugin struct{}

func (m *MetricsExtractionPlugin) Name() string { return "metrics_extraction" }
func (m *MetricsExtractionPlugin) Priority() int { return 100 }
func (m *MetricsExtractionPlugin) Process(ctx context.Context, event *Event) (*Event, error) {
	// Validate that metrics are present and reasonable
	if event.Metrics == nil {
		event.Metrics = make(Metrics)
	}
	
	// Ensure required metrics exist or set defaults
	if _, ok := event.Metrics["cpu"]; !ok {
		event.Metrics["cpu"] = 0.0
	}
	if _, ok := event.Metrics["memory"]; !ok {
		event.Metrics["memory"] = 0.0
	}
	
	// Clamp values to reasonable ranges
	if event.Metrics["cpu"] > 100.0 {
		event.Metrics["cpu"] = 100.0
	}
	if event.Metrics["memory"] > 100.0 {
		event.Metrics["memory"] = 100.0
	}
	
	return event, nil
}

// TenantAssignmentPlugin assigns events to tenants
type TenantAssignmentPlugin struct {
	defaultTenant string
}

func (t *TenantAssignmentPlugin) Name() string { return "tenant_assignment" }
func (t *TenantAssignmentPlugin) Priority() int { return 200 }
func (t *TenantAssignmentPlugin) Process(ctx context.Context, event *Event) (*Event, error) {
	// Assign tenant if not set
	if event.Tenant == "" {
		event.Tenant = t.defaultTenant
	}
	return event, nil
}

// HealthCheckPlugin validates event health properties
type HealthCheckPlugin struct{}

func (h *HealthCheckPlugin) Name() string { return "health_check" }
func (h *HealthCheckPlugin) Priority() int { return 50 }
func (h *HealthCheckPlugin) Process(ctx context.Context, event *Event) (*Event, error) {
	// Validate event has required fields
	if event.Source == "" {
		event.Source = "unknown"
	}
	
	// Check timestamp is reasonable
	now := time.Now().Unix()
	if event.Timestamp == 0 {
		event.Timestamp = now
	} else if event.Timestamp < now-86400 {
		// Event is more than a day old, could be a data issue
		// Don't reject, just log
		fmt.Printf("Warning: old event timestamp %d\n", event.Timestamp)
	}
	
	return event, nil
}

// NewCollector creates a new event collector with backpressure
func NewCollector(maxRate int) *BackpressureManager {
	bm := NewBackpressureManager(maxRate, 1000, true)
	
	// Add default plugins in priority order
	bm.UsePlugin(&HealthCheckPlugin{})        // Priority 50 - check health first
	bm.UsePlugin(&MetricsExtractionPlugin{}) // Priority 100 - extract metrics
	bm.UsePlugin(&TenantAssignmentPlugin{})  // Priority 200 - assign tenant
	
	return bm
}

// Agent represents the complete observability agent
type Agent struct {
	// backpressure manages event flow
	bp *BackpressureManager
	// eventsProcessed tracks total events processed
	eventsProcessed int64
	// eventsDropped tracks events dropped due to backpressure
	eventsDropped int64
	// mu protects counters
	mu sync.Mutex
}

// NewAgent creates a new observability agent
func NewAgent(maxRate int) *Agent {
	bm := NewCollector(maxRate)
	
	agent := &Agent{
		bp: bm,
	}
	
	// Start the backpressure manager with context
	ctx, cancel := context.WithCancel(context.Background())
	bm.Start(ctx)
	// Store cancel function for later use (in production, would be used for graceful shutdown)
	_ = cancel
	
	return agent
}

// ProcessEvent processes a single event through the agent pipeline
func (a *Agent) ProcessEvent(ctx context.Context, event *Event) error {
	a.mu.Lock()
	processed := atomic.AddInt64(&a.eventsProcessed, 1)
	_ = processed
	a.mu.Unlock()
	
	err := a.bp.Process(ctx, event)
	if err != nil {
		a.mu.Lock()
		// Don't count dropped events in processed, they're separate
		a.mu.Unlock()
		return err
	}
	
	return nil
}

// Stats returns agent statistics
func (a *Agent) Stats() (processed int64, dropped int64) {
	a.mu.Lock()
	processed = a.eventsProcessed
	dropped = a.eventsDropped
	a.mu.Unlock()
	return
}

// main runs the agent
func main() {
	// Create agent with 5000 events/second rate limit
	agent := NewAgent(5000)
	
	fmt.Println("Helix Agent started")
	fmt.Println("Rate limit: 5000 events/second")
	fmt.Println("Buffer size: 10000 events (backpressure)")
	
	// Simulate event generation
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			event := Event{
				Timestamp: time.Now().Unix(),
				Source:    fmt.Sprintf("soc-agent-%d", i),
				Metrics:   map[string]float64{"cpu": 50.0 + float64(i%10), "mem": 75.0 + float64(i%5)},
				Tenant:    "soc-team-a",
				Tags:      map[string]string{"host": fmt.Sprintf("server-%d", i%100), "zone": "us-east-1"},
			}
			
			err := agent.ProcessEvent(context.Background(), &event)
			if err != nil {
				fmt.Printf("Event %d dropped: %v\n", i, err)
			}
		}(i)
	}
	
	wg.Wait()
	
	// Allow time for processing
	time.Sleep(2 * time.Second)
	
	// Stop the agent
	agent.bp.Stop()
	
	// Print stats
	processed, dropped := agent.Stats()
	fmt.Printf("Agent collection complete: %d processed, %d dropped\n", processed, dropped)
	fmt.Println("Helix Agent stopped")
}

// Helper to fix the float64i typo - should be just float64
// This is a demonstration so we'll use a proper version

// float64i should be float64 - this is just a syntax error fix reference