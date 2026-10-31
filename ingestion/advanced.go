// Package ingestion provides additional ingestion features
// including advanced rate limiting, tenant management, and message queuing
package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// AdvancedRateLimiter implements advanced rate limiting with multiple algorithms
type AdvancedRateLimiter struct {
	algorithm    string
	buckets      map[string]*TokenBucket
	slidingWindows map[string]*SlidingWindow
	leakyBuckets map[string]*LeakyBucket
	mu           sync.RWMutex
}

// SlidingWindow implements sliding window rate limiting
type SlidingWindow struct {
	windowSize time.Duration
	maxRequests int
	requests   []time.Time
	mu         sync.Mutex
}

func NewSlidingWindow(windowSize time.Duration, maxRequests int) *SlidingWindow {
	return &SlidingWindow{
		windowSize:  windowSize,
		maxRequests: maxRequests,
		requests:    make([]time.Time, 0),
	}
}

func (sw *SlidingWindow) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-sw.windowSize)
	
	// Remove old requests
	var valid []time.Time
	for _, t := range sw.requests {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	sw.requests = valid
	
	if len(sw.requests) < sw.maxRequests {
		sw.requests = append(sw.requests, now)
		return true
	}
	return false
}

func (sw *SlidingWindow) Count() int {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return len(sw.requests)
}

// LeakyBucket implements leaky bucket rate limiting
type LeakyBucket struct {
	capacity   int
	leakRate   float64
	water      float64
	lastLeak   time.Time
	mu         sync.Mutex
}

func NewLeakyBucket(capacity int, leakRate float64) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		leakRate: leakRate,
		water:    0,
		lastLeak: time.Now(),
	}
}

func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	
	now := time.Now()
	elapsed := now.Sub(lb.lastLeak).Seconds()
	lb.water = max(0, lb.water-elapsed*lb.leakRate)
	lb.lastLeak = now
	
	if lb.water+1 <= float64(lb.capacity) {
		lb.water++
		return true
	}
	return false
}

func (lb *LeakyBucket) Level() float64 {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	return lb.water
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// NewAdvancedRateLimiter creates a new advanced rate limiter
func NewAdvancedRateLimiter(algorithm string) *AdvancedRateLimiter {
	return &AdvancedRateLimiter{
		algorithm:      algorithm,
		buckets:        make(map[string]*TokenBucket),
		slidingWindows: make(map[string]*SlidingWindow),
		leakyBuckets:   make(map[string]*LeakyBucket),
	}
}

func (arl *AdvancedRateLimiter) Allow(tenant string) bool {
	arl.mu.RLock()
	defer arl.mu.RUnlock()
	
	switch arl.algorithm {
	case "token_bucket":
		if bucket, ok := arl.buckets[tenant]; ok {
			return bucket.Allow()
		}
	case "sliding_window":
		if window, ok := arl.slidingWindows[tenant]; ok {
			return window.Allow()
		}
	case "leaky_bucket":
		if bucket, ok := arl.leakyBuckets[tenant]; ok {
			return bucket.Allow()
		}
	}
	return true
}

func (arl *AdvancedRateLimiter) SetupTokenBucket(tenant string, capacity, refillRate int) {
	arl.mu.Lock()
	defer arl.mu.Unlock()
	arl.buckets[tenant] = NewTokenBucket(capacity, refillRate)
}

func (arl *AdvancedRateLimiter) SetupSlidingWindow(tenant string, windowSize time.Duration, maxRequests int) {
	arl.mu.Lock()
	defer arl.mu.Unlock()
	arl.slidingWindows[tenant] = NewSlidingWindow(windowSize, maxRequests)
}

func (arl *AdvancedRateLimiter) SetupLeakyBucket(tenant string, capacity int, leakRate float64) {
	arl.mu.Lock()
	defer arl.mu.Unlock()
	arl.leakyBuckets[tenant] = NewLeakyBucket(capacity, leakRate)
}

// TenantManager manages tenants
type TenantManager struct {
	tenants map[string]*Tenant
	mu      sync.RWMutex
	store   TenantStore
}

// TenantStore interface for tenant persistence
type TenantStore interface {
	Get(ctx context.Context, id string) (*Tenant, error)
	Put(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*Tenant, error)
}

// NewTenantManager creates a new tenant manager
func NewTenantManager(store TenantStore) *TenantManager {
	return &TenantManager{
		tenants: make(map[string]*Tenant),
		store:   store,
	}
}

func (tm *TenantManager) Create(tenant *Tenant) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	if _, exists := tm.tenants[tenant.ID]; exists {
		return fmt.Errorf("tenant %s already exists", tenant.ID)
	}
	
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()
	tenant.Active = true
	
	tm.tenants[tenant.ID] = tenant
	
	if tm.store != nil {
		return tm.store.Put(context.Background(), tenant)
	}
	return nil
}

