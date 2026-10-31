// Package api contains tests for the API server
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer_Health(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	
	server.handleHealth(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	
	if response["status"] != "healthy" {
		t.Errorf("Expected status healthy, got %v", response["status"])
	}
}

func TestServer_Metrics(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	
	server.handleMetrics(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestServer_CreateEvent(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	body := `{"source":"test","metrics":{"cpu":50.0},"tags":{"host":"test"},"tenant":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	server.handleEvents(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestServer_ListEvents(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	req := httptest.NewRequest(http.MethodGet, "/v1/events", nil)
	w := httptest.NewRecorder()
	
	server.handleEvents(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestServer_ExecuteQuery(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	req := httptest.NewRequest(http.MethodGet, "/v1/query?query=cpu_usage", nil)
	w := httptest.NewRecorder()
	
	server.handleQuery(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	
	if response["resultType"] != "vector" {
		t.Errorf("Expected resultType vector, got %v", response["resultType"])
	}
}

func TestServer_ListTenants(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	req := httptest.NewRequest(http.MethodGet, "/v1/tenants", nil)
	w := httptest.NewRecorder()
	
	server.handleTenants(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestServer_CreateTenant(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	body := `{"name":"test-tenant","rateLimit":1000,"quota":10000}`
	req := httptest.NewRequest(http.MethodPost, "/v1/tenants", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	server.handleTenants(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestServer_ListDashboards(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	req := httptest.NewRequest(http.MethodGet, "/v1/dashboards", nil)
	w := httptest.NewRecorder()
	
	server.handleDashboards(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestServer_CreateDashboard(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	body := `{"title":"Test Dashboard","description":"Test"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/dashboards", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	server.handleDashboards(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestServer_ListAlerts(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	req := httptest.NewRequest(http.MethodGet, "/v1/alerts", nil)
	w := httptest.NewRecorder()
	
	server.handleAlerts(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestServer_CreateAlert(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	body := `{"name":"Test Alert","expression":"cpu > 90","severity":"warning"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/alerts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	server.handleAlerts(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestServer_MethodNotAllowed(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	req := httptest.NewRequest(http.MethodPut, "/v1/events", nil)
	w := httptest.NewRecorder()
	
	server.handleEvents(w, req)
	
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestServer_CORS(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	req := httptest.NewRequest(http.MethodOptions, "/v1/events", nil)
	w := httptest.NewRecorder()
	
	handler := server.corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	handler.ServeHTTP(w, req)
	
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected CORS header")
	}
}

func TestServer_DefaultConfig(t *testing.T) {
	config := DefaultServerConfig()
	
	if config.HTTPPort != 8080 {
		t.Errorf("Expected HTTP port 8080, got %d", config.HTTPPort)
	}
	if config.GRPCPort != 50051 {
		t.Errorf("Expected gRPC port 50051, got %d", config.GRPCPort)
	}
	if config.Host != "0.0.0.0" {
		t.Errorf("Expected host 0.0.0.0, got %s", config.Host)
	}
}

func TestServer_CreateServer(t *testing.T) {
	server := CreateServer()
	if server == nil {
		t.Fatal("Server is nil")
	}
	if server.config == nil {
		t.Fatal("Config is nil")
	}
}