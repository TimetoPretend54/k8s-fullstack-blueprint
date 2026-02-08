package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"k8s-fullstack-blueprint-backend/service"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	healthService *service.HealthService
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(healthService *service.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// Check handles GET /health
func (hh *HealthHandler) Check(c echo.Context) error {
	healthData := hh.healthService.GetHealthStatus()
	response := HealthResponse{
		Status:    healthData["status"],
		Timestamp: healthData["timestamp"],
	}
	return c.JSON(http.StatusOK, response)
}

// AppInfo represents application info
type AppInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	GoVersion string `json:"go_version"`
	Hostname  string `json:"hostname"`
	Timestamp string `json:"timestamp"`
}

// Info handles GET /info
func (hh *HealthHandler) Info(c echo.Context) error {
	appInfo := hh.healthService.GetAppInfo()
	info := AppInfo{
		Name:      appInfo["name"],
		Version:   appInfo["version"],
		GoVersion: appInfo["go_version"],
		Hostname:  appInfo["hostname"],
		Timestamp: appInfo["timestamp"],
	}
	return c.JSON(http.StatusOK, info)
}

// Root handles GET /
func (hh *HealthHandler) Root(c echo.Context) error {
	return c.String(http.StatusOK, "Hello from k8s-fullstack-blueprint Go backend!")
}
