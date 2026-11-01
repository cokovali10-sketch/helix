// Package agent provides additional agent features
// including service mesh, load balancing, and failover
package agent

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ServiceMesh manages service mesh for agents
type ServiceMesh struct {
	services map[string]*Service
	mu       sync.RWMutex
}

// Service represents a service in the mesh
type Service struct {
	ID       string
	Name     string
	Endpoint string
	Weight   int
	Healthy  bool
	Metadata map[string]string
	LastSeen time.Time
}

// NewServiceMesh creates a new service mesh
func NewServiceMesh() *ServiceMesh {
	return &ServiceMesh{
		services: make(map[string]*Service),
	}
}

// Register registers a service
func (sm *ServiceMesh) Register(service *Service) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	service.LastSeen = time.Now()
	sm.services[service.ID] = service
}

// Unregister unregisters a service
func (sm *ServiceMesh) Unregister(id string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.services, id)
}

// Get gets a service by ID
func (sm *ServiceMesh) Get(id string) (*Service, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	service, ok := sm.services[id]
	return service, ok
}

// List lists all services
func (sm *ServiceMesh) List() []*Service {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	services := make([]*Service, 0, len(sm.services))
	for _, s := range sm.services {
		services = append(services, s)
	}
	return services
}

// GetHealthy gets healthy services
func (sm *ServiceMesh) GetHealthy() []*Service {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	var healthy []*Service
	for _, s := range sm.services {
		if s.Healthy {
			healthy = append(healthy, s)
		}
	}
	return healthy
}

// LoadBalancer implements load balancing
type LoadBalancer struct {
	strategy string
	mesh     *ServiceMesh
	mu       sync.RWMutex
}

// NewLoadBalancer creates a new load balancer
func NewLoadBalancer(mesh *ServiceMesh, strategy string) *LoadBalancer {
	return &LoadBalancer{
		strategy: strategy,
		mesh:     mesh,
	}
}

// Select selects a service
func (lb *LoadBalancer) Select(serviceName string) (*Service, error) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	healthy := lb.mesh.GetHealthy()
	if len(healthy) == 0 {
		return nil, fmt.Errorf("no healthy services")
	}

	switch lb.strategy {
	case "round_robin":
		return lb.roundRobin(healthy), nil
	case "random":
		return lb.random(healthy), nil
	case "least_connections":
		return lb.leastConnections(healthy), nil
	case "weighted":
		return lb.weighted(healthy), nil
	default:
		return lb.roundRobin(healthy), nil
	}
}

func (lb *LoadBalancer) roundRobin(services []*Service) *Service {
	// Simplified round robin
	return services[0]
}

func (lb *LoadBalancer) random(services []*Service) *Service {
	return services[rand.Intn(len(services))]
}

func (lb *LoadBalancer) leastConnections(services []*Service) *Service {
	// Simplified - would track connections
	return services[0]
}

func (lb *LoadBalancer) weighted(services []*Service) *Service {
	totalWeight := 0
	for _, s := range services {
		totalWeight += s.Weight
	}

	r := rand.Intn(totalWeight)
	current := 0
	for _, s := range services {
		current += s.Weight
		if r < current {
			return s
		}
	}
	return services[0]
}

// FailoverManager manages failover
type FailoverManager struct {
	primary   *Service
	backup    *Service
	mesh      *ServiceMesh
	mu        sync.RWMutex
	failover  bool
}

// NewFailoverManager creates a new failover manager
func NewFailoverManager(mesh *ServiceMesh) *FailoverManager {
	return &FailoverManager{
		mesh: mesh,
	}
}

// SetPrimary sets the primary service
func (fm *FailoverManager) SetPrimary(service *Service) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.primary = service
}

// SetBackup sets the backup service
func (fm *FailoverManager) SetBackup(service *Service) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.backup = service
}

// Failover performs failover
func (fm *FailoverManager) Failover() error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if fm.primary == nil || fm.backup == nil {
		return fmt.Errorf("primary or backup not set")
	}

	fm.primary.Healthy = false
	fm.backup.Healthy = true
	fm.failover = true
	return nil
}

// IsFailover returns if in failover mode
func (fm *FailoverManager) IsFailover() bool {
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	return fm.failover
}

// Recover recovers from failover
func (fm *FailoverManager) Recover() {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if fm.primary != nil {
		fm.primary.Healthy = true
	}
	fm.failover = false
}

// CircuitBreaker implements circuit breaker pattern
type CircuitBreaker struct {
	name        string
	maxFailures int
	timeout     time.Duration
	state       CircuitState
	failures    int
	lastFailure time.Time
	mu          sync.Mutex
}

// CircuitState represents circuit breaker state
type CircuitState string

const (
	CircuitClosed   CircuitState = "closed"
	CircuitOpen     CircuitState = "open"
	CircuitHalfOpen CircuitState = "half_open"
)

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(name string, maxFailures int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		name:        name,
		maxFailures: maxFailures,
		timeout:     timeout,
		state:       CircuitClosed,
	}
}

// Execute executes a function with circuit breaker
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == CircuitOpen {
		if time.Since(cb.lastFailure) > cb.timeout {
			cb.state = CircuitHalfOpen
		} else {
			return fmt.Errorf("circuit breaker is open")
		}
	}

	err := fn()
	if err != nil {
		cb.failures++
		cb.lastFailure = time.Now()
		if cb.failures >= cb.maxFailures {
			cb.state = CircuitOpen
		}
		return err
	}

	if cb.state == CircuitHalfOpen {
		cb.state = CircuitClosed
		cb.failures = 0
	}
	return nil
}

// GetState returns the current state
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Reset resets the circuit breaker
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = CircuitClosed
	cb.failures = 0
}

// RetryPolicy implements retry logic
type RetryPolicy struct {
	MaxRetries    int
	InitialDelay  time.Duration
	MaxDelay      time.Duration
	BackoffFactor  float64
	Jitter        bool
}

// NewRetryPolicy creates a new retry policy
func NewRetryPolicy(maxRetries int, initialDelay, maxDelay time.Duration, backoffFactor float64) *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:   maxRetries,
		InitialDelay:  initialDelay,
		MaxDelay:     maxDelay,
		BackoffFactor: backoffFactor,
		Jitter:       true,
	}
}

// Execute executes a function with retry
func (rp *RetryPolicy) Execute(fn func() error) error {
	var lastErr error
	delay := rp.InitialDelay

	for attempt := 0; attempt <= rp.MaxRetries; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}
		lastErr = err

		if attempt < rp.MaxRetries {
			if rp.Jitter {
				delay = time.Duration(float64(delay) * (1 + rand.Float64()*0.1))
			}
			time.Sleep(delay)
			delay = time.Duration(float64(delay) * rp.BackoffFactor)
			if delay > rp.MaxDelay {
				delay = rp.MaxDelay
			}
		}
	}

	return lastErr
}

// Bulkhead implements bulkhead pattern
type Bulkhead struct {
	name     string
	maxConc  int
	sem      chan struct{}
	mu       sync.Mutex
}

// NewBulkhead creates a new bulkhead
func NewBulkhead(name string, maxConcurrent int) *Bulkhead {
	return &Bulkhead{
		name:    name,
		maxConc: maxConcurrent,
		sem:     make(chan struct{}, maxConcurrent),
	}
}

// Execute executes a function with bulkhead
func (b *Bulkhead) Execute(fn func() error) error {
	select {
	case b.sem <- struct{}{}:
		defer func() { <-b.sem }()
		return fn()
	default:
		return fmt.Errorf("bulkhead %s is full", b.name)
	}
}

// GetAvailable returns available slots
func (b *Bulkhead) GetAvailable() int {
	return b.maxConc - len(b.sem)
}

// GetMax returns max concurrent
func (b *Bulkhead) GetMax() int {
	return b.maxConc
}

// RateLimiter implements rate limiting
type RateLimiter struct {
	rate     int
	burst    int
	tokens   int
	lastTime time.Time
	mu       sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, burst int) *RateLimiter {
	return &RateLimiter{
		rate:     rate,
		burst:    burst,
		tokens:   burst,
		lastTime: time.Now(),
	}
}

// Allow checks if request is allowed
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastTime)
	tokensToAdd := int(elapsed.Seconds()) * rl.rate

	if tokensToAdd > 0 {
		rl.tokens = min(rl.tokens+tokensToAdd, rl.burst)
		rl.lastTime = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

// GetTokens returns available tokens
func (rl *RateLimiter) GetTokens() int {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.tokens
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}