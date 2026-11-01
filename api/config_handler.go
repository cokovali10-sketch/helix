// Package api provides additional API handlers
// including configuration, settings, and preferences
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// ConfigHandler handles configuration endpoints
type ConfigHandler struct {
	config map[string]interface{}
	mu     sync.RWMutex
}

// NewConfigHandler creates a new config handler
func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{
		config: make(map[string]interface{}),
	}
}

// Set sets a config value
func (ch *ConfigHandler) Set(key string, value interface{}) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	ch.config[key] = value
}

// Get gets a config value
func (ch *ConfigHandler) Get(key string) (interface{}, bool) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	v, ok := ch.config[key]
	return v, ok
}

// GetAll gets all config values
func (ch *ConfigHandler) GetAll() map[string]interface{} {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	result := make(map[string]interface{})
	for k, v := range ch.config {
		result[k] = v
	}
	return result
}

// Delete deletes a config value
func (ch *ConfigHandler) Delete(key string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	delete(ch.config, key)
}

// ServeHTTP serves config endpoint
func (ch *ConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		key := r.URL.Query().Get("key")
		if key != "" {
			value, ok := ch.Get(key)
			if !ok {
				http.Error(w, "Config not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"key":   key,
				"value": value,
			})
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ch.GetAll())
		}
	case http.MethodPost:
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for key, value := range body {
			ch.Set(key, value)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	case http.MethodDelete:
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "Key required", http.StatusBadRequest)
			return
		}
		ch.Delete(key)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// SettingsHandler handles settings endpoints
type SettingsHandler struct {
	settings map[string]interface{}
	mu       sync.RWMutex
}

// NewSettingsHandler creates a new settings handler
func NewSettingsHandler() *SettingsHandler {
	return &SettingsHandler{
		settings: make(map[string]interface{}),
	}
}

// Get gets a setting
func (sh *SettingsHandler) Get(key string) (interface{}, bool) {
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	v, ok := sh.settings[key]
	return v, ok
}

// Set sets a setting
func (sh *SettingsHandler) Set(key string, value interface{}) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	sh.settings[key] = value
}

// GetAll gets all settings
func (sh *SettingsHandler) GetAll() map[string]interface{} {
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	result := make(map[string]interface{})
	for k, v := range sh.settings {
		result[k] = v
	}
	return result
}

