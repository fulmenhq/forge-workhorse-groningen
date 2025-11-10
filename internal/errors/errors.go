package errors

import (
	"context"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/server/middleware"
	"github.com/fulmenhq/gofulmen/errors"
	"github.com/google/uuid"
)

// Error creation helpers for common error types

// User Errors (400-level)
func NewInvalidInputError(message string) *errors.ErrorEnvelope {
	return errors.NewErrorEnvelope("INVALID_INPUT", message)
}

func NewNotFoundError(message string) *errors.ErrorEnvelope {
	return errors.NewErrorEnvelope("NOT_FOUND", message)
}

func NewUnauthorizedError(message string) *errors.ErrorEnvelope {
	return errors.NewErrorEnvelope("UNAUTHORIZED", message)
}

func NewForbiddenError(message string) *errors.ErrorEnvelope {
	return errors.NewErrorEnvelope("FORBIDDEN", message)
}

func NewMethodNotAllowedError(message string) *errors.ErrorEnvelope {
	return errors.NewErrorEnvelope("METHOD_NOT_ALLOWED", message)
}

func NewConflictError(message string) *errors.ErrorEnvelope {
	return errors.NewErrorEnvelope("CONFLICT", message)
}

func NewValidationError(message string) *errors.ErrorEnvelope {
	return errors.NewErrorEnvelope("VALIDATION_FAILED", message)
}

// Server Errors (500-level)
func NewInternalError(message string) *errors.ErrorEnvelope {
	return errors.NewErrorEnvelope("INTERNAL_ERROR", message)
}

func NewDatabaseError(message string) *errors.ErrorEnvelope {
	return errors.NewErrorEnvelope("DATABASE_ERROR", message)
}

func NewExternalServiceError(message string) *errors.ErrorEnvelope {
	return errors.NewErrorEnvelope("EXTERNAL_SERVICE_ERROR", message)
}

func NewTimeoutError(message string) *errors.ErrorEnvelope {
	return errors.NewErrorEnvelope("TIMEOUT", message)
}

// Application-Specific Errors
func NewDataProcessingError(message string) *errors.ErrorEnvelope {
	return errors.NewErrorEnvelope("DATA_PROCESSING_ERROR", message)
}

func NewConfigInvalidError(message string) *errors.ErrorEnvelope {
	return errors.NewErrorEnvelope("CONFIG_INVALID", message)
}

// Wrap functions for existing errors
// These functions accept a context to extract correlation/trace IDs from the request context

func WrapInvalidInput(ctx context.Context, err error, message string) *errors.ErrorEnvelope {
	envelope := errors.NewErrorEnvelope("INVALID_INPUT", message)
	envelope = envelope.WithCorrelationID(extractCorrelationID(ctx))
	envelope = envelope.WithTraceID(extractTraceID(ctx))
	envelope, _ = envelope.WithContext(map[string]interface{}{
		"wrapped_error": err.Error(),
	})
	return envelope
}

func WrapNotFound(ctx context.Context, err error, message string) *errors.ErrorEnvelope {
	envelope := errors.NewErrorEnvelope("NOT_FOUND", message)
	envelope = envelope.WithCorrelationID(extractCorrelationID(ctx))
	envelope = envelope.WithTraceID(extractTraceID(ctx))
	envelope, _ = envelope.WithContext(map[string]interface{}{
		"wrapped_error": err.Error(),
	})
	return envelope
}

func WrapUnauthorized(ctx context.Context, err error, message string) *errors.ErrorEnvelope {
	envelope := errors.NewErrorEnvelope("UNAUTHORIZED", message)
	envelope = envelope.WithCorrelationID(extractCorrelationID(ctx))
	envelope = envelope.WithTraceID(extractTraceID(ctx))
	envelope, _ = envelope.WithContext(map[string]interface{}{
		"wrapped_error": err.Error(),
	})
	return envelope
}

func WrapForbidden(ctx context.Context, err error, message string) *errors.ErrorEnvelope {
	envelope := errors.NewErrorEnvelope("FORBIDDEN", message)
	envelope = envelope.WithCorrelationID(extractCorrelationID(ctx))
	envelope = envelope.WithTraceID(extractTraceID(ctx))
	envelope, _ = envelope.WithContext(map[string]interface{}{
		"wrapped_error": err.Error(),
	})
	return envelope
}

