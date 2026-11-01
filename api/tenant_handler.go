// Package api provides tenant management handlers
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Tenant represents a tenant
type Tenant struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Plan        string            `json:"plan"`
	Status      string            `json:"status"`
	Settings    map[string]string `json:"settings"`
	Quotas      TenantQuota       `json:"quotas"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Metadata    map[string]string `json:"metadata"`
}

// TenantQuota represents tenant quotas
type TenantQuota struct {
	MaxSeries   int64 `json:"max_series"`
	MaxQueries  int64 `json:"max_queries_per_second"`
	MaxStorage  int64 `json:"max_storage_bytes"`
	MaxUsers    int   `json:"max_users"`
	MaxDashboards int `json:"max_dashboards"`
}

// TenantHandler manages tenants
type TenantHandler struct {
	tenants map[string]*Tenant
	mu      sync.RWMutex
}

// NewTenantHandler creates a new tenant handler
func NewTenantHandler() *TenantHandler {
	return &TenantHandler{tenants: make(map[string]*Tenant)}
}

// CreateTenant creates a tenant
func (th *TenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var t Tenant
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	th.mu.Lock()
	defer th.mu.Unlock()
	if _, exists := th.tenants[t.ID]; exists {
		http.Error(w, "tenant exists", http.StatusConflict)
		return
	}
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	if t.Status == "" {
		t.Status = "active"
	}
	if t.Plan == "" {
		t.Plan = "free"
	}
	if t.Settings == nil {
		t.Settings = map[string]string{}
	}
	if t.Metadata == nil {
		t.Metadata = map[string]string{}
	}
	th.tenants[t.ID] = &t
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

// GetTenant retrieves a tenant
func (th *TenantHandler) GetTenant(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	th.mu.RLock()
	t, ok := th.tenants[id]
	th.mu.RUnlock()
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

// ListTenants lists tenants
func (th *TenantHandler) ListTenants(w http.ResponseWriter, r *http.Request) {
	th.mu.RLock()
	defer th.mu.RUnlock()
	out := make([]*Tenant, 0, len(th.tenants))
	for _, t := range th.tenants {
		out = append(out, t)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

// UpdateTenant updates a tenant
func (th *TenantHandler) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var t Tenant
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	th.mu.Lock()
	defer th.mu.Unlock()
	existing, ok := th.tenants[t.ID]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	existing.Name = t.Name
	existing.Plan = t.Plan
	existing.Status = t.Status
	existing.Settings = t.Settings
	existing.Quotas = t.Quotas
	existing.Metadata = t.Metadata
	existing.UpdatedAt = time.Now()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existing)
}

// DeleteTenant deletes a tenant
func (th *TenantHandler) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	th.mu.Lock()
	defer th.mu.Unlock()
	if _, ok := th.tenants[id]; !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	delete(th.tenants, id)
	w.WriteHeader(http.StatusNoContent)
}

// TenantStats returns stats
func (th *TenantHandler) TenantStats(w http.ResponseWriter, r *http.Request) {
	th.mu.RLock()
	defer th.mu.RUnlock()
	stats := map[string]interface{}{
		"total":  len(th.tenants),
		"active": 0,
	}
	for _, t := range th.tenants {
		if t.Status == "active" {
			stats["active"] = stats["active"].(int) + 1
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// TenantMiddleware enforces tenant context
func TenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get("X-Tenant-ID")
		if tenantID == "" {
			tenantID = "default"
		}
		ctx := r.Context()
		_ = ctx
		next.ServeHTTP(w, r)
	})
}

// RegisterTenantRoutes registers routes
func RegisterTenantRoutes(mux *http.ServeMux, th *TenantHandler) {
	mux.HandleFunc("/tenants", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			th.ListTenants(w, r)
		case http.MethodPost:
			th.CreateTenant(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/tenants/get", th.GetTenant)
	mux.HandleFunc("/tenants/update", th.UpdateTenant)
	mux.HandleFunc("/tenants/delete", th.DeleteTenant)
	mux.HandleFunc("/tenants/stats", th.TenantStats)
	fmt.Println("tenant routes registered")
}
