// Package api provides error handling
package api

import (
	"encoding/json"
	"net/http"
)

// APIError represents an API error
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// WriteError writes error response
func WriteError(w http.ResponseWriter, code int, message, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(APIError{Code: code, Message: message, Details: details})
}

// BadRequest writes 400
func BadRequest(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusBadRequest, message, "")
}

// NotFound writes 404
func NotFound(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusNotFound, message, "")
}

// InternalError writes 500
func InternalError(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusInternalServerError, message, "")
}

// Unauthorized writes 401
func Unauthorized(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusUnauthorized, message, "")
}

// Forbidden writes 403
func Forbidden(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusForbidden, message, "")
}