func (tm *TenantManager) Get(id string) (*Tenant, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	tenant, ok := tm.tenants[id]
	return tenant, ok
}

func (tm *TenantManager) Update(tenant *Tenant) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	if _, exists := tm.tenants[tenant.ID]; !exists {
		return fmt.Errorf("tenant %s not found", tenant.ID)
	}
	
	tenant.UpdatedAt = time.Now()
	tm.tenants[tenant.ID] = tenant
	
	if tm.store != nil {
		return tm.store.Put(context.Background(), tenant)
	}
	return nil
}

func (tm *TenantManager) Delete(id string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	if _, exists := tm.tenants[id]; !exists {
		return fmt.Errorf("tenant %s not found", id)
	}
	
	delete(tm.tenants, id)
	
	if tm.store != nil {
		return tm.store.Delete(context.Background(), id)
	}
	return nil
}

func (tm *TenantManager) List() []*Tenant {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	
	tenants := make([]*Tenant, 0, len(tm.tenants))
	for _, t := range tm.tenants {
		tenants = append(tenants, t)
	}
	return tenants
}

func (tm *TenantManager) Count() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return len(tm.tenants)
}

func (tm *TenantManager) ActiveCount() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	
	count := 0
	for _, t := range tm.tenants {
		if t.Active {
			count++
		}
	}
	return count
}

// MessageQueue for queuing events
type MessageQueue struct {
	queue    chan Event
	consumers []Consumer
	mu       sync.RWMutex
	running  bool
	ctx      context.Context
	cancel   context.CancelFunc
}

// Consumer consumes messages
type Consumer interface {
	Consume(event Event) error
	Name() string
}

// NewMessageQueue creates a new message queue
func NewMessageQueue(bufferSize int) *MessageQueue {
	return &MessageQueue{
		queue:     make(chan Event, bufferSize),
		consumers: make([]Consumer, 0),
	}
}

func (mq *MessageQueue) Start(ctx context.Context) {
	mq.ctx, mq.cancel = context.WithCancel(ctx)
	mq.running = true
	go mq.run()
}

func (mq *MessageQueue) Stop() {
	mq.running = false
	if mq.cancel != nil {
		mq.cancel()
	}
	close(mq.queue)
}

func (mq *MessageQueue) Enqueue(event Event) bool {
	select {
	case mq.queue <- event:
		return true
	default:
		return false
	}
}

func (mq *MessageQueue) AddConsumer(consumer Consumer) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.consumers = append(mq.consumers, consumer)
}

func (mq *MessageQueue) run() {
	for {
		select {
		case <-mq.ctx.Done():
			return
		case event, ok := <-mq.queue:
			if !ok {
				return
			}
			mq.mu.RLock()
			consumers := make([]Consumer, len(mq.consumers))
			copy(consumers, mq.consumers)
			mq.mu.RUnlock()
			
			for _, consumer := range consumers {
				if err := consumer.Consume(event); err != nil {
					// Log error
				}
			}
		}
	}
}

func (mq *MessageQueue) Size() int {
	return len(mq.queue)
}

func (mq *MessageQueue) IsRunning() bool {
	mq.mu.RLock()
	defer mq.mu.RUnlock()
	return mq.running
}

// InMemoryTenantStore implements TenantStore
type InMemoryTenantStore struct {
	tenants map[string]*Tenant
	mu      sync.RWMutex
}

func NewInMemoryTenantStore() *InMemoryTenantStore {
	return &InMemoryTenantStore{
		tenants: make(map[string]*Tenant),
	}
}

func (s *InMemoryTenantStore) Get(ctx context.Context, id string) (*Tenant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if tenant, ok := s.tenants[id]; ok {
		return tenant, nil
	}
	return nil, fmt.Errorf("tenant %s not found", id)
}

func (s *InMemoryTenantStore) Put(ctx context.Context, tenant *Tenant) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tenants[tenant.ID] = tenant
	return nil
}

func (s *InMemoryTenantStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tenants, id)
	return nil
}

func (s *InMemoryTenantStore) List(ctx context.Context) ([]*Tenant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tenants := make([]*Tenant, 0, len(s.tenants))
	for _, t := range s.tenants {
		tenants = append(tenants, t)
	}
	return tenants, nil
}

// Event transformer
type EventTransformer struct {
	rules []TransformRule
}

