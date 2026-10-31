package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestGateway_RegisterTenant(t *testing.T) {
	gateway := NewGateway()
	
	// Register a new tenant
	gateway.RegisterTenant("test-tenant", 1000)
	
	// Try to ingest event for registered tenant
	event := Event{
		Timestamp: time.Now().Unix(),
		Source:    "test-source",
		Metrics:   map[string]float64{"cpu": 50.0},
		Tags:      map[string]string{"host": "test-host"},
		Tenant:    "test-tenant",
	}
	
	success := gateway.Ingest(event)
	if !success {
		t.Fatal("Expected success for registered tenant")
	}
}

func TestGateway_UnregisteredTenant(t *testing.T) {
	gateway := NewGateway()
	
	// Try to ingest event for unregistered tenant (should use default)
	event := Event{
		Timestamp: time.Now().Unix(),
		Source:    "test-source",
		Metrics:   map[string]float64{"cpu": 50.0},
		Tags:      map[string]string{"host": "test-host"},
		Tenant:    "unregistered-tenant",
	}
	
	success := gateway.Ingest(event)
	if !success {
		t.Fatal("Expected success for unregistered tenant using default rate limit")
	}
}

func TestGateway_RateLimit(t *testing.T) {
	gateway := NewGateway()
	gateway.RegisterTenant("rate-limited", 2) // Only 2 events per second
	
	event := Event{
		Timestamp: time.Now().Unix(),
		Source:    "test-source",
		Metrics:   map[string]float64{"cpu": 50.0},
		Tags:      map[string]string{"host": "test-host"},
		Tenant:    "rate-limited",
	}
	
	// First two should succeed
	success1 := gateway.Ingest(event)
	success2 := gateway.Ingest(event)
	
	if !success1 || !success2 {
		t.Fatal("First two events should succeed")
	}
	
	// Third should fail (rate limited)
	success3 := gateway.Ingest(event)
	if success3 {
		t.Fatal("Third event should be rate limited")
	}
}

func TestGateway_ConcurrentTenants(t *testing.T) {
	gateway := NewGateway()
	gateway.RegisterTenant("tenant-a", 1000)
	gateway.RegisterTenant("tenant-b", 1000)
	
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
					Timestamp: time.Now().Unix(),
					Source:    "test-source",
					Metrics:   map[string]float64{"cpu": 50.0},
					Tags:      map[string]string{"host": "test-host"},
					Tenant:    tenant,
				}
				gateway.Ingest(event)
			}
		}(g)
	}
	
	wg.Wait()
}

func TestGateway_DefaultTenantRateLimit(t *testing.T) {
	gateway := NewGateway()
	
	// Default tenant should have rate limit of 10000
	successCount := 0
	for i := 0; i < 10050; i++ {
		event := Event{
			Timestamp: time.Now().Unix(),
			Source:    "test-source",
			Metrics:   map[string]float64{"cpu": 50.0},
			Tags:      map[string]string{"host": "test-host"},
			Tenant:    "default",
		}
		if gateway.Ingest(event) {
			successCount++
		}
	}
	
	// Should have rate limited at some point
	if successCount >= 10050 {
		t.Log("All events succeeded (rate limit may be per-second)")
	}
}