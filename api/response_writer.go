// Package api provides response writer helpers
package api

import (
	"encoding/json"
	"net/http"
)

// WriteJSON writes JSON response
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// WriteSuccess writes success
func WriteSuccess(w http.ResponseWriter, v interface{}) {
	WriteJSON(w, http.StatusOK, v)
}

// WriteCreated writes created
func WriteCreated(w http.ResponseWriter, v interface{}) {
	WriteJSON(w, http.StatusCreated, v)
}

// WriteNoContent writes no content
func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
