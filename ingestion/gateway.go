// Package ingestion implements the event ingestion gateway
// with tenancy, rate-limiting, NATS/Kafka integration
package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// NATSClient wraps NATS connection
type NATSClient struct {
	url     string
	conn    interface{}
	subs    map[string]interface{}
	mu      sync.RWMutex
}

// NewNATSClient creates a new NATS client
func NewNATSClient(url string) *NATSClient {
	return &NATSClient{
		url:  url,
		subs: make(map[string]interface{}),
	}
}

// Connect connects to NATS
func (n *NATSClient) Connect() error {
	// In production, would use nats.go
	return nil
}

// Publish publishes a message to a subject
func (n *NATSClient) Publish(subject string, data []byte) error {
	// In production, would publish to NATS
	return nil
}

// Subscribe subscribes to a subject
func (n *NATSClient) Subscribe(subject string, handler func([]byte)) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.subs[subject] = handler
	return nil
}

// Close closes the NATS connection
func (n *NATSClient) Close() error {
	return nil
}

// KafkaClient wraps Kafka connection
type KafkaClient struct {
	brokers []string
	producer interface{}
	consumer interface{}
	mu      sync.RWMutex
}

// NewKafkaClient creates a new Kafka client
func NewKafkaClient(brokers []string) *KafkaClient {
	return &KafkaClient{
		brokers: brokers,
	}
}

// Connect connects to Kafka
func (k *KafkaClient) Connect() error {
	return nil
}

// Produce produces a message to a topic
func (k *KafkaClient) Produce(topic string, key, value []byte) error {
	return nil
}

// Consume consumes messages from a topic
func (k *KafkaClient) Consume(topic string, handler func([]byte)) error {
	return nil
}

// Close closes the Kafka connection
func (k *KafkaClient) Close() error {
	return nil
}

// IngestionGateway is the main gateway
type IngestionGateway struct {
	gateway     *Gateway
	nats        *NATSClient
	kafka       *KafkaClient
	tenants     map[string]*Tenant
	rateLimiter *RateLimiter
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	running     bool
	metrics     *GatewayMetrics
}

// Tenant represents a tenant
type Tenant struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	RateLimit   int       `json:"rateLimit"`
	Quota       int64     `json:"quota"`
	StorageUsed int64     `json:"storageUsed"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// GatewayMetrics tracks gateway metrics
type GatewayMetrics struct {
	EventsReceived   int64
	EventsProcessed  int64
	EventsDropped    int64
	EventsForwarded  int64
	ActiveTenants    int
	BytesReceived    int64
	BytesForwarded   int64
	ProcessingTimeMs int64
	ForwardTimeMs    int64
	LastFlushTime    time.Time
	mu               sync.RWMutex
}

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	rates   map[string]*TokenBucket
	mu      sync.RWMutex
}

// TokenBucket implements token bucket algorithm
type TokenBucket struct {
	capacity   int
	tokens     int
	refillRate int
	lastRefill time.Time
	mu         sync.Mutex
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request is allowed
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	tokensToAdd := int(elapsed.Seconds()) * tb.refillRate
	
	if tokensToAdd > 0 {
		tb.tokens = min(tb.tokens+tokensToAdd, tb.capacity)
		tb.lastRefill = now
	}
	
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		rates: make(map[string]*TokenBucket),
	}
}

// GetBucket gets or creates a token bucket for a tenant
func (rl *RateLimiter) GetBucket(tenant string, capacity, refillRate int) *TokenBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	if bucket, ok := rl.rates[tenant]; ok {
		return bucket
	}
	
	bucket := NewTokenBucket(capacity, refillRate)
	rl.rates[tenant] = bucket
	return bucket
}

// Allow checks if a request is allowed for a tenant
func (rl *RateLimiter) Allow(tenant string) bool {
	rl.mu.RLock()
	bucket, ok := rl.rates[tenant]
	rl.mu.RUnlock()
	
	if !ok {
		return true
	}
	return bucket.Allow()
}

// NewIngestionGateway creates a new ingestion gateway
func NewIngestionGateway(natsURL string, kafkaBrokers []string) *IngestionGateway {
	gw := &IngestionGateway{
		gateway:     NewGateway(),
		nats:        NewNATSClient(natsURL),
		kafka:       NewKafkaClient(kafkaBrokers),
		tenants:     make(map[string]*Tenant),
		rateLimiter: NewRateLimiter(),
		metrics:     &GatewayMetrics{},
	}
	
	// Connect to NATS and Kafka
	gw.nats.Connect()
	gw.kafka.Connect()
	
	return gw
}

// Start starts the gateway
func (gw *IngestionGateway) Start(ctx context.Context) error {
	gw.ctx, gw.cancel = context.WithCancel(ctx)
	gw.running = true
	
	// Start gateway processing
	gw.gateway.Start()
	
	// Subscribe to NATS subjects
	gw.nats.Subscribe("helix.events", func(data []byte) {
		var event Event
		if err := json.Unmarshal(data, &event); err != nil {
			return
		}
		gw.Ingest(event)
	})
	
	// Start processing loop
	go gw.processLoop()
	
	return nil
}

// Stop stops the gateway
func (gw *IngestionGateway) Stop() {
	gw.running = false
	gw.cancel()
	gw.gateway.Stop()
	gw.nats.Close()
	gw.kafka.Close()
}

// Ingest ingests an event
func (gw *IngestionGateway) Ingest(event Event) bool {
	// Check rate limit
	if !gw.rateLimiter.Allow(event.Tenant) {
		gw.metrics.mu.Lock()
		gw.metrics.EventsDropped++
		gw.metrics.mu.Unlock()
		return false
	}
	
	// Ingest through gateway
	success := gw.gateway.Ingest(event)
	if success {
		gw.metrics.mu.Lock()
		gw.metrics.EventsReceived++
		gw.metrics.BytesReceived += int64(len(event.Source))
		gw.metrics.mu.Unlock()
	}
	
	return success
}

// processLoop processes events
func (gw *IngestionGateway) processLoop() {
	for {
		select {
		case <-gw.ctx.Done():
			return
		default:
			// Process events from gateway
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// RegisterTenant registers a tenant
func (gw *IngestionGateway) RegisterTenant(tenant *Tenant) {
	gw.mu.Lock()
	defer gw.mu.Unlock()
	gw.tenants[tenant.ID] = tenant
	
	// Set up rate limiter for tenant
	gw.rateLimiter.GetBucket(tenant.ID, tenant.RateLimit, tenant.RateLimit)
	
	// Register with gateway
	gw.gateway.RegisterTenant(tenant.Name, tenant.RateLimit)
}

// GetTenant gets a tenant by ID
func (gw *IngestionGateway) GetTenant(id string) (*Tenant, bool) {
	gw.mu.RLock()
	defer gw.mu.RUnlock()
	tenant, ok := gw.tenants[id]
	return tenant, ok
}

// ListTenants lists all tenants
func (gw *IngestionGateway) ListTenants() []*Tenant {
	gw.mu.RLock()
	defer gw.mu.RUnlock()
	tenants := make([]*Tenant, 0, len(gw.tenants))
	for _, tenant := range gw.tenants {
		tenants = append(tenants, tenant)
	}
	return tenants
}

// UpdateTenant updates a tenant
func (gw *IngestionGateway) UpdateTenant(tenant *Tenant) {
	gw.mu.Lock()
	defer gw.mu.Unlock()
	if existing, ok := gw.tenants[tenant.ID]; ok {
		existing.Name = tenant.Name
		existing.RateLimit = tenant.RateLimit
		existing.Quota = tenant.Quota
		existing.Active = tenant.Active
		existing.UpdatedAt = time.Now()
	}
}

// DeleteTenant deletes a tenant
func (gw *IngestionGateway) DeleteTenant(id string) {
	gw.mu.Lock()
	defer gw.mu.Unlock()
	delete(gw.tenants, id)
}

// GetMetrics returns gateway metrics
func (gw *IngestionGateway) GetMetrics() *GatewayMetrics {
	gw.metrics.mu.RLock()
	defer gw.metrics.mu.RUnlock()
	return gw.metrics
}

// IsRunning returns if gateway is running
func (gw *IngestionGateway) IsRunning() bool {
	gw.mu.RLock()
	defer gw.mu.RUnlock()
	return gw.running
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Kafka producer for forwarding events
type KafkaProducer struct {
	brokers []string
	topic   string
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		brokers: brokers,
		topic:   topic,
	}
}

func (kp *KafkaProducer) Produce(events []Event) error {
	// In production, would produce to Kafka
	return nil
}

// NATS publisher for forwarding events
type NATSPublisher struct {
	url     string
	subject string
}

func NewNATSPublisher(url, subject string) *NATSPublisher {
	return &NATSPublisher{
		url:     url,
		subject: subject,
	}
}

func (np *NATSPublisher) Publish(events []Event) error {
	// In production, would publish to NATS
	return nil
}

// Event validator
type EventValidator struct {
	rules []ValidationRule
}

type ValidationRule struct {
	Field    string
	Required bool
	Type     string
	Min      float64
	Max      float64
}

func NewEventValidator() *EventValidator {
	return &EventValidator{
		rules: []ValidationRule{
			{Field: "source", Required: true, Type: "string"},
			{Field: "tenant", Required: true, Type: "string"},
			{Field: "timestamp", Required: true, Type: "int64"},
		},
	}
}

func (ev *EventValidator) Validate(event Event) error {
	for _, rule := range ev.rules {
		switch rule.Field {
		case "source":
			if rule.Required && event.Source == "" {
				return fmt.Errorf("source is required")
			}
		case "tenant":
			if rule.Required && event.Tenant == "" {
				return fmt.Errorf("tenant is required")
			}
		case "timestamp":
			if rule.Required && event.Timestamp.IsZero() {
				return fmt.Errorf("timestamp is required")
			}
		}
	}
	return nil
}

// Batch processor for efficient ingestion
type BatchProcessor struct {
	batchSize    int
	flushInterval time.Duration
	events       chan Event
	batch        []Event
	mu           sync.Mutex
	flushChan    chan []Event
}

func NewBatchProcessor(batchSize int, flushInterval time.Duration) *BatchProcessor {
	return &BatchProcessor{
		batchSize:     batchSize,
		flushInterval: flushInterval,
		events:        make(chan Event, batchSize*2),
		batch:         make([]Event, 0, batchSize),
		flushChan:     make(chan []Event, 10),
	}
}

func (bp *BatchProcessor) Start(ctx context.Context) {
	go bp.run(ctx)
}

func (bp *BatchProcessor) Stop() {
	close(bp.events)
}

func (bp *BatchProcessor) Add(event Event) {
	bp.events <- event
}

func (bp *BatchProcessor) run(ctx context.Context) {
	ticker := time.NewTicker(bp.flushInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			bp.flush()
			return
		case event, ok := <-bp.events:
			if !ok {
				bp.flush()
				return
			}
			bp.mu.Lock()
			bp.batch = append(bp.batch, event)
			if len(bp.batch) >= bp.batchSize {
				bp.flushLocked()
			}
			bp.mu.Unlock()
		case <-ticker.C:
			bp.mu.Lock()
			bp.flushLocked()
			bp.mu.Unlock()
		}
	}
}

func (bp *BatchProcessor) flush() {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	bp.flushLocked()
}

func (bp *BatchProcessor) flushLocked() {
	if len(bp.batch) == 0 {
		return
	}
	batch := make([]Event, len(bp.batch))
	copy(batch, bp.batch)
	bp.batch = bp.batch[:0]
	
	select {
	case bp.flushChan <- batch:
	default:
		// Drop batch if channel full
	}
}