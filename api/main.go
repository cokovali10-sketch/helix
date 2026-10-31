// Package api provides gRPC and REST API definitions for the Helix observability platform
// This module defines the API contracts for the platform's services
package main

import (
	"fmt"
	"net/http"
)

// API represents the Helix platform API server
type API struct {
	version    string
	host       string
	port       string
	grpcPort   string
	routes     map[string]http.Handler
}

// NewAPI creates a new API server
func NewAPI(version, host, grpcPort, httpPort string) *API {
	return &API{
		version:   version,
		host:      host,
		grpcPort:  grpcPort,
		port:      httpPort,
		routes:    make(map[string]http.Handler),
	}
}

// RegisterRoute registers an HTTP route handler
func (a *API) RegisterRoute(pattern string, handler http.Handler) {
	a.routes[pattern] = handler
	fmt.Printf "API route registered: %s\n", pattern)
}

// Start starts the API server
func (a *API) Start() error {
	fmt.Printf("Starting Helix API v%s on %s:%s\n", a.version, a.host, a.port)
	fmt.Printf("gRPC server on %s:%s\n", a.host, a.grpcPort)
	
	// TODO: Start HTTP and gRPC servers
	
	return nil
}

// RegisterService registers a service with the API
func (a *API) RegisterService(name string, addr string) {
	fmt.Printf("Registering service: %s at %s\n", name, addr)
	// TODO: Implement service registration
}

// ---------------------------------------------------------------------------
// gRPC Service Definitions
// ---------------------------------------------------------------------------

// EventService defines the gRPC service for event ingestion
type EventService interface {
	// Ingest receives events from agents
	Ingest(*EventRequest) (*EventResponse, error)
	// Query executes a query against stored events
	Query(*QueryRequest) (*QueryResponse, error)
}

// EventRequest represents an event ingestion request
type EventRequest struct {
	Timestamp int64
	Source    string
	Metrics   map[string]float64
	Tags      map[string]string
	Tenant    string
}

// EventResponse represents an event ingestion response
type EventResponse struct {
	Success bool
	Message string
	ID      string
}

// QueryRequest represents a query request
type QueryRequest struct {
	Expression string
	TimeRange  []int64
	Tenant     string
}

// QueryResponse represents a query response
type QueryResponse struct {
	Results    []*QueryResult
	Error      string
	ExecutionTime float64
}

// QueryResult represents a single query result
type QueryResult struct {
	Metric     map[string]string
	Value      float64
	Timestamp  int64
}

// ---------------------------------------------------------------------------
// REST API Endpoints
// ---------------------------------------------------------------------------

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status   string `json:"status"`
	Version  string `json:"version"`
	Uptime     string `json:"uptime"`
}

// MetricsResponse represents the metrics response
type MetricsResponse struct {
	EventsReceived   int64 `json:"events_received"`
	EventsProcessed  int64 `json:"events_processed"`
	ActiveTenants    int   `json:"active_tenants"`
	MemoryUsageBytes int64 `json:"memory_usage_bytes"`
}

// SetupRoutes sets up the default API routes
func (a *API) SetupRoutes() {
	// Health check endpoint
	a.RegisterRoute("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		a.writeJSON(w, HealthCheckResponse{
			Status:  "healthy",
			Version: a.version,
		})
	}))
	
	// Metrics endpoint
	a.RegisterRoute("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.writeJSON(w, MetricsResponse{
			EventsReceived:   0,
			EventsProcessed:  0,
			ActiveTenants:    0,
			MemoryUsageBytes: 0,
		})
	}))
	
	// Event ingestion endpoint
	a.RegisterRoute("/v1/events", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Parse and ingest event
		w.WriteHeader(http.StatusOK)
		a.writeJSON(w, EventResponse{
			Success: true,
			Message: "Event received",
		})
	}))
	
	// Query endpoint
	a.RegisterRoute("/v1/query", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Parse query and execute
		w.WriteHeader(http.StatusOK)
		a.writeJSON(w, QueryResponse{
			Results:    nil,
			Error:      "",
			ExecutionTime: 0,
		})
	}))
}

// writeJSON writes a JSON response
func (a *API) writeJSON(w http.ResponseWriter, v interface{}) {
	// TODO: Proper JSON encoding
	w.Header().Set("Content-Type", "application/json")
	_ = v
}

func main() {
	api := NewAPI("0.1.0", "0.0.0.0", "50051", "8080")
	
	// Set up default routes
	api.SetupRoutes()
	
	// Start the API server
	if err := api.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start API: %v\n", err)
		os.Exit(1)
	}
	
	// Keep running
	select {}
}