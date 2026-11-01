// Package api provides audit logging handlers
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// AuditEntry represents an audit log entry
type AuditEntry struct {
	ID        string            `json:"id"`
	Actor     string            `json:"actor"`
	Action    string            `json:"action"`
	Resource  string            `json:"resource"`
	Details   map[string]string `json:"details"`
	IP        string            `json:"ip"`
	UserAgent string            `json:"user_agent"`
	Timestamp time.Time         `json:"timestamp"`
	Status    string            `json:"status"`
}

// AuditHandler manages audit logs
type AuditHandler struct {
	entries []AuditEntry
	mu      sync.RWMutex
}

// NewAuditHandler creates a new audit handler
func NewAuditHandler() *AuditHandler {
	return &AuditHandler{entries: []AuditEntry{}}
}

// Log records an audit entry
func (ah *AuditHandler) Log(e AuditEntry) {
	ah.mu.Lock()
	defer ah.mu.Unlock()
	if e.ID == "" {
		e.ID = fmt.Sprintf("audit-%d", time.Now().UnixNano())
	}
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	ah.entries = append(ah.entries, e)
}

// List returns audit entries
func (ah *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	ah.mu.RLock()
	defer ah.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ah.entries)
}

// Filter returns filtered entries
func (ah *AuditHandler) Filter(w http.ResponseWriter, r *http.Request) {
	actor := r.URL.Query().Get("actor")
	action := r.URL.Query().Get("action")
	ah.mu.RLock()
	defer ah.mu.RUnlock()
	out := []AuditEntry{}
	for _, e := range ah.entries {
		if actor != "" && e.Actor != actor {
			continue
		}
		if action != "" && e.Action != action {
			continue
		}
		out = append(out, e)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

// Stats returns audit stats
func (ah *AuditHandler) Stats(w http.ResponseWriter, r *http.Request) {
	ah.mu.RLock()
	defer ah.mu.RUnlock()
	byAction := map[string]int{}
	for _, e := range ah.entries {
		byAction[e.Action]++
	}
	stats := map[string]interface{}{
		"total":     len(ah.entries),
		"by_action": byAction,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// Clear clears entries
func (ah *AuditHandler) Clear(w http.ResponseWriter, r *http.Request) {
	ah.mu.Lock()
	defer ah.mu.Unlock()
	ah.entries = []AuditEntry{}
	w.WriteHeader(http.StatusNoContent)
}

// RegisterAuditRoutes registers routes
func RegisterAuditRoutes(mux *http.ServeMux, ah *AuditHandler) {
	mux.HandleFunc("/audit", ah.List)
	mux.HandleFunc("/audit/filter", ah.Filter)
	mux.HandleFunc("/audit/stats", ah.Stats)
	mux.HandleFunc("/audit/clear", ah.Clear)
	fmt.Println("audit routes registered")
}
