// Package ingestion contains tests for the ingestion gateway
package ingestion

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestIngestionGateway_Create(t *testing.T) {
	gw := NewIngestionGateway("nats://localhost:4222", []string{"localhost:9092"})
	if gw == nil {
		t.Fatal("Gateway is nil")
	}
	if gw.running {
		t.Error("Gateway should not be running initially")
	}
}

func TestIngestionGateway_StartStop(t *testing.T) {
	gw := NewIngestionGateway("nats://localhost:4222", []string{"localhost:9092"})
	ctx := context.Background()
	
	if err := gw.Start(ctx); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}
	
	if !gw.IsRunning() {
		t.Error("Gateway should be running")
	}
	
	gw.Stop()
	
	if gw.IsRunning() {
		t.Error("Gateway should not be running after stop")
	}
}

func TestIngestionGateway_RegisterTenant(t *testing.T) {
	gw := NewIngestionGateway("nats://localhost:4222", []string{"localhost:9092"})
	
	tenant := &Tenant{
		ID:        "tenant-1",
		Name:      "Test Tenant",
		RateLimit: 1000,
		Quota:     10000,
		Active:    true,
	}
	
	gw.RegisterTenant(tenant)
	
	got, ok := gw.GetTenant("tenant-1")
	if !ok {
		t.Fatal("Tenant not found")
	}
	if got.Name != "Test Tenant" {
		t.Errorf("Expected Test Tenant, got %s", got.Name)
	}
}

func TestIngestionGateway_Ingest(t *testing.T) {
	gw := NewIngestionGateway("nats://localhost:4222", []string{"localhost:9092"})
	
	event := Event{
		Timestamp: time.Now(),
		Source:    "test-source",
		Metrics:   map[string]float64{"cpu": 50.0},
		Tags:      map[string]string{"host": "test-host"},
		Tenant:    "default",
	}
	
	success := gw.Ingest(event)
	if !success {
		t.Error("Expected ingest to succeed")
	}
}

func TestIngestionGateway_RateLimit(t *testing.T) {
	gw := NewIngestionGateway("nats://localhost:4222", []string{"localhost:9092"})
	
	tenant := &Tenant{
		ID:        "rate-limited",
		Name:      "Rate Limited Tenant",
		RateLimit: 2,
		Quota:     10000,
		Active:    true,
	}
	gw.RegisterTenant(tenant)
	
	event := Event{
		Timestamp: time.Now(),
		Source:    "test-source",
		Metrics:   map[string]float64{"cpu": 50.0},
		Tags:      map[string]string{"host": "test-host"},
		Tenant:    "rate-limited",
	}
	
	// First two should succeed
	success1 := gw.Ingest(event)
	success2 := gw.Ingest(event)
	
	if !success1 || !success2 {
		t.Fatal("First two events should succeed")
	}
	
	// Third should fail
	success3 := gw.Ingest(event)
	if success3 {
		t.Fatal("Third event should be rate limited")
	}
}

func TestIngestionGateway_ConcurrentTenants(t *testing.T) {
	gw := NewIngestionGateway("nats://localhost:4222", []string{"localhost:9092"})
	
	gw.RegisterTenant(&Tenant{ID: "tenant-a", Name: "Tenant A", RateLimit: 1000, Quota: 10000, Active: true})
	gw.RegisterTenant(&Tenant{ID: "tenant-b", Name: "Tenant B", RateLimit: 1000, Quota: 10000, Active: true})
	
	var wg sync.WaitGroup
	goroutines := 10
	eventsPerGoroutine := 50
	
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			tenant := "tenant-a"
			if g%2 == 1 {
				tenant = "tenant-b"
			}
			
			for i := 0; i < eventsPerGoroutine; i++ {
				event := Event{
					Timestamp: time.Now(),
					Source:    "test-source",
					Metrics:   map[string]float64{"cpu": 50.0},
					Tags:      map[string]string{"host": "test-host"},
					Tenant:    tenant,
				}
				gw.Ingest(event)
			}
		}(g)
	}
	
	wg.Wait()
}

func TestIngestionGateway_ListTenants(t *testing.T) {
	gw := NewIngestionGateway("nats://localhost:4222", []string{"localhost:9092"})
	
	gw.RegisterTenant(&Tenant{ID: "tenant-1", Name: "Tenant 1", RateLimit: 1000, Quota: 10000, Active: true})
	gw.RegisterTenant(&Tenant{ID: "tenant-2", Name: "Tenant 2", RateLimit: 1000, Quota: 10000, Active: true})
	
	tenants := gw.ListTenants()
	if len(tenants) != 2 {
		t.Errorf("Expected 2 tenants, got %d", len(tenants))
	}
}

