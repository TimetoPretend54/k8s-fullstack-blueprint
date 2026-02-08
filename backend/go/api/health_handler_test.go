package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"k8s-fullstack-blueprint-backend/service"
)

func TestHealthHandler_Root(t *testing.T) {
	// Setup
	healthService := service.NewHealthService()
	handler := NewHealthHandler(healthService)

	e := echo.New()
	e.GET("/", handler.Root)

	// Execute
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expected := "Hello from k8s-fullstack-blueprint Go backend!"
	if rec.Body.String() != expected {
		t.Errorf("expected body '%s', got '%s'", expected, rec.Body.String())
	}
}

func TestHealthHandler_Check(t *testing.T) {
	healthService := service.NewHealthService()
	handler := NewHealthHandler(healthService)

	e := echo.New()
	e.GET("/health", handler.Check)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Verify JSON response contains status and timestamp
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", rec.Header().Get("Content-Type"))
	}
}

func TestHealthHandler_Info(t *testing.T) {
	healthService := service.NewHealthService()
	handler := NewHealthHandler(healthService)

	e := echo.New()
	e.GET("/info", handler.Info)

	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", rec.Header().Get("Content-Type"))
	}
}
