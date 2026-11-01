// Package catalog provides role management
package catalog

import (
	"fmt"
	"sync"
)

// Role represents a role
type Role struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	Description string   `json:"description"`
}

// RoleManager manages roles
type RoleManager struct {
	roles map[string]*Role
	mu    sync.RWMutex
}

// NewRoleManager creates manager
func NewRoleManager() *RoleManager {
	return &RoleManager{roles: make(map[string]*Role)}
}

// Create creates role
func (rm *RoleManager) Create(name, desc string, perms []string) *Role {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	r := &Role{Name: name, Description: desc, Permissions: perms}
	rm.roles[name] = r
	return r
}

// Get returns role
func (rm *RoleManager) Get(name string) (*Role, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	r, ok := rm.roles[name]
	return r, ok
}

// Delete deletes role
func (rm *RoleManager) Delete(name string) bool {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if _, ok := rm.roles[name]; !ok {
		return false
	}
	delete(rm.roles, name)
	return true
}

// List lists roles
func (rm *RoleManager) List() []*Role {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	out := []*Role{}
	for _, r := range rm.roles {
		out = append(out, r)
	}
	return out
}

// HasPermission checks permission
func (rm *RoleManager) HasPermission(roleName, perm string) bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	r, ok := rm.roles[roleName]
	if !ok {
		return false
	}
	for _, p := range r.Permissions {
		if p == perm || p == "*" {
			return true
		}
	}
	return false
}

// AddPermission adds permission
func (rm *RoleManager) AddPermission(roleName, perm string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	r, ok := rm.roles[roleName]
	if !ok {
		return fmt.Errorf("role not found")
	}
	r.Permissions = append(r.Permissions, perm)
	return nil
}

// RemovePermission removes permission
func (rm *RoleManager) RemovePermission(roleName, perm string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	r, ok := rm.roles[roleName]
	if !ok {
		return
	}
	out := []string{}
	for _, p := range r.Permissions {
		if p != perm {
			out = append(out, p)
		}
	}
	r.Permissions = out
}

// Count returns count
func (rm *RoleManager) Count() int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return len(rm.roles)
}
