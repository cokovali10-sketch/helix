// Package api implements the gRPC and REST API for the Helix platform
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server represents the Helix API server
type Server struct {
	httpServer *http.Server
	grpcServer *grpc.Server
	config     *ServerConfig
	mu         sync.RWMutex
	running    bool
	ctx        context.Context
	cancel     context.CancelFunc
}

// ServerConfig holds server configuration
type ServerConfig struct {
	HTTPPort     int
	GRPCPort     int
	Host         string
	TLSEnabled   bool
	CertFile     string
	KeyFile      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	MaxHeaderBytes int
}

// NewServer creates a new API server
func NewServer(config *ServerConfig) *Server {
	return &Server{
		config: config,
	}
}

// Start starts the API server
func (s *Server) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)
	
	// Start HTTP server
	go s.startHTTP()
	
	// Start gRPC server
	go s.startGRPC()
	
	s.mu.Lock()
	s.running = true
	s.mu.Unlock()
	
	return nil
}

// Stop stops the API server
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if !s.running {
		return fmt.Errorf("server not running")
	}
	
	s.running = false
	s.cancel()
	
	if s.httpServer != nil {
		s.httpServer.Shutdown(context.Background())
	}
	
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
	
	return nil
}

// startHTTP starts the HTTP server
func (s *Server) startHTTP() {
	mux := http.NewServeMux()
	
	// Health check
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/metrics", s.handleMetrics)
	
	// API v1
	mux.HandleFunc("/v1/events", s.handleEvents)
	mux.HandleFunc("/v1/events/", s.handleEventByID)
	mux.HandleFunc("/v1/query", s.handleQuery)
	mux.HandleFunc("/v1/tenants", s.handleTenants)
	mux.HandleFunc("/v1/tenants/", s.handleTenantByID)
	mux.HandleFunc("/v1/dashboards", s.handleDashboards)
	mux.HandleFunc("/v1/dashboards/", s.handleDashboardByID)
	mux.HandleFunc("/v1/alerts", s.handleAlerts)
	mux.HandleFunc("/v1/alerts/", s.handleAlertByID)
	
	// CORS middleware
	handler := s.corsMiddleware(mux)
	
	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.config.Host, s.config.HTTPPort),
		Handler:        handler,
		ReadTimeout:    s.config.ReadTimeout,
		WriteTimeout:   s.config.WriteTimeout,
		IdleTimeout:    s.config.IdleTimeout,
		MaxHeaderBytes: s.config.MaxHeaderBytes,
	}
	
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("HTTP server error: %v\n", err)
	}
}

// startGRPC starts the gRPC server
func (s *Server) startGRPC() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.GRPCPort))
	if err != nil {
		fmt.Printf("Failed to listen: %v\n", err)
		return
	}
	
	s.grpcServer = grpc.NewServer()
	
	// Register services
	// In production, would register gRPC service implementations
	
	// Enable reflection
	reflection.Register(s.grpcServer)
	
	if err := s.grpcServer.Serve(lis); err != nil {
		fmt.Printf("gRPC server error: %v\n", err)
	}
}

// Handlers

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"version": "0.1.0",
		"uptime":  time.Since(time.Now()).String(),
	})
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"events_received":    0,
		"events_processed":   0,
		"active_tenants":     0,
		"memory_usage_bytes": 0,
	})
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listEvents(w, r)
	case http.MethodPost:
		s.createEvent(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleEventByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/events/"):]
	switch r.Method {
	case http.MethodGet:
		s.getEvent(w, r, id)
	case http.MethodDelete:
		s.deleteEvent(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleQuery(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.executeQuery(w, r)
	case http.MethodPost:
		s.executeQuery(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleTenants(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listTenants(w, r)
	case http.MethodPost:
		s.createTenant(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleTenantByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/tenants/"):]
	switch r.Method {
	case http.MethodGet:
		s.getTenant(w, r, id)
	case http.MethodPut:
		s.updateTenant(w, r, id)
	case http.MethodDelete:
		s.deleteTenant(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleDashboards(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listDashboards(w, r)
	case http.MethodPost:
		s.createDashboard(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleDashboardByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/dashboards/"):]
	switch r.Method {
	case http.MethodGet:
		s.getDashboard(w, r, id)
	case http.MethodPut:
		s.updateDashboard(w, r, id)
	case http.MethodDelete:
		s.deleteDashboard(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleAlerts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listAlerts(w, r)
	case http.MethodPost:
		s.createAlert(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleAlertByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/alerts/"):]
	switch r.Method {
	case http.MethodGet:
		s.getAlert(w, r, id)
	case http.MethodPut:
		s.updateAlert(w, r, id)
	case http.MethodDelete:
		s.deleteAlert(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Implementation stubs

func (s *Server) listEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"events": []interface{}{},
		"total":  0,
	})
}

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	var event map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      "evt-123",
	})
}

func (s *Server) getEvent(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}

func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request, id string) {
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) executeQuery(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"resultType": "vector",
		"result":     []interface{}{},
		"stats": map[string]interface{}{
			"rowsRead":         0,
			"rowsFiltered":     0,
			"executionTimeMs":  0,
			"cacheHits":        0,
			"cacheMisses":      0,
		},
	})
	_ = query
}

func (s *Server) listTenants(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tenants": []interface{}{},
	})
}

func (s *Server) createTenant(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": "tenant-1",
	})
}

func (s *Server) getTenant(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}

func (s *Server) updateTenant(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}

func (s *Server) deleteTenant(w http.ResponseWriter, r *http.Request, id string) {
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listDashboards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"dashboards": []interface{}{},
	})
}

func (s *Server) createDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": "dashboard-1",
	})
}

func (s *Server) getDashboard(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}

func (s *Server) updateDashboard(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}

func (s *Server) deleteDashboard(w http.ResponseWriter, r *http.Request, id string) {
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listAlerts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"alerts": []interface{}{},
	})
}

func (s *Server) createAlert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": "alert-1",
	})
}

func (s *Server) getAlert(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}

func (s *Server) updateAlert(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}

func (s *Server) deleteAlert(w http.ResponseWriter, r *http.Request, id string) {
	w.WriteHeader(http.StatusNoContent)
}

// Middleware

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Rate limiting logic would go here
		next.ServeHTTP(w, r)
	})
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		fmt.Printf("%s %s %v\n", r.Method, r.URL.Path, duration)
	})
}

// DefaultServerConfig returns default configuration
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		HTTPPort:       8080,
		GRPCPort:       50051,
		Host:           "0.0.0.0",
		TLSEnabled:     false,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

// CreateServer creates a server with default config
func CreateServer() *Server {
	return NewServer(DefaultServerConfig())
}