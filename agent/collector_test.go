package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestCollector_StartStop(t *testing.T) {
	collector := NewCollector(1000)
	collector.Start(context.Background())
	
	time.Sleep(10 * time.Millisecond)
	collector.Stop()
}

func TestCollector_ProcessEvent(t *testing.T) {
	collector := NewCollector(1000)
	collector.Start(context.Background())
	defer collector.Stop()
	
	event := Event{
		Timestamp: time.Now().Unix(),
		Source:    "test-source",
		Metrics:   map[string]float64{"cpu": 50.0},
		Tags:      map[string]string{"host": "test-host"},
	}
	
	select {
	case collector.events <- event:
	default:
		t.Fatal("Failed to send event to collector")
	}
	
	time.Sleep(10 * time.Millisecond)
}

func TestCollector_Backpressure(t *testing.T) {
	collector := NewCollector(10)
	collector.Start(context.Background())
	defer collector.Stop()
	
	var wg sync.WaitGroup
	dropped := 0
	
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			event := Event{
				Timestamp: time.Now().Unix(),
				Source:    "test-source",
				Metrics:   map[string]float64{"cpu": 50.0},
				Tags:      map[string]string{"host": "test-host"},
			}
			
			select {
			case collector.events <- event:
			default:
				dropped++
			}
		}(i)
	}
	
	wg.Wait()
	
	if dropped == 0 {
		t.Log("No events dropped (might be OK depending on timing)")
	} else {
		t.Logf("Dropped %d events due to backpressure", dropped)
	}
}

func TestCollector_MultipleEvents(t *testing.T) {
	collector := NewCollector(1000)
	collector.Start(context.Background())
	defer collector.Stop()
	
	var wg sync.WaitGroup
	eventCount := 50
	
	for i := 0; i < eventCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			event := Event{
				Timestamp: time.Now().Unix(),
				Source:    "test-source",
				Metrics:   map[string]float64{"cpu": 50.0},
				Tags:      map[string]string{"host": "test-host"},
			}
			collector.events <- event
		}(i)
	}
	
	wg.Wait()
	time.Sleep(50 * time.Millisecond)
}

func TestCollector_ConcurrentAccess(t *testing.T) {
	collector := NewCollector(1000)
	collector.Start(context.Background())
	defer collector.Stop()
	
	var wg sync.WaitGroup
	goroutines := 10
	eventsPerGoroutine := 20
	
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < eventsPerGoroutine; i++ {
				event := Event{
					Timestamp: time.Now().Unix(),
					Source:    "test-source",
					Metrics:   map[string]float64{"cpu": 50.0},
					Tags:      map[string]string{"host": "test-host"},
				}
				collector.events <- event
			}
		}(g)
	}
	
	wg.Wait()
	time.Sleep(50 * time.Millisecond)
}
