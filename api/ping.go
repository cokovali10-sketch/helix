// Package api provides ping
package api

import (
	"net/http"
	"time"
)

// PingHandler handles ping
type PingHandler struct{}

// NewPingHandler creates handler
func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

// ServeHTTP serves
func (p *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("pong"))
}

// Ping returns timestamp
func Ping() int64 {
	return time.Now().Unix()
}
