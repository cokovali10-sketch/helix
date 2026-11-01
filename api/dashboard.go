// Package api provides additional API handlers
// including dashboard, panel, and widget handlers
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// DashboardHandler handles dashboard endpoints
type DashboardHandler struct {
	dashboards map[string]*Dashboard
	mu         sync.RWMutex
}

// Dashboard represents a dashboard
type Dashboard struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Panels      []Panel   `json:"panels"`
	Layout      string    `json:"layout"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Owner       string    `json:"owner"`
	IsPublic    bool      `json:"isPublic"`
	Tags        []string  `json:"tags"`
}

// Panel represents a dashboard panel
type Panel struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Query    string `json:"query"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Options  map[string]interface{} `json:"options"`
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{
		dashboards: make(map[string]*Dashboard),
	}
}

// Create creates a new dashboard
func (dh *DashboardHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var dashboard Dashboard
	if err := json.NewDecoder(r.Body).Decode(&dashboard); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dh.mu.Lock()
	defer dh.mu.Unlock()

	dashboard.ID = fmt.Sprintf("dashboard-%d", time.Now().UnixNano())
	dashboard.CreatedAt = time.Now()
	dashboard.UpdatedAt = time.Now()

	dh.dashboards[dashboard.ID] = &dashboard

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dashboard)
}

// Get gets a dashboard by ID
func (dh *DashboardHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	dh.mu.RLock()
	defer dh.mu.RUnlock()

	dashboard, ok := dh.dashboards[id]
	if !ok {
		http.Error(w, "Dashboard not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}

// List lists all dashboards
func (dh *DashboardHandler) List(w http.ResponseWriter, r *http.Request) {
	dh.mu.RLock()
	defer dh.mu.RUnlock()

	dashboards := make([]*Dashboard, 0, len(dh.dashboards))
	for _, d := range dh.dashboards {
		dashboards = append(dashboards, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"dashboards": dashboards,
		"total":      len(dashboards),
	})
}

// Update updates a dashboard
func (dh *DashboardHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	var updates Dashboard
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dh.mu.Lock()
	defer dh.mu.Unlock()

	dashboard, ok := dh.dashboards[id]
	if !ok {
		http.Error(w, "Dashboard not found", http.StatusNotFound)
		return
	}

	if updates.Title != "" {
		dashboard.Title = updates.Title
	}
	if updates.Description != "" {
		dashboard.Description = updates.Description
	}
	if updates.Panels != nil {
		dashboard.Panels = updates.Panels
	}
	if updates.Layout != "" {
		dashboard.Layout = updates.Layout
	}
	if updates.Tags != nil {
		dashboard.Tags = updates.Tags
	}
	dashboard.IsPublic = updates.IsPublic
	dashboard.UpdatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}

// Delete deletes a dashboard
func (dh *DashboardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	dh.mu.Lock()
	defer dh.mu.Unlock()

	if _, ok := dh.dashboards[id]; !ok {
		http.Error(w, "Dashboard not found", http.StatusNotFound)
		return
	}

	delete(dh.dashboards, id)
	w.WriteHeader(http.StatusNoContent)
}

// AddPanel adds a panel to a dashboard
func (dh *DashboardHandler) AddPanel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dashboardID := r.URL.Query().Get("dashboardId")
	if dashboardID == "" {
		http.Error(w, "Dashboard ID required", http.StatusBadRequest)
		return
	}

	var panel Panel
	if err := json.NewDecoder(r.Body).Decode(&panel); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dh.mu.Lock()
	defer dh.mu.Unlock()

	dashboard, ok := dh.dashboards[dashboardID]
	if !ok {
		http.Error(w, "Dashboard not found", http.StatusNotFound)
		return
	}

	panel.ID = fmt.Sprintf("panel-%d", time.Now().UnixNano())
	dashboard.Panels = append(dashboard.Panels, panel)
	dashboard.UpdatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(panel)
}

// RemovePanel removes a panel from a dashboard
func (dh *DashboardHandler) RemovePanel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dashboardID := r.URL.Query().Get("dashboardId")
	panelID := r.URL.Query().Get("panelId")
	if dashboardID == "" || panelID == "" {
		http.Error(w, "Dashboard ID and Panel ID required", http.StatusBadRequest)
		return
	}

	dh.mu.Lock()
	defer dh.mu.Unlock()

	dashboard, ok := dh.dashboards[dashboardID]
	if !ok {
		http.Error(w, "Dashboard not found", http.StatusNotFound)
		return
	}

	for i, panel := range dashboard.Panels {
		if panel.ID == panelID {
			dashboard.Panels = append(dashboard.Panels[:i], dashboard.Panels[i+1:]...)
			dashboard.UpdatedAt = time.Now()
			break
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdatePanel updates a panel
func (dh *DashboardHandler) UpdatePanel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dashboardID := r.URL.Query().Get("dashboardId")
	panelID := r.URL.Query().Get("panelId")
	if dashboardID == "" || panelID == "" {
		http.Error(w, "Dashboard ID and Panel ID required", http.StatusBadRequest)
		return
	}

	var updates Panel
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dh.mu.Lock()
	defer dh.mu.Unlock()

	dashboard, ok := dh.dashboards[dashboardID]
	if !ok {
		http.Error(w, "Dashboard not found", http.StatusNotFound)
		return
	}

	for i, panel := range dashboard.Panels {
		if panel.ID == panelID {
			if updates.Title != "" {
				dashboard.Panels[i].Title = updates.Title
			}
			if updates.Type != "" {
				dashboard.Panels[i].Type = updates.Type
			}
			if updates.Query != "" {
				dashboard.Panels[i].Query = updates.Query
			}
			if updates.Width > 0 {
				dashboard.Panels[i].Width = updates.Width
			}
			if updates.Height > 0 {
				dashboard.Panels[i].Height = updates.Height
			}
			if updates.X >= 0 {
				dashboard.Panels[i].X = updates.X
			}
			if updates.Y >= 0 {
				dashboard.Panels[i].Y = updates.Y
			}
			dashboard.UpdatedAt = time.Now()
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// ServeHTTP serves dashboard endpoint
func (dh *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/api/v1/dashboards" || path == "/api/v1/dashboards/":
		switch r.Method {
		case http.MethodGet:
			dh.List(w, r)
		case http.MethodPost:
			dh.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case strings.HasPrefix(path, "/api/v1/dashboards/"):
		id := strings.TrimPrefix(path, "/api/v1/dashboards/")
		switch r.Method {
		case http.MethodGet:
			r.URL.RawQuery = "id=" + id
			dh.Get(w, r)
		case http.MethodPut:
			r.URL.RawQuery = "id=" + id
			dh.Update(w, r)
		case http.MethodDelete:
			r.URL.RawQuery = "id=" + id
			dh.Delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.NotFound(w, r)
	}
}
