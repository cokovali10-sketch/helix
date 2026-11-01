// Package catalog provides permission management
package catalog

import (
	"fmt"
	"sync"
)

// Permission represents a permission
type Permission struct {
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

// PermissionManager manages permissions
type PermissionManager struct {
	perms map[string]*Permission
	mu    sync.RWMutex
}

// NewPermissionManager creates manager
func NewPermissionManager() *PermissionManager {
	return &PermissionManager{perms: make(map[string]*Permission)}
}

// Create creates permission
func (pm *PermissionManager) Create(name, resource, action, desc string) *Permission {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	p := &Permission{Name: name, Resource: resource, Action: action, Description: desc}
	pm.perms[name] = p
	return p
}

// Get returns permission
func (pm *PermissionManager) Get(name string) (*Permission, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	p, ok := pm.perms[name]
	return p, ok
}

// Delete deletes permission
func (pm *PermissionManager) Delete(name string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if _, ok := pm.perms[name]; !ok {
		return false
	}
	delete(pm.perms, name)
	return true
}

// List lists permissions
func (pm *PermissionManager) List() []*Permission {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	out := []*Permission{}
	for _, p := range pm.perms {
		out = append(out, p)
	}
	return out
}

// ByResource returns permissions by resource
func (pm *PermissionManager) ByResource(resource string) []*Permission {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	out := []*Permission{}
	for _, p := range pm.perms {
		if p.Resource == resource {
			out = append(out, p)
		}
	}
	return out
}

// Count returns count
func (pm *PermissionManager) Count() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return len(pm.perms)
}

// Exists checks existence
func (pm *PermissionManager) Exists(name string) bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	_, ok := pm.perms[name]
	return ok
}

// Validate validates permission
func (pm *PermissionManager) Validate(name string) error {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	if _, ok := pm.perms[name]; !ok {
		return fmt.Errorf("permission not found")
	}
	return nil
}
