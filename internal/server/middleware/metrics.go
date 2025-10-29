package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/observability"
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

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		// Emit counter
		observability.TelemetrySystem.Counter(
			"http_requests_total",
			1,
			map[string]string{
				"method": r.Method,
				"path":   r.URL.Path,
				"status": strconv.Itoa(wrapped.statusCode),
			},
		)

		// Emit histogram for duration
		observability.TelemetrySystem.Histogram(
			"http_request_duration_ms",
			duration,
			map[string]string{
				"method": r.Method,
				"path":   r.URL.Path,
			},
		)
	})
}
