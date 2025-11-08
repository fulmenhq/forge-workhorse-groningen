package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubChecker struct {
	err error
}

func (s stubChecker) CheckHealth(ctx context.Context) error {
	return s.err
}

func TestHealthHandlerReturnsHealthyStatus(t *testing.T) {
	manager := NewHealthManager("1.2.3")
	manager.RegisterChecker("ok", stubChecker{err: nil})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	manager.HealthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "healthy" {
		t.Fatalf("expected healthy status, got %s", resp.Status)
	}

	if resp.Version != "1.2.3" {
		t.Fatalf("expected version 1.2.3, got %s", resp.Version)
	}

	if resp.Checks["ok"] != "healthy" {
		t.Fatalf("expected ok check to be healthy, got %s", resp.Checks["ok"])
	}
}

func TestHealthHandlerReturnsServiceUnavailableWhenUnhealthy(t *testing.T) {
	manager := NewHealthManager("1.2.3")
	manager.RegisterChecker("db", stubChecker{err: errors.New("down")})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	manager.HealthHandler(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", rec.Code)
	}

	var resp HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "unhealthy" {
		t.Fatalf("expected unhealthy status, got %s", resp.Status)
	}

	if resp.Checks["db"] != "unhealthy" {
		t.Fatalf("expected db check to be unhealthy, got %s", resp.Checks["db"])
	}
}

func TestDetermineOverallStatusTreatsTimeoutAsDegraded(t *testing.T) {
	manager := NewHealthManager("dev")

	status := manager.determineOverallStatus(map[string]string{
		"db": "timeout",
	})

	if status != "degraded" {
		t.Fatalf("expected degraded status, got %s", status)
	}
}
