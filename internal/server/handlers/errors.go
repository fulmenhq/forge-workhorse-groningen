package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents standardized error response format
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// WriteError writes a standardized error response
func WriteError(w http.ResponseWriter, statusCode int, code, message string, details interface{}) {
	resp := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}

// Common error codes and messages
const (
	ErrCodeInternalServer     = "INTERNAL_SERVER_ERROR"
	ErrCodeNotFound           = "NOT_FOUND"
	ErrCodeBadRequest         = "BAD_REQUEST"
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeForbidden          = "FORBIDDEN"
	ErrCodeMethodNotAllowed   = "METHOD_NOT_ALLOWED"
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
	ErrCodeTimeout            = "TIMEOUT"
	ErrCodeTooManyRequests    = "TOO_MANY_REQUESTS"
)

// Convenience functions for common error responses
func WriteInternalError(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusInternalServerError, ErrCodeInternalServer, message, nil)
}

func WriteNotFoundError(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusNotFound, ErrCodeNotFound, message, nil)
}

func WriteBadRequest(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusBadRequest, ErrCodeBadRequest, message, nil)
}

func WriteUnauthorizedError(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusUnauthorized, ErrCodeUnauthorized, message, nil)
}

func WriteForbiddenError(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusForbidden, ErrCodeForbidden, message, nil)
}

func WriteMethodNotAllowed(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, message, nil)
}

func WriteServiceUnavailable(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusServiceUnavailable, ErrCodeServiceUnavailable, message, nil)
}

func WriteTimeoutError(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusRequestTimeout, ErrCodeTimeout, message, nil)
}

func WriteTooManyRequests(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusTooManyRequests, ErrCodeTooManyRequests, message, nil)
}

// NotFoundHandler handles 404 errors with standardized response
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	WriteNotFoundError(w, "The requested resource was not found")
}

// MethodNotAllowedHandler handles 405 errors with standardized response
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	WriteMethodNotAllowed(w, "The requested method is not allowed for this resource")
}
