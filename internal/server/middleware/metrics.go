package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/observability"
	"go.uber.org/zap"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// RequestMetrics middleware captures HTTP request metrics
func RequestMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if observability.TelemetrySystem == nil {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		requestID := GetRequestID(r.Context())

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		// Common labels for all metrics
		commonLabels := map[string]string{
			"method":    r.Method,
			"path":      r.URL.Path,
			"status":    strconv.Itoa(wrapped.statusCode),
			"requestID": requestID,
		}

		// Emit counter
		_ = observability.TelemetrySystem.Counter(
			"http_requests_total",
			1,
			commonLabels,
		)

		// Emit histogram for duration
		_ = observability.TelemetrySystem.Histogram(
			"http_request_duration_ms",
			duration,
			map[string]string{
				"method":    r.Method,
				"path":      r.URL.Path,
				"status":    strconv.Itoa(wrapped.statusCode),
				"requestID": requestID,
			},
		)

		// Log request with request ID for tracing
		if observability.ServerLogger != nil {
			observability.ServerLogger.Info("HTTP request completed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", wrapped.statusCode),
				zap.Duration("duration", duration),
				zap.String("requestID", requestID),
			)
		}
	})
}
