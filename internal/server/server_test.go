package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/server/handlers"
)

func TestServerUsesStandardErrorHandlers(t *testing.T) {
	srv := New("127.0.0.1", 0)

	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}

	var body handlers.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if body.Error.Code != handlers.ErrCodeNotFound {
		t.Fatalf("expected error code %s, got %s", handlers.ErrCodeNotFound, body.Error.Code)
	}
}
