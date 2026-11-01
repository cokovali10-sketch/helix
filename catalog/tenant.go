// Package catalog provides tenant management
package catalog

import (
	"fmt"
	"sync"
	"time"
)

// TenantInfo represents tenant info
type TenantInfo struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Plan      string            `json:"plan"`
	Status    string            `json:"status"`
	Settings  map[string]string `json:"settings"`
	CreatedAt time.Time         `json:"created_at"`
}

// TenantManager manages tenants
type TenantManager struct {
	tenants map[string]*TenantInfo
	mu      sync.RWMutex
}

// NewTenantManager creates manager
func NewTenantManager() *TenantManager {
	return &TenantManager{tenants: make(map[string]*TenantInfo)}
}

// Create creates tenant
func (tm *TenantManager) Create(id, name, plan string) *TenantInfo {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	t := &TenantInfo{
		ID:        id,
		Name:      name,
		Plan:      plan,
		Status:    "active",
		Settings:  map[string]string{},
		CreatedAt: time.Now(),
	}
	tm.tenants[id] = t
	return t
}

// Get returns tenant
func (tm *TenantManager) Get(id string) (*TenantInfo, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	t, ok := tm.tenants[id]
	return t, ok
}

// Delete deletes tenant
func (tm *TenantManager) Delete(id string) bool {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if _, ok := tm.tenants[id]; !ok {
		return false
	}
	delete(tm.tenants, id)
	return true
}

// List lists tenants
func (tm *TenantManager) List() []*TenantInfo {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	out := []*TenantInfo{}
	for _, t := range tm.tenants {
		out = append(out, t)
	}
	return out
}

// Count returns count
func (tm *TenantManager) Count() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return len(tm.tenants)
}

// UpdateStatus updates status
func (tm *TenantManager) UpdateStatus(id, status string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	t, ok := tm.tenants[id]
	if !ok {
		return fmt.Errorf("tenant not found")
	}
	t.Status = status
	return nil
}
