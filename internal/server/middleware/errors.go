package middleware

import (
	"net/http"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/server/handlers"
)

// Recovery middleware recovers from panics and returns standardized error responses
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic (would use proper logger in production)
				// For now, just return a standardized internal server error
				handlers.WriteInternalError(w, "Internal server error occurred")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// ErrorHandler wraps a handler to provide standardized error responses
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to capture the status code
		wrapped := &errorResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		// If status code indicates an error and no content was written, use standardized error response
		if wrapped.statusCode >= 400 && wrapped.statusCode != http.StatusNotFound && !wrapped.written {
			switch wrapped.statusCode {
			case http.StatusMethodNotAllowed:
				handlers.WriteMethodNotAllowed(w, "The requested method is not allowed for this resource")
			case http.StatusRequestTimeout:
				handlers.WriteTimeoutError(w, "Request timed out")
			case http.StatusTooManyRequests:
				handlers.WriteTooManyRequests(w, "Too many requests")
			case http.StatusInternalServerError:
				handlers.WriteInternalError(w, "Internal server error occurred")
			case http.StatusServiceUnavailable:
				handlers.WriteServiceUnavailable(w, "Service temporarily unavailable")
			default:
				handlers.WriteError(w, wrapped.statusCode, "UNKNOWN_ERROR", "An error occurred", nil)
			}
		}

		// Handle 404 case specifically
		if wrapped.statusCode == http.StatusNotFound && !wrapped.written {
			handlers.WriteNotFoundError(w, "The requested resource was not found")
		}
	})
}

// errorResponseWriter wraps http.ResponseWriter to capture status code and written state
type errorResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *errorResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *errorResponseWriter) Write(data []byte) (int, error) {
	rw.written = true
	return rw.ResponseWriter.Write(data)
}