// ServeHTTP serves settings endpoint
func (sh *SettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		key := r.URL.Query().Get("key")
		if key != "" {
			value, ok := sh.Get(key)
			if !ok {
				http.Error(w, "Setting not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"key":   key,
				"value": value,
			})
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(sh.GetAll())
		}
	case http.MethodPost, http.MethodPut:
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for key, value := range body {
			sh.Set(key, value)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// PreferencesHandler handles user preferences
type PreferencesHandler struct {
	preferences map[string]map[string]interface{}
	mu          sync.RWMutex
}

// NewPreferencesHandler creates a new preferences handler
func NewPreferencesHandler() *PreferencesHandler {
	return &PreferencesHandler{
		preferences: make(map[string]map[string]interface{}),
	}
}

// Get gets preferences for a user
func (ph *PreferencesHandler) Get(userID string) map[string]interface{} {
	ph.mu.RLock()
	defer ph.mu.RUnlock()
	if prefs, ok := ph.preferences[userID]; ok {
		result := make(map[string]interface{})
		for k, v := range prefs {
			result[k] = v
		}
		return result
	}
	return make(map[string]interface{})
}

// Set sets a preference for a user
func (ph *PreferencesHandler) Set(userID, key string, value interface{}) {
	ph.mu.Lock()
	defer ph.mu.Unlock()

	if ph.preferences[userID] == nil {
		ph.preferences[userID] = make(map[string]interface{})
	}
	ph.preferences[userID][key] = value
}

// SetAll sets all preferences for a user
func (ph *PreferencesHandler) SetAll(userID string, prefs map[string]interface{}) {
	ph.mu.Lock()
	defer ph.mu.Unlock()

	if ph.preferences[userID] == nil {
		ph.preferences[userID] = make(map[string]interface{})
	}
	for k, v := range prefs {
		ph.preferences[userID][k] = v
	}
}

// Delete deletes a preference for a user
func (ph *PreferencesHandler) Delete(userID, key string) {
	ph.mu.Lock()
	defer ph.mu.Unlock()

	if prefs, ok := ph.preferences[userID]; ok {
		delete(prefs, key)
	}
}

// ServeHTTP serves preferences endpoint
func (ph *PreferencesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userId required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		prefs := ph.Get(userID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prefs)
	case http.MethodPost, http.MethodPut:
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ph.SetAll(userID, body)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	case http.MethodDelete:
		key := r.URL.Query().Get("key")
		if key != "" {
			ph.Delete(userID, key)
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// FeatureFlagHandler handles feature flags
type FeatureFlagHandler struct {
	flags map[string]bool
	mu    sync.RWMutex
}

// NewFeatureFlagHandler creates a new feature flag handler
func NewFeatureFlagHandler() *FeatureFlagHandler {
	return &FeatureFlagHandler{
		flags: make(map[string]bool),
	}
}

// IsEnabled checks if a feature is enabled
func (ffh *FeatureFlagHandler) IsEnabled(name string) bool {
	ffh.mu.RLock()
	defer ffh.mu.RUnlock()
	return ffh.flags[name]
}

// Set sets a feature flag
func (ffh *FeatureFlagHandler) Set(name string, enabled bool) {
	ffh.mu.Lock()
	defer ffh.mu.Unlock()
	ffh.flags[name] = enabled
}

// GetAll gets all feature flags
func (ffh *FeatureFlagHandler) GetAll() map[string]bool {
	ffh.mu.RLock()
	defer ffh.mu.RUnlock()
	result := make(map[string]bool)
	for k, v := range ffh.flags {
		result[k] = v
	}
	return result
}

// ServeHTTP serves feature flags endpoint
func (ffh *FeatureFlagHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		name := r.URL.Query().Get("name")
		if name != "" {
			enabled := ffh.IsEnabled(name)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"name":    name,
				"enabled": enabled,
			})
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ffh.GetAll())
		}
	case http.MethodPost, http.MethodPut:
		var body map[string]bool
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for name, enabled := range body {
			ffh.Set(name, enabled)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// RateLimitHandler handles rate limit endpoints
type RateLimitHandler struct {
	limits map[string]*RateLimit
	mu     sync.RWMutex
}

// RateLimit represents a rate limit
type RateLimit struct {
	Tenant    string `json:"tenant"`
	Limit     int    `json:"limit"`
	Remaining int    `json:"remaining"`
	ResetAt   time.Time `json:"resetAt"`
}

// NewRateLimitHandler creates a new rate limit handler
func NewRateLimitHandler() *RateLimitHandler {
	return &RateLimitHandler{
		limits: make(map[string]*RateLimit),
	}
}

// Get gets a rate limit
func (rlh *RateLimitHandler) Get(tenant string) (*RateLimit, bool) {
	rlh.mu.RLock()
	defer rlh.mu.RUnlock()
	limit, ok := rlh.limits[tenant]
	return limit, ok
}

// Set sets a rate limit
func (rlh *RateLimitHandler) Set(tenant string, limit int) {
	rlh.mu.Lock()
	defer rlh.mu.Unlock()
	rlh.limits[tenant] = &RateLimit{
		Tenant:    tenant,
		Limit:     limit,
		Remaining: limit,
		ResetAt:   time.Now().Add(time.Hour),
	}
}

// ServeHTTP serves rate limit endpoint
func (rlh *RateLimitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tenant := r.URL.Query().Get("tenant")
	if tenant == "" {
		http.Error(w, "Tenant required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		limit, ok := rlh.Get(tenant)
		if !ok {
			http.Error(w, "Rate limit not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(limit)
	case http.MethodPost, http.MethodPut:
		var body struct {
			Limit int `json:"limit"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rlh.Set(tenant, body.Limit)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// AuditLogHandler handles audit log endpoints
type AuditLogHandler struct {
	logs []AuditLog
	mu   sync.RWMutex
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"userAgent"`
	Details   string    `json:"details"`
}

// NewAuditLogHandler creates a new audit log handler
func NewAuditLogHandler() *AuditLogHandler {
	return &AuditLogHandler{
		logs: make([]AuditLog, 0),
	}
}

// Log adds an audit log entry
func (alh *AuditLogHandler) Log(userID, action, resource, ip, userAgent, details string) {
	alh.mu.Lock()
	defer alh.mu.Unlock()

	log := AuditLog{
		ID:        fmt.Sprintf("audit-%d", time.Now().UnixNano()),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
		Timestamp: time.Now(),
		IP:        ip,
		UserAgent: userAgent,
		Details:   details,
	}
	alh.logs = append(alh.logs, log)
}

// List lists audit logs
func (alh *AuditLogHandler) List(userID, action, resource string, limit int) []AuditLog {
	alh.mu.RLock()
	defer alh.mu.RUnlock()

	var result []AuditLog
	for _, log := range alh.logs {
		if userID != "" && log.UserID != userID {
			continue
		}
		if action != "" && log.Action != action {
			continue
		}
		if resource != "" && log.Resource != resource {
			continue
		}
		result = append(result, log)
		if limit > 0 && len(result) >= limit {
			break
		}
	}
	return result
}

// ServeHTTP serves audit log endpoint
func (alh *AuditLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("userId")
	action := r.URL.Query().Get("action")
	resource := r.URL.Query().Get("resource")
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	logs := alh.List(userID, action, resource, limit)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"logs":  logs,
		"total": len(logs),
	})
}