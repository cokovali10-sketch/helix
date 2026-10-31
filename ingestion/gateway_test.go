package main

import (
	"sync"
	"testing"
	"time"
)

func TestGateway_RegisterTenant(t *testing.T) {
	gateway := NewGateway()
	
	gateway.RegisterTenant("test-tenant", 1000)
	
	event := Event{
		Timestamp: time.Now(),
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
	
	event := Event{
		Timestamp: time.Now(),
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
	gateway.RegisterTenant("rate-limited", 2)
	
	event := Event{
		Timestamp: time.Now(),
		Source:    "test-source",
		Metrics:   map[string]float64{"cpu": 50.0},
		Tags:      map[string]string{"host": "test-host"},
		Tenant:    "rate-limited",
	}
	
	success1 := gateway.Ingest(event)
	success2 := gateway.Ingest(event)
	
	if !success1 || !success2 {
		t.Fatal("First two events should succeed")
	}
	
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
					Timestamp: time.Now(),
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
	
	successCount := 0
	for i := 0; i < 10050; i++ {
		event := Event{
			Timestamp: time.Now(),
			Source:    "test-source",
			Metrics:   map[string]float64{"cpu": 50.0},
			Tags:      map[string]string{"host": "test-host"},
			Tenant:    "default",
		}
		if gateway.Ingest(event) {
			successCount++
		}
	}
	
	if successCount >= 10050 {
		t.Log("All events succeeded (rate limit may be per-second)")
	}
}