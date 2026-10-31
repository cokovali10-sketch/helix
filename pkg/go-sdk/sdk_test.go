package helix

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_IngestEvent(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	client := NewClient(server.URL, "test-key")
	
	event := &Event{
		Timestamp: time.Now().Unix(),
		Source:    "test-source",
		Metrics:   map[string]float64{"cpu": 50.0},
		Tags:      map[string]string{"host": "test-host"},
		Tenant:    "test-tenant",
	}
	
	err := client.IngestEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestClient_Query(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":[{"metric":{"host":"test-host"},"value":50.0,"timestamp":1700000000}]}`))
	}))
	defer server.Close()
	
	client := NewClient(server.URL, "test-key")
	
	results, err := client.Query(context.Background(), &QueryOptions{
		Expression: "rate(cpu_usage[30s])",
		Tenant:     "test-tenant",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}
}

func TestClient_HealthCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	client := NewClient(server.URL, "test-key")
	
	healthy, err := client.HealthCheck(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if !healthy {
		t.Fatal("Expected healthy")
	}
}

func TestClient_Metrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"events_received":1000,"events_processed":950,"active_tenants":2,"memory_usage_bytes":1024000}`))
	}))
	defer server.Close()
	
	client := NewClient(server.URL, "test-key")
	
	metrics, err := client.Metrics(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if metrics["events_received"] != 1000 {
		t.Errorf("Expected events_received=1000, got %d", metrics["events_received"])
	}
}

func TestClient_ListTenants(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"tenants":[{"id":"tenant-1","name":"Tenant 1","rate_limit":1000,"quota":10000}]}`))
	}))
	defer server.Close()
	
	client := NewClient(server.URL, "test-key")
	
	tenants, err := client.ListTenants(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if len(tenants) != 1 {
		t.Fatalf("Expected 1 tenant, got %d", len(tenants))
	}
}

func TestClient_SearchEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"events":[{"id":"1","timestamp":1700000000,"source":"test","metrics":{"cpu":50.0},"tags":{"host":"test"},"tenant":"test"}]}`))
	}))
	defer server.Close()
	
	client := NewClient(server.URL, "test-key")
	
	events, err := client.SearchEvents(context.Background(), "cpu > 50", map[string]string{"host": "test"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}
}