type TransformRule struct {
	Field     string
	Operation string // rename, set, delete, convert
	NewField  string
	Value     interface{}
}

func NewEventTransformer() *EventTransformer {
	return &EventTransformer{
		rules: make([]TransformRule, 0),
	}
}

func (et *EventTransformer) AddRule(rule TransformRule) {
	et.rules = append(et.rules, rule)
}

func (et *EventTransformer) Transform(event Event) Event {
	result := event
	for _, rule := range et.rules {
		switch rule.Operation {
		case "rename":
			if val, ok := result.Tags[rule.Field]; ok {
				delete(result.Tags, rule.Field)
				result.Tags[rule.NewField] = val
			}
		case "set":
			if result.Tags == nil {
				result.Tags = make(map[string]string)
			}
			result.Tags[rule.Field] = fmt.Sprintf("%v", rule.Value)
		case "delete":
			delete(result.Tags, rule.Field)
		case "convert":
			// Type conversion logic
		}
	}
	return result
}

// Event filter
type EventFilter struct {
	rules []FilterRule
}

type FilterRule struct {
	Field    string
	Operator string // eq, ne, contains, gt, lt, gte, lte
	Value    interface{}
	Action   string // drop, keep, modify
}

func NewEventFilter() *EventFilter {
	return &EventFilter{
		rules: make([]FilterRule, 0),
	}
}

func (ef *EventFilter) AddRule(rule FilterRule) {
	ef.rules = append(ef.rules, rule)
}

func (ef *EventFilter) Filter(event Event) bool {
	for _, rule := range ef.rules {
		if !ef.matches(event, rule) {
			if rule.Action == "drop" {
				return false
			}
		}
	}
	return true
}

func (ef *EventFilter) matches(event Event, rule FilterRule) bool {
	var value interface{}
	switch rule.Field {
	case "source":
		value = event.Source
	case "tenant":
		value = event.Tenant
	default:
		if v, ok := event.Tags[rule.Field]; ok {
			value = v
		} else if v, ok := event.Metrics[rule.Field]; ok {
			value = v
		}
	}
	
	switch rule.Operator {
	case "eq":
		return value == rule.Value
	case "ne":
		return value != rule.Value
	case "contains":
		if s, ok := value.(string); ok {
			if substr, ok := rule.Value.(string); ok {
				return containsStr(s, substr)
			}
		}
		return false
	case "gt":
		if v, ok := value.(float64); ok {
			if threshold, ok := rule.Value.(float64); ok {
				return v > threshold
			}
		}
		return false
	case "lt":
		if v, ok := value.(float64); ok {
			if threshold, ok := rule.Value.(float64); ok {
				return v < threshold
			}
		}
		return false
	case "gte":
		if v, ok := value.(float64); ok {
			if threshold, ok := rule.Value.(float64); ok {
				return v >= threshold
			}
		}
		return false
	case "lte":
		if v, ok := value.(float64); ok {
			if threshold, ok := rule.Value.(float64); ok {
				return v <= threshold
			}
		}
		return false
	}
	return true
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Event router
type EventRouter struct {
	routes []Route
	mu     sync.RWMutex
}

type Route struct {
	Name      string
	Condition func(Event) bool
	Target    string
	Priority  int
}

func NewEventRouter() *EventRouter {
	return &EventRouter{
		routes: make([]Route, 0),
	}
}

func (er *EventRouter) AddRoute(route Route) {
	er.mu.Lock()
	defer er.mu.Unlock()
	er.routes = append(er.routes, route)
	// Sort by priority
	// Simple insertion sort
	for i := len(er.routes) - 1; i > 0; i-- {
		if er.routes[i].Priority > er.routes[i-1].Priority {
			er.routes[i], er.routes[i-1] = er.routes[i-1], er.routes[i]
		} else {
			break
		}
	}
}

func (er *EventRouter) Route(event Event) []string {
	er.mu.RLock()
	defer er.mu.RUnlock()
	
	var targets []string
	for _, route := range er.routes {
		if route.Condition(event) {
			targets = append(targets, route.Target)
		}
	}
	return targets
}

func (er *EventRouter) RemoveRoute(name string) {
	er.mu.Lock()
	defer er.mu.Unlock()
	
	for i, route := range er.routes {
		if route.Name == name {
			er.routes = append(er.routes[:i], er.routes[i+1:]...)
			break
		}
	}
}

func (er *EventRouter) ListRoutes() []Route {
	er.mu.RLock()
	defer er.mu.RUnlock()
	routes := make([]Route, len(er.routes))
	copy(routes, er.routes)
	return routes
}