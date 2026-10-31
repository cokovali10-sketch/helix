// Package catalogapi contains tests for the catalog API
package catalogapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"helix.io/catalog"
)

func TestHandler_ListMetrics(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	// Register a metric
	cat.RegisterMetric(&catalog.MetricMetadata{
		Metadata: catalog.Metadata{
			ID:          "metric-1",
			Type:        "gauge",
			Name:        "cpu_usage",
			Description: "CPU usage percentage",
		},
		MetricName: "cpu_usage",
		Type:       "gauge",
		Unit:       "percent",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/metrics", nil)
	w := httptest.NewRecorder()

	handler.listMetrics(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if response["total"].(float64) != 1 {
		t.Errorf("Expected 1 metric, got %v", response["total"])
	}
}

func TestHandler_GetMetric(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	cat.RegisterMetric(&catalog.MetricMetadata{
		Metadata: catalog.Metadata{
			ID:   "metric-1",
			Type: "gauge",
			Name: "cpu_usage",
		},
		MetricName: "cpu_usage",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/metrics/metric-1", nil)
	w := httptest.NewRecorder()

	handler.getMetric(w, req, "metric-1")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandler_GetMetric_NotFound(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/metrics/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.getMetric(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestHandler_CreateMetric(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	body := `{"id":"metric-2","type":"gauge","name":"memory_usage","metricName":"memory_usage","unit":"bytes"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/catalog/metrics", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createMetric(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestHandler_ListDashboards(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	cat.RegisterDashboard(&catalog.DashboardMetadata{
		Metadata: catalog.Metadata{
			ID:   "dashboard-1",
			Type: "dashboard",
			Name: "Test Dashboard",
		},
		Title: "Test Dashboard",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/dashboards", nil)
	w := httptest.NewRecorder()

	handler.listDashboards(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandler_ListAlerts(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	cat.RegisterAlert(&catalog.AlertMetadata{
		Metadata: catalog.Metadata{
			ID:   "alert-1",
			Type: "alert",
			Name: "High CPU",
		},
		Severity: "critical",
		Status:   "firing",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/alerts", nil)
	w := httptest.NewRecorder()

	handler.listAlerts(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandler_ListTenants(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	cat.RegisterTenant(&catalog.TenantMetadata{
		Metadata: catalog.Metadata{
			ID:   "tenant-1",
			Type: "tenant",
			Name: "Test Tenant",
		},
		Active: true,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/tenants", nil)
	w := httptest.NewRecorder()

	handler.listTenants(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandler_Search(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	cat.RegisterMetric(&catalog.MetricMetadata{
		Metadata: catalog.Metadata{
			ID:          "metric-1",
			Type:        "gauge",
			Name:        "cpu_usage",
			Description: "CPU usage percentage",
		},
		MetricName: "cpu_usage",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/search?q=cpu", nil)
	w := httptest.NewRecorder()

	handler.handleSearch(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandler_Search_MissingQuery(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/search", nil)
	w := httptest.NewRecorder()

	handler.handleSearch(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandler_ListRoles(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/roles", nil)
	w := httptest.NewRecorder()

	handler.listRoles(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	// Should have default roles
	if response["total"].(float64) < 4 {
		t.Errorf("Expected at least 4 roles, got %v", response["total"])
	}
}

func TestHandler_GetRole(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/roles/admin", nil)
	w := httptest.NewRecorder()

	handler.getRole(w, req, "admin")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandler_GetRole_NotFound(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/roles/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.getRole(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestHandler_ListUsers(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/users", nil)
	w := httptest.NewRecorder()

	handler.listUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandler_CreateUser(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	body := `{"id":"user-1","username":"testuser","email":"test@example.com","roles":["viewer"],"tenantId":"tenant-1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/catalog/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createUser(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestHandler_CheckAccess(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	// Create user and assign role
	rbac.CreateUser(&catalog.User{ID: "user-1", Username: "test", Email: "test@test.com"})
	rbac.AssignRole("user-1", "admin")

	body := `{"userId":"user-1","resource":"metrics","action":"read"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/catalog/check-access", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.handleCheckAccess(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if response["allowed"] != true {
		t.Errorf("Expected allowed true, got %v", response["allowed"])
	}
}

func TestHandler_Permissions(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/permissions", nil)
	w := httptest.NewRecorder()

	handler.handlePermissions(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if response["total"].(float64) < 16 {
		t.Errorf("Expected at least 16 permissions, got %v", response["total"])
	}
}

func TestHandler_CreateRole(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	body := `{"id":"custom-role","name":"Custom Role","description":"Custom","permissions":["metrics:read"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/catalog/roles", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createRole(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestHandler_DeleteRole(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	// Create a role first
	rbac.CreateRole(&catalog.Role{ID: "test-role", Name: "Test", Description: "Test", Permissions: []string{}})

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/catalog/roles/test-role", nil)
	w := httptest.NewRecorder()

	handler.deleteRole(w, req, "test-role")

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestHandler_DeleteRole_NotFound(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/catalog/roles/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.deleteRole(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestHandler_DeleteUser(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	rbac.CreateUser(&catalog.User{ID: "user-1", Username: "test", Email: "test@test.com"})

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/catalog/users/user-1", nil)
	w := httptest.NewRecorder()

	handler.deleteUser(w, req, "user-1")

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestHandler_GetUser(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	rbac.CreateUser(&catalog.User{ID: "user-1", Username: "test", Email: "test@test.com"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/users/user-1", nil)
	w := httptest.NewRecorder()

	handler.getUser(w, req, "user-1")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandler_GetUser_NotFound(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/users/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.getUser(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestHandler_UpdateTenant(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	cat.RegisterTenant(&catalog.TenantMetadata{
		Metadata: catalog.Metadata{ID: "tenant-1", Type: "tenant", Name: "Test"},
		Active:   true,
	})

	body := `{"id":"tenant-1","name":"Updated Tenant","active":true}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/catalog/tenants/tenant-1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.updateTenant(w, req, "tenant-1")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandler_DeleteTenant(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/catalog/tenants/tenant-1", nil)
	w := httptest.NewRecorder()

	handler.deleteTenant(w, req, "tenant-1")

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestHandler_UpdateDashboard(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	cat.RegisterDashboard(&catalog.DashboardMetadata{
		Metadata: catalog.Metadata{ID: "dashboard-1", Type: "dashboard", Name: "Test"},
		Title:    "Test Dashboard",
	})

	body := `{"id":"dashboard-1","title":"Updated Dashboard"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/catalog/dashboards/dashboard-1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.updateDashboard(w, req, "dashboard-1")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandler_DeleteDashboard(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/catalog/dashboards/dashboard-1", nil)
	w := httptest.NewRecorder()

	handler.deleteDashboard(w, req, "dashboard-1")

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestHandler_UpdateAlert(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	cat.RegisterAlert(&catalog.AlertMetadata{
		Metadata: catalog.Metadata{ID: "alert-1", Type: "alert", Name: "Test"},
		Severity: "warning",
		Status:   "firing",
	})

	body := `{"id":"alert-1","severity":"critical","status":"resolved"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/catalog/alerts/alert-1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.updateAlert(w, req, "alert-1")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandler_DeleteAlert(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/catalog/alerts/alert-1", nil)
	w := httptest.NewRecorder()

	handler.deleteAlert(w, req, "alert-1")

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestHandler_UpdateMetric(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	cat.RegisterMetric(&catalog.MetricMetadata{
		Metadata: catalog.Metadata{ID: "metric-1", Type: "gauge", Name: "cpu_usage"},
		MetricName: "cpu_usage",
	})

	body := `{"id":"metric-1","name":"updated_cpu","type":"gauge"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/catalog/metrics/metric-1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.updateMetric(w, req, "metric-1")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandler_DeleteMetric(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/catalog/metrics/metric-1", nil)
	w := httptest.NewRecorder()

	handler.deleteMetric(w, req, "metric-1")

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/catalog/metrics", nil)
	w := httptest.NewRecorder()

	handler.handleMetrics(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandler_RegisterRoutes(t *testing.T) {
	store := catalog.NewInMemoryMetadataStore()
	cat := catalog.NewCatalog(store)
	rbac := catalog.NewRBAC()
	handler := NewHandler(cat, rbac)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Just verify it doesn't panic
}