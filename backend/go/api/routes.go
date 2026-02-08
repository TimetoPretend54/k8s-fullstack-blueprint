package api

import (
	"github.com/labstack/echo/v4"
)

// SetupRoutes configures all API routes
func SetupRoutes(e *echo.Echo, healthHandler *HealthHandler, demoDataHandler *DemoDataHandler) {
	// Health check endpoints
	e.GET("/", healthHandler.Root)
	e.GET("/health", healthHandler.Check)
	e.GET("/info", healthHandler.Info)

	// Demo data endpoints
	e.GET("/api/demo-data", demoDataHandler.GetAll)
	e.GET("/api/demo-data/:id", demoDataHandler.GetByID)
	e.POST("/api/demo-data", demoDataHandler.Upsert)
}
