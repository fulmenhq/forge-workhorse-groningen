package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/metrics"
	"github.com/fulmenhq/forge-workhorse-groningen/internal/observability"
	"github.com/fulmenhq/forge-workhorse-groningen/internal/server/middleware"
	"github.com/fulmenhq/gofulmen/errors"
)

// HandleError central handler for all errors
func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	// Extract structured error
	var envelope *errors.ErrorEnvelope
	if ee, ok := err.(*errors.ErrorEnvelope); ok {
		envelope = ee
	} else {
		// Wrap unexpected errors
		envelope = errors.NewErrorEnvelope("INTERNAL_ERROR", "unexpected error")
		envelope, _ = envelope.WithContext(map[string]interface{}{
			"wrapped_error": err.Error(),
		})
	}

	// Extract correlation ID safely (handle nil request)
	if r != nil {
		if envelope.CorrelationID == "" {
			envelope = envelope.WithCorrelationID(middleware.GetRequestID(r.Context()))
		}
	} else {
		// Request is nil - ensure we have some correlation ID
		if envelope.CorrelationID == "" {
			// Generate fallback ID for non-HTTP errors
			envelope = envelope.WithCorrelationID(generateFallbackCorrelationID())
		}
	}

	// Determine HTTP status
	statusCode := errorToHTTPStatus(envelope)

	// Build response
	response := ErrorResponse{
		Error: ErrorDetail{
			Code:      envelope.Code,
			Message:   envelope.Message,
			Details:   envelope.Context,
			RequestID: envelope.CorrelationID,
		},
	}

	// Log with severity
	if observability.ServerLogger != nil {
		fields := []zap.Field{
			zap.String("error_code", envelope.Code),
			zap.Int("http_status", statusCode),
		}
		if envelope.Context != nil {
			for k, v := range envelope.Context {
				if str, ok := v.(string); ok {
					fields = append(fields, zap.String(k, str))
				}
			}
		}
		if envelope.CorrelationID != "" {
			fields = append(fields, zap.String("request_id", envelope.CorrelationID))
		}

		switch envelope.Severity {
		case errors.SeverityCritical:
			observability.ServerLogger.Error(envelope.Message, fields...)
		case errors.SeverityHigh:
			observability.ServerLogger.Error(envelope.Message, fields...)
		case errors.SeverityMedium:
			observability.ServerLogger.Warn(envelope.Message, fields...)
		default:
			observability.ServerLogger.Info(envelope.Message, fields...)
		}
	}

	// Record metric
	metrics.RecordError(envelope.Code, statusCode)

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

// ErrorResponse structure per API standards
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
}

// errorToHTTPStatus maps error codes to status codes
func errorToHTTPStatus(envelope *errors.ErrorEnvelope) int {
	switch envelope.Code {
	case "INVALID_INPUT", "VALIDATION_FAILED":
		return http.StatusBadRequest
	case "NOT_FOUND":
		return http.StatusNotFound
	case "UNAUTHORIZED":
		return http.StatusUnauthorized
	case "FORBIDDEN":
		return http.StatusForbidden
	case "METHOD_NOT_ALLOWED":
		return http.StatusMethodNotAllowed
	case "CONFLICT":
		return http.StatusConflict
	case "TIMEOUT":
		return http.StatusGatewayTimeout
	case "EXTERNAL_SERVICE_ERROR":
		return http.StatusBadGateway
	case "SERVICE_UNAVAILABLE":
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// WriteServiceUnavailable helper for 503 errors
func WriteServiceUnavailable(w http.ResponseWriter, message string) {
	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    "SERVICE_UNAVAILABLE",
			Message: message,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	_ = json.NewEncoder(w).Encode(response)
}

// Legacy compatibility - redirect to handlers for backward compatibility
func WriteInternalError(w http.ResponseWriter, message string) {
	envelope := errors.NewErrorEnvelope("INTERNAL_SERVER_ERROR", message)
	HandleError(w, nil, envelope)
}

// generateFallbackCorrelationID creates a correlation ID when request context is unavailable
func generateFallbackCorrelationID() string {
	// Generate UUID for fallback scenarios (CLI errors, non-HTTP errors)
	// Prefix with "fallback-" to distinguish from request-based IDs
	return "fallback-" + generateUUID()
}

// generateUUID creates a new UUID string
func generateUUID() string {
	return uuid.New().String()
}
