// Package agent provides load balancer
package agent

import (
	"sync"
)

// LoadBalancer balances load
type LoadBalancer struct {
	targets []string
	index   int
	mu      sync.Mutex
}

// NewLoadBalancer creates balancer
func NewLoadBalancer(targets []string) *LoadBalancer {
	return &LoadBalancer{targets: targets}
}

// Next returns next target
func (lb *LoadBalancer) Next() string {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	if len(lb.targets) == 0 {
		return ""
	}
	target := lb.targets[lb.index%len(lb.targets)]
	lb.index++
	return target
}

// Add adds target
func (lb *LoadBalancer) Add(target string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.targets = append(lb.targets, target)
}

// Remove removes target
func (lb *LoadBalancer) Remove(target string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	out := []string{}
	for _, t := range lb.targets {
		if t != target {
			out = append(out, t)
		}
	}
	lb.targets = out
}

// Targets returns targets
func (lb *LoadBalancer) Targets() []string {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	out := make([]string, len(lb.targets))
	copy(out, lb.targets)
	return out
}

// Count returns count
func (lb *LoadBalancer) Count() int {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	return len(lb.targets)
}
