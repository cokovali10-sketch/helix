// Package agent provides service registry
package agent

import (
	"fmt"
	"sync"
	"time"
)

// ServiceInstance represents service instance
type ServiceInstance struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Metadata map[string]string `json:"metadata"`
	RegisteredAt time.Time     `json:"registered_at"`
}

// ServiceRegistry manages services
type ServiceRegistry struct {
	services map[string][]*ServiceInstance
	mu       sync.RWMutex
}

// NewServiceRegistry creates registry
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{services: make(map[string][]*ServiceInstance)}
}

// Register registers instance
func (sr *ServiceRegistry) Register(inst *ServiceInstance) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	inst.RegisteredAt = time.Now()
	sr.services[inst.Name] = append(sr.services[inst.Name], inst)
}

// Deregister deregisters
func (sr *ServiceRegistry) Deregister(name, id string) bool {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	list := sr.services[name]
	for i, inst := range list {
		if inst.ID == id {
			sr.services[name] = append(list[:i], list[i+1:]...)
			return true
		}
	}
	return false
}

// Discover returns instances
func (sr *ServiceRegistry) Discover(name string) []*ServiceInstance {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	out := []*ServiceInstance{}
	for _, inst := range sr.services[name] {
		out = append(out, inst)
	}
	return out
}

// Services returns service names
func (sr *ServiceRegistry) Services() []string {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	out := []string{}
	for name := range sr.services {
		out = append(out, name)
	}
	return out
}

// Count returns total instances
func (sr *ServiceRegistry) Count() int {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	n := 0
	for _, list := range sr.services {
		n += len(list)
	}
	return n
}

// String returns summary
func (sr *ServiceRegistry) String() string {
	return fmt.Sprintf("ServiceRegistry: %d services", len(sr.services))
}
