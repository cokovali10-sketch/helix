// Package ingestion provides routing
package ingestion

import (
	"fmt"
	"sync"
)

// Route represents a route
type Route struct {
	Name     string
	Matcher  func(map[string]interface{}) bool
	Target   string
	Priority int
}

// Router routes records
type Router struct {
	routes []*Route
	mu     sync.RWMutex
}

// NewRouter creates router
func NewRouter() *Router {
	return &Router{routes: []*Route{}}
}

// AddRoute adds route
func (r *Router) AddRoute(route *Route) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes = append(r.routes, route)
}

// Route returns target for record
func (r *Router) Route(record map[string]interface{}) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, route := range r.routes {
		if route.Matcher(record) {
			return route.Target
		}
	}
	return "default"
}

// Routes returns routes
func (r *Router) Routes() []*Route {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Route, len(r.routes))
	copy(out, r.routes)
	return out
}

// RemoveRoute removes route
func (r *Router) RemoveRoute(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, route := range r.routes {
		if route.Name == name {
			r.routes = append(r.routes[:i], r.routes[i+1:]...)
			return true
		}
	}
	return false
}

// Count returns count
func (r *Router) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.routes)
}

// FieldRoute creates route by field value
func FieldRoute(field, value, target string) *Route {
	return &Route{
		Name:   fmt.Sprintf("%s=%v", field, value),
		Target: target,
		Matcher: func(rec map[string]interface{}) bool {
			v, ok := rec[field]
			if !ok {
				return false
			}
			return fmt.Sprintf("%v", v) == value
		},
	}
}
