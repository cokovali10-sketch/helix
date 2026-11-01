// Package agent contains tests for the agent
package agent

import (
	"context"
	"testing"
	"time"
)

func TestAgent_Create(t *testing.T) {
	config := DefaultAgentConfig()
	agent := NewAgent(config)
	if agent == nil {
		t.Fatal("Agent is nil")
	}
}

func TestAgent_StartStop(t *testing.T) {
	config := DefaultAgentConfig()
	agent := NewAgent(config)
	ctx := context.Background()

	if err := agent.Start(ctx); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}

	if !agent.IsRunning() {
		t.Error("Agent should be running")
	}

	if err := agent.Stop(); err != nil {
		t.Fatalf("Failed to stop: %v", err)
	}

	if agent.IsRunning() {
		t.Error("Agent should not be running after stop")
	}
}

func TestAgent_ProcessEvent(t *testing.T) {
	config := DefaultAgentConfig()
	agent := NewAgent(config)
	ctx := context.Background()

	agent.Start(ctx)
	defer agent.Stop()

	event := Event{
		Timestamp: time.Now(),
		Source:    "test-source",
		Metrics:   map[string]float64{"cpu": 50.0},
		Tags:      map[string]string{"host": "test-host"},
		Tenant:    "test-tenant",
	}

	agent.processEvent(event)
}
func TestAgent_AddProcessor(t *testing.T) {
	config := DefaultAgentConfig()
	agent := NewAgent(config)

	processor := &ValidationProcessor{}
	agent.AddProcessor(processor)

	if len(agent.processors) == 0 {
		t.Error("Processor not added")
	}
}

func TestAgent_AddSource(t *testing.T) {
	config := DefaultAgentConfig()
	agent := NewAgent(config)

	source := NewFileSource("/tmp/test.log", 5*time.Second)
	agent.AddSource(source)

	if len(agent.sources) == 0 {
		t.Error("Source not added")
	}
}

func TestAgent_AddExporter(t *testing.T) {
	config := DefaultAgentConfig()
	agent := NewAgent(config)

	exporter := NewHTTPExporter("http://localhost:8080", 100)
	agent.AddExporter(exporter)

	if len(agent.exporters) == 0 {
		t.Error("Exporter not added")
	}
}

func TestAgent_GetMetrics(t *testing.T) {
	config := DefaultAgentConfig()
	agent := NewAgent(config)

	metrics := agent.GetMetrics()
	if metrics == nil {
		t.Fatal("Metrics is nil")
	}
}

func TestValidationProcessor(t *testing.T) {
	processor := &ValidationProcessor{}
	event := &Event{
		Source:  "",
		Metrics: nil,
		Tags:    nil,
	}

	result, err := processor.Process(context.Background(), event)
	if err != nil {
		t.Fatalf("Failed to process: %v", err)
	}
	if result.Source != "unknown" {
		t.Errorf("Expected unknown, got %s", result.Source)
	}
	if result.Metrics == nil {
		t.Error("Metrics should be initialized")
	}
	if result.Tags == nil {
		t.Error("Tags should be initialized")
	}
}

func TestEnrichmentProcessor(t *testing.T) {
	processor := &EnrichmentProcessor{hostname: "test-host", region: "us-east"}
	event := &Event{
		Source: "test",
		Tags:   make(map[string]string),
	}

	result, err := processor.Process(context.Background(), event)
	if err != nil {
		t.Fatalf("Failed to process: %v", err)
	}
	if result.Tags["host"] != "test-host" {
		t.Errorf("Expected test-host, got %s", result.Tags["host"])
	}
	if result.Tags["region"] != "us-east" {
		t.Errorf("Expected us-east, got %s", result.Tags["region"])
	}
}

func TestFilterProcessor(t *testing.T) {
	processor := &FilterProcessor{
		rules: []FilterRule{
			{Field: "source", Operator: "equals", Value: "blocked", Action: "drop"},
		},
	}

	event := &Event{Source: "allowed"}
	result, err := processor.Process(context.Background(), event)
	if err != nil {
		t.Fatalf("Failed to process: %v", err)
	}
	if result == nil {
		t.Error("Event should not be dropped")
	}

	blockedEvent := &Event{Source: "blocked"}
	_, err = processor.Process(context.Background(), blockedEvent)
	if err == nil {
		t.Error("Expected error for blocked event")
	}
}

func TestFileSource_Create(t *testing.T) {
	source := NewFileSource("/tmp/test.log", 5*time.Second)
	if source == nil {
		t.Fatal("Source is nil")
	}
	if source.Name() != "file" {
		t.Errorf("Expected file, got %s", source.Name())
	}
}

func TestHTTPSource_Create(t *testing.T) {
	source := NewHTTPSource("http://localhost:8080", 5*time.Second)
	if source == nil {
		t.Fatal("Source is nil")
	}
	if source.Name() != "http" {
		t.Errorf("Expected http, got %s", source.Name())
	}
}

func TestHTTPExporter_Create(t *testing.T) {
	exporter := NewHTTPExporter("http://localhost:8080", 100)
	if exporter == nil {
		t.Fatal("Exporter is nil")
	}
	if exporter.Name() != "http" {
		t.Errorf("Expected http, got %s", exporter.Name())
	}
}

func TestKafkaExporter_Create(t *testing.T) {
	exporter := NewKafkaExporter([]string{"localhost:9092"}, "test-topic")
	if exporter == nil {
		t.Fatal("Exporter is nil")
	}
	if exporter.Name() != "kafka" {
		t.Errorf("Expected kafka, got %s", exporter.Name())
	}
}

func TestNATSExporter_Create(t *testing.T) {
	exporter := NewNATSExporter("nats://localhost:4222", "test-subject")
	if exporter == nil {
		t.Fatal("Exporter is nil")
	}
	if exporter.Name() != "nats" {
		t.Errorf("Expected nats, got %s", exporter.Name())
	}
}

func TestPluginRegistry(t *testing.T) {
	RegisterPlugin("test", func() Processor { return &ValidationProcessor{} })
	
	plugin, ok := GetPlugin("test")
	if !ok {
		t.Fatal("Plugin not found")
	}
	if plugin == nil {
		t.Fatal("Plugin is nil")
	}

	plugins := ListPlugins()
	if len(plugins) == 0 {
		t.Error("No plugins listed")
	}
}

func TestDefaultAgentConfig(t *testing.T) {
	config := DefaultAgentConfig()
	if config == nil {
		t.Fatal("Config is nil")
	}
	if config.MaxEventsPerSec != 10000 {
		t.Errorf("Expected 10000, got %d", config.MaxEventsPerSec)
	}
}

func TestCreateDefaultAgent(t *testing.T) {
	agent := CreateDefaultAgent()
	if agent == nil {
		t.Fatal("Agent is nil")
	}
	if !agent.IsRunning() {
		// Agent not started yet
	}
}










