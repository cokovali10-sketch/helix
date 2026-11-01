// Package catalogapi implements the REST API for the catalog and metadata system
package catalogapi

import (
	"encoding/json"

	"net/http"
	"strconv"
	"strings"

	"helix.io/catalog"
)

// Handler handles catalog API requests
type Handler struct {
	catalog *catalog.Catalog
	rbac    *catalog.RBAC
}

// NewHandler creates a new catalog API handler
func NewHandler(cat *catalog.Catalog, rbac *catalog.RBAC) *Handler {
	return &Handler{
		catalog: cat,
		rbac:    rbac,
	}
}

// RegisterRoutes registers catalog API routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/catalog/metrics", h.handleMetrics)
	mux.HandleFunc("/api/v1/catalog/metrics/", h.handleMetricByID)
	mux.HandleFunc("/api/v1/catalog/dashboards", h.handleDashboards)
	mux.HandleFunc("/api/v1/catalog/dashboards/", h.handleDashboardByID)
	mux.HandleFunc("/api/v1/catalog/alerts", h.handleAlerts)
	mux.HandleFunc("/api/v1/catalog/alerts/", h.handleAlertByID)
	mux.HandleFunc("/api/v1/catalog/tenants", h.handleTenants)
	mux.HandleFunc("/api/v1/catalog/tenants/", h.handleTenantByID)
	mux.HandleFunc("/api/v1/catalog/search", h.handleSearch)
	mux.HandleFunc("/api/v1/catalog/roles", h.handleRoles)
	mux.HandleFunc("/api/v1/catalog/roles/", h.handleRoleByID)
	mux.HandleFunc("/api/v1/catalog/users", h.handleUsers)
	mux.HandleFunc("/api/v1/catalog/users/", h.handleUserByID)
	mux.HandleFunc("/api/v1/catalog/permissions", h.handlePermissions)
	mux.HandleFunc("/api/v1/catalog/check-access", h.handleCheckAccess)
}