func WrapConflict(ctx context.Context, err error, message string) *errors.ErrorEnvelope {
	envelope := errors.NewErrorEnvelope("CONFLICT", message)
	envelope = envelope.WithCorrelationID(extractCorrelationID(ctx))
	envelope = envelope.WithTraceID(extractTraceID(ctx))
	envelope, _ = envelope.WithContext(map[string]interface{}{
		"wrapped_error": err.Error(),
	})
	return envelope
}

func WrapValidationError(ctx context.Context, err error, message string) *errors.ErrorEnvelope {
	envelope := errors.NewErrorEnvelope("VALIDATION_FAILED", message)
	envelope = envelope.WithCorrelationID(extractCorrelationID(ctx))
	envelope = envelope.WithTraceID(extractTraceID(ctx))
	envelope, _ = envelope.WithContext(map[string]interface{}{
		"wrapped_error": err.Error(),
	})
	return envelope
}

func WrapInternal(ctx context.Context, err error, message string) *errors.ErrorEnvelope {
	envelope := errors.NewErrorEnvelope("INTERNAL_ERROR", message)
	envelope = envelope.WithCorrelationID(extractCorrelationID(ctx))
	envelope = envelope.WithTraceID(extractTraceID(ctx))
	envelope, _ = envelope.WithContext(map[string]interface{}{
		"wrapped_error": err.Error(),
	})
	return envelope
}

func WrapDatabaseError(ctx context.Context, err error, message string) *errors.ErrorEnvelope {
	envelope := errors.NewErrorEnvelope("DATABASE_ERROR", message)
	envelope = envelope.WithCorrelationID(extractCorrelationID(ctx))
	envelope = envelope.WithTraceID(extractTraceID(ctx))
	envelope, _ = envelope.WithContext(map[string]interface{}{
		"wrapped_error": err.Error(),
	})
	return envelope
}

func WrapExternalService(ctx context.Context, err error, message string) *errors.ErrorEnvelope {
	envelope := errors.NewErrorEnvelope("EXTERNAL_SERVICE_ERROR", message)
	envelope = envelope.WithCorrelationID(extractCorrelationID(ctx))
	envelope = envelope.WithTraceID(extractTraceID(ctx))
	envelope, _ = envelope.WithContext(map[string]interface{}{
		"wrapped_error": err.Error(),
	})
	return envelope
}

func WrapTimeout(ctx context.Context, err error, message string) *errors.ErrorEnvelope {
	envelope := errors.NewErrorEnvelope("TIMEOUT", message)
	envelope = envelope.WithCorrelationID(extractCorrelationID(ctx))
	envelope = envelope.WithTraceID(extractTraceID(ctx))
	envelope, _ = envelope.WithContext(map[string]interface{}{
		"wrapped_error": err.Error(),
	})
	return envelope
}

func WrapDataProcessing(ctx context.Context, err error, message string) *errors.ErrorEnvelope {
	envelope := errors.NewErrorEnvelope("DATA_PROCESSING_ERROR", message)
	envelope = envelope.WithCorrelationID(extractCorrelationID(ctx))
	envelope = envelope.WithTraceID(extractTraceID(ctx))
	envelope, _ = envelope.WithContext(map[string]interface{}{
		"wrapped_error": err.Error(),
	})
	return envelope
}

func WrapConfigInvalid(ctx context.Context, err error, message string) *errors.ErrorEnvelope {
	envelope := errors.NewErrorEnvelope("CONFIG_INVALID", message)
	envelope = envelope.WithCorrelationID(extractCorrelationID(ctx))
	envelope = envelope.WithTraceID(extractTraceID(ctx))
	envelope, _ = envelope.WithContext(map[string]interface{}{
		"wrapped_error": err.Error(),
	})
	return envelope
}

// Helper functions for ID generation

// extractCorrelationID gets correlation ID from context, falls back to generating new UUID
func extractCorrelationID(ctx context.Context) string {
	if ctx != nil {
		if requestID := middleware.GetRequestID(ctx); requestID != "" {
			return requestID
		}
	}
	// Fallback: generate new UUID when context is nil or has no request ID
	return uuid.New().String()
}

// extractTraceID gets trace ID from context, falls back to generating new UUID
func extractTraceID(ctx context.Context) string {
	// TODO: Extract from OpenTelemetry or other tracing system when implemented
	// For now, use correlation ID as trace ID
	return extractCorrelationID(ctx)
}