func TestIngestionGateway_UpdateTenant(t *testing.T) {
	gw := NewIngestionGateway("nats://localhost:4222", []string{"localhost:9092"})
	
	gw.RegisterTenant(&Tenant{ID: "tenant-1", Name: "Tenant 1", RateLimit: 1000, Quota: 10000, Active: true})
	
	gw.UpdateTenant(&Tenant{ID: "tenant-1", Name: "Updated Tenant", RateLimit: 2000, Quota: 20000, Active: true})
	
	tenant, ok := gw.GetTenant("tenant-1")
	if !ok {
		t.Fatal("Tenant not found")
	}
	if tenant.Name != "Updated Tenant" {
		t.Errorf("Expected Updated Tenant, got %s", tenant.Name)
	}
	if tenant.RateLimit != 2000 {
		t.Errorf("Expected rate limit 2000, got %d", tenant.RateLimit)
	}
}

func TestIngestionGateway_DeleteTenant(t *testing.T) {
	gw := NewIngestionGateway("nats://localhost:4222", []string{"localhost:9092"})
	
	gw.RegisterTenant(&Tenant{ID: "tenant-1", Name: "Tenant 1", RateLimit: 1000, Quota: 10000, Active: true})
	gw.DeleteTenant("tenant-1")
	
	_, ok := gw.GetTenant("tenant-1")
	if ok {
		t.Fatal("Tenant should be deleted")
	}
}

func TestTokenBucket_Allow(t *testing.T) {
	tb := NewTokenBucket(10, 5)
	
	// Should allow up to capacity
	for i := 0; i < 10; i++ {
		if !tb.Allow() {
			t.Errorf("Expected allow on iteration %d", i)
		}
	}
	
	// Should deny when empty
	if tb.Allow() {
		t.Error("Expected deny when bucket empty")
	}
}

func TestTokenBucket_Refill(t *testing.T) {
	tb := NewTokenBucket(10, 100)
	
	// Drain bucket
	for i := 0; i < 10; i++ {
		tb.Allow()
	}
	
	// Wait for refill
	time.Sleep(100 * time.Millisecond)
	
	// Should have refilled
	if !tb.Allow() {
		t.Error("Expected allow after refill")
	}
}

func TestEventValidator_Valid(t *testing.T) {
	validator := NewEventValidator()
	
	event := Event{
		Timestamp: time.Now(),
		Source:    "test-source",
		Tenant:    "test-tenant",
	}
	
	if err := validator.Validate(event); err != nil {
		t.Errorf("Expected valid, got %v", err)
	}
}

func TestEventValidator_MissingSource(t *testing.T) {
	validator := NewEventValidator()
	
	event := Event{
		Timestamp: time.Now(),
		Tenant:    "test-tenant",
	}
	
	if err := validator.Validate(event); err == nil {
		t.Error("Expected error for missing source")
	}
}

func TestEventValidator_MissingTenant(t *testing.T) {
	validator := NewEventValidator()
	
	event := Event{
		Timestamp: time.Now(),
		Source:    "test-source",
	}
	
	if err := validator.Validate(event); err == nil {
		t.Error("Expected error for missing tenant")
	}
}

func TestBatchProcessor_Add(t *testing.T) {
	bp := NewBatchProcessor(10, 1*time.Second)
	
	event := Event{
		Timestamp: time.Now(),
		Source:    "test-source",
		Tenant:    "test-tenant",
	}
	
	bp.Add(event)
	// Just verify it doesn't panic
}

func TestNATSClient_Create(t *testing.T) {
	client := NewNATSClient("nats://localhost:4222")
	if client == nil {
		t.Fatal("Client is nil")
	}
}

func TestKafkaClient_Create(t *testing.T) {
	client := NewKafkaClient([]string{"localhost:9092"})
	if client == nil {
		t.Fatal("Client is nil")
	}
}

func TestGatewayMetrics(t *testing.T) {
	metrics := &GatewayMetrics{}
	metrics.mu.Lock()
	metrics.EventsReceived = 100
	metrics.EventsProcessed = 95
	metrics.EventsDropped = 5
	metrics.mu.Unlock()
	
	if metrics.EventsReceived != 100 {
		t.Errorf("Expected 100, got %d", metrics.EventsReceived)
	}
}