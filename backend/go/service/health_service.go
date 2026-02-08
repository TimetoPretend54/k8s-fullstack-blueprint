package service

import (
	"os"
	"runtime"
	"time"
)

// HealthService handles health-related business logic
type HealthService struct{}

// NewHealthService creates a new health service
func NewHealthService() *HealthService {
	return &HealthService{}
}

// GetHealthStatus returns the current health status
func (hs *HealthService) GetHealthStatus() map[string]string {
	return map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
}

// GetAppInfo returns application information
func (hs *HealthService) GetAppInfo() map[string]string {
	hostname, _ := os.Hostname()
	return map[string]string{
		"name":       "k8s-fullstack-blueprint-backend",
		"version":    "1.0.0",
		"go_version": runtime.Version(),
		"hostname":   hostname,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	}
}