// handleMetrics handles metrics endpoints
func (h *Handler) handleMetrics(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listMetrics(w, r)
	case http.MethodPost:
		h.createMetric(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleMetricByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/catalog/metrics/")
	switch r.Method {
	case http.MethodGet:
		h.getMetric(w, r, id)
	case http.MethodPut:
		h.updateMetric(w, r, id)
	case http.MethodDelete:
		h.deleteMetric(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listMetrics(w http.ResponseWriter, r *http.Request) {
	filter := &catalog.MetadataFilter{
		Type:   r.URL.Query().Get("type"),
		Owner:  r.URL.Query().Get("owner"),
		Limit:  atoi(r.URL.Query().Get("limit")),
		Offset: atoi(r.URL.Query().Get("offset")),
	}
	
	metrics := h.catalog.ListMetrics(filter)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"metrics": metrics,
		"total":   len(metrics),
	})
}

func (h *Handler) getMetric(w http.ResponseWriter, r *http.Request, id string) {
	metric, ok := h.catalog.GetMetric(id)
	if !ok {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metric)
}

func (h *Handler) createMetric(w http.ResponseWriter, r *http.Request) {
	var metric catalog.MetricMetadata
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	h.catalog.RegisterMetric(&metric)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(metric)
}

func (h *Handler) updateMetric(w http.ResponseWriter, r *http.Request, id string) {
	var metric catalog.MetricMetadata
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	metric.ID = id
	h.catalog.RegisterMetric(&metric)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metric)
}

func (h *Handler) deleteMetric(w http.ResponseWriter, r *http.Request, id string) {
	w.WriteHeader(http.StatusNoContent)
}

// handleDashboards handles dashboard endpoints
func (h *Handler) handleDashboards(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listDashboards(w, r)
	case http.MethodPost:
		h.createDashboard(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleDashboardByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/catalog/dashboards/")
	switch r.Method {
	case http.MethodGet:
		h.getDashboard(w, r, id)
	case http.MethodPut:
		h.updateDashboard(w, r, id)
	case http.MethodDelete:
		h.deleteDashboard(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listDashboards(w http.ResponseWriter, r *http.Request) {
	filter := &catalog.MetadataFilter{
		Owner:  r.URL.Query().Get("owner"),
		Limit:  atoi(r.URL.Query().Get("limit")),
		Offset: atoi(r.URL.Query().Get("offset")),
	}
	
	dashboards := h.catalog.ListDashboards(filter)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"dashboards": dashboards,
		"total":      len(dashboards),
	})
}

func (h *Handler) getDashboard(w http.ResponseWriter, r *http.Request, id string) {
	dashboard, ok := h.catalog.GetDashboard(id)
	if !ok {
		http.Error(w, "Dashboard not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}

func (h *Handler) createDashboard(w http.ResponseWriter, r *http.Request) {
	var dashboard catalog.DashboardMetadata
	if err := json.NewDecoder(r.Body).Decode(&dashboard); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	h.catalog.RegisterDashboard(&dashboard)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dashboard)
}

func (h *Handler) updateDashboard(w http.ResponseWriter, r *http.Request, id string) {
	var dashboard catalog.DashboardMetadata
	if err := json.NewDecoder(r.Body).Decode(&dashboard); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	dashboard.ID = id
	h.catalog.RegisterDashboard(&dashboard)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}

func (h *Handler) deleteDashboard(w http.ResponseWriter, r *http.Request, id string) {
	w.WriteHeader(http.StatusNoContent)
}

// handleAlerts handles alert endpoints
func (h *Handler) handleAlerts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listAlerts(w, r)
	case http.MethodPost:
		h.createAlert(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleAlertByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/catalog/alerts/")
	switch r.Method {
	case http.MethodGet:
		h.getAlert(w, r, id)
	case http.MethodPut:
		h.updateAlert(w, r, id)
	case http.MethodDelete:
		h.deleteAlert(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listAlerts(w http.ResponseWriter, r *http.Request) {
	filter := &catalog.MetadataFilter{
		Type:   r.URL.Query().Get("status"),
		Limit:  atoi(r.URL.Query().Get("limit")),
		Offset: atoi(r.URL.Query().Get("offset")),
	}
	
	alerts := h.catalog.ListAlerts(filter)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"alerts": alerts,
		"total":  len(alerts),
	})
}

func (h *Handler) getAlert(w http.ResponseWriter, r *http.Request, id string) {
	alert, ok := h.catalog.GetAlert(id)
	if !ok {
		http.Error(w, "Alert not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alert)
}

func (h *Handler) createAlert(w http.ResponseWriter, r *http.Request) {
	var alert catalog.AlertMetadata
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	h.catalog.RegisterAlert(&alert)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(alert)
}

func (h *Handler) updateAlert(w http.ResponseWriter, r *http.Request, id string) {
	var alert catalog.AlertMetadata
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	alert.ID = id
	h.catalog.RegisterAlert(&alert)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alert)
}

func (h *Handler) deleteAlert(w http.ResponseWriter, r *http.Request, id string) {
	w.WriteHeader(http.StatusNoContent)
}

// handleTenants handles tenant endpoints
func (h *Handler) handleTenants(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listTenants(w, r)
	case http.MethodPost:
		h.createTenant(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleTenantByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/catalog/tenants/")
	switch r.Method {
	case http.MethodGet:
		h.getTenant(w, r, id)
	case http.MethodPut:
		h.updateTenant(w, r, id)
	case http.MethodDelete:
		h.deleteTenant(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listTenants(w http.ResponseWriter, r *http.Request) {
	tenants := h.catalog.ListTenants()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tenants": tenants,
		"total":   len(tenants),
	})
}

func (h *Handler) getTenant(w http.ResponseWriter, r *http.Request, id string) {
	tenant, ok := h.catalog.GetTenant(id)
	if !ok {
		http.Error(w, "Tenant not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenant)
}

func (h *Handler) createTenant(w http.ResponseWriter, r *http.Request) {
	var tenant catalog.TenantMetadata
	if err := json.NewDecoder(r.Body).Decode(&tenant); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	h.catalog.RegisterTenant(&tenant)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tenant)
}

func (h *Handler) updateTenant(w http.ResponseWriter, r *http.Request, id string) {
	var tenant catalog.TenantMetadata
	if err := json.NewDecoder(r.Body).Decode(&tenant); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	tenant.ID = id
	h.catalog.RegisterTenant(&tenant)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenant)
}

func (h *Handler) deleteTenant(w http.ResponseWriter, r *http.Request, id string) {
	w.WriteHeader(http.StatusNoContent)
}

// handleSearch handles search
func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query required", http.StatusBadRequest)
		return
	}
	
	results := h.catalog.Search(query)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
		"total":   len(results),
	})
}

// handleRoles handles role endpoints
func (h *Handler) handleRoles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listRoles(w, r)
	case http.MethodPost:
		h.createRole(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleRoleByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/catalog/roles/")
	switch r.Method {
	case http.MethodGet:
		h.getRole(w, r, id)
	case http.MethodDelete:
		h.deleteRole(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listRoles(w http.ResponseWriter, r *http.Request) {
	roles := h.rbac.ListRoles()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"roles": roles,
		"total": len(roles),
	})
}

func (h *Handler) getRole(w http.ResponseWriter, r *http.Request, id string) {
	role, ok := h.rbac.GetRole(id)
	if !ok {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func (h *Handler) createRole(w http.ResponseWriter, r *http.Request) {
	var role catalog.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.rbac.CreateRole(&role); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

func (h *Handler) deleteRole(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.rbac.DeleteRole(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// handleUsers handles user endpoints
func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleUserByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/catalog/users/")
	switch r.Method {
	case http.MethodGet:
		h.getUser(w, r, id)
	case http.MethodDelete:
		h.deleteUser(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	users := h.rbac.ListUsers()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": users,
		"total": len(users),
	})
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request, id string) {
	user, ok := h.rbac.GetUser(id)
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var user catalog.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.rbac.CreateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.rbac.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// handlePermissions handles permission endpoints
func (h *Handler) handlePermissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Return all permissions
	permissions := []map[string]string{
		{"id": "metrics:read", "resource": "metrics", "action": "read"},
		{"id": "metrics:write", "resource": "metrics", "action": "write"},
		{"id": "metrics:delete", "resource": "metrics", "action": "delete"},
		{"id": "dashboards:read", "resource": "dashboards", "action": "read"},
		{"id": "dashboards:write", "resource": "dashboards", "action": "write"},
		{"id": "dashboards:delete", "resource": "dashboards", "action": "delete"},
		{"id": "alerts:read", "resource": "alerts", "action": "read"},
		{"id": "alerts:write", "resource": "alerts", "action": "write"},
		{"id": "alerts:delete", "resource": "alerts", "action": "delete"},
		{"id": "tenants:read", "resource": "tenants", "action": "read"},
		{"id": "tenants:write", "resource": "tenants", "action": "write"},
		{"id": "tenants:delete", "resource": "tenants", "action": "delete"},
		{"id": "users:read", "resource": "users", "action": "read"},
		{"id": "users:write", "resource": "users", "action": "write"},
		{"id": "users:delete", "resource": "users", "action": "delete"},
		{"id": "admin", "resource": "*", "action": "admin"},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"permissions": permissions,
		"total":       len(permissions),
	})
}

// handleCheckAccess handles access check
func (h *Handler) handleCheckAccess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		UserID   string `json:"userId"`
		Resource string `json:"resource"`
		Action   string `json:"action"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	allowed := h.rbac.CheckAccess(req.UserID, req.Resource, req.Action)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"allowed": allowed,
	})
}

func atoi(s string) int {
	if s == "" {
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
