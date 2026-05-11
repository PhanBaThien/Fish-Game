package utils

import (
	"encoding/json"
	"net/http"
)

// APIResponse is the standard envelope for all API responses.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// WriteJSON writes a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// WriteSuccess sends a 200 OK response with data.
func WriteSuccess(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: data})
}

// WriteError sends an error response with the given HTTP status.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, APIResponse{Success: false, Error: message})
}

// Ptr returns a pointer to the given value — useful for optional struct fields.
func Ptr[T any](v T) *T {
	return &v
}
