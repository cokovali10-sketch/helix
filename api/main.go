// Package api provides gRPC and REST API definitions
package api

import (
	"fmt"
	"net/http"
)

// API represents the Helix platform API server
type API struct {
	version  string
	host     string
	port     string
	grpcPort string
}

// NewAPI creates a new API server
func NewAPI(version, host, grpcPort, httpPort string) *API {
	return &API{
		version:  version,
		host:     host,
		port:     httpPort,
		grpcPort: grpcPort,
	}
}

// Start starts the API server
func (a *API) Start() error {
	fmt.Printf("Starting Helix API v%s on %s:%s\n", a.version, a.host, a.port)
	return nil
}

// SetupRoutes sets up default routes
func (a *API) SetupRoutes() {
	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})
}

// RegisterService registers a service
func (a *API) RegisterService(name string, addr string) {
	fmt.Printf("Registering service: %s at %s\n", name, addr)
}