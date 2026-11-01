// Package catalog provides audit logging
package catalog

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// CatalogAuditEntry represents catalog audit
type CatalogAuditEntry struct {
	ID        string    `json:"id"`
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`
	Entity    string    `json:"entity"`
	EntityID  string    `json:"entity_id"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details"`
}

// CatalogAuditor audits catalog operations
type CatalogAuditor struct {
	entries []CatalogAuditEntry
	mu      sync.RWMutex
}

// NewCatalogAuditor creates auditor
func NewCatalogAuditor() *CatalogAuditor {
	return &CatalogAuditor{entries: []CatalogAuditEntry{}}
}

// Record records an entry
func (ca *CatalogAuditor) Record(actor, action, entity, entityID, details string) {
	ca.mu.Lock()
	defer ca.mu.Unlock()
	ca.entries = append(ca.entries, CatalogAuditEntry{
		ID:        fmt.Sprintf("ca-%d", time.Now().UnixNano()),
		Actor:     actor,
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		Timestamp: time.Now(),
		Details:   details,
	})
}

// List returns entries
func (ca *CatalogAuditor) List() []CatalogAuditEntry {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	out := make([]CatalogAuditEntry, len(ca.entries))
	copy(out, ca.entries)
	return out
}

// FilterByActor filters
func (ca *CatalogAuditor) FilterByActor(actor string) []CatalogAuditEntry {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	out := []CatalogAuditEntry{}
	for _, e := range ca.entries {
		if e.Actor == actor {
			out = append(out, e)
		}
	}
	return out
}

// FilterByEntity filters
func (ca *CatalogAuditor) FilterByEntity(entity string) []CatalogAuditEntry {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	out := []CatalogAuditEntry{}
	for _, e := range ca.entries {
		if e.Entity == entity {
			out = append(out, e)
		}
	}
	return out
}

// Count returns count
func (ca *CatalogAuditor) Count() int {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	return len(ca.entries)
}

// Clear clears entries
func (ca *CatalogAuditor) Clear() {
	ca.mu.Lock()
	defer ca.mu.Unlock()
	ca.entries = []CatalogAuditEntry{}
}

// MarshalJSON
func (ca *CatalogAuditor) MarshalJSON() ([]byte, error) {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	return json.Marshal(ca.entries)
}
