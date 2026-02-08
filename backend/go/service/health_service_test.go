package service

import (
	"testing"
)

func TestHealthService_GetHealthStatus(t *testing.T) {
	hs := NewHealthService()
	health := hs.GetHealthStatus()

	if health["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got '%s'", health["status"])
	}

	if health["timestamp"] == "" {
		t.Error("expected timestamp to be set")
	}
}

func TestHealthService_GetAppInfo(t *testing.T) {
	hs := NewHealthService()
	info := hs.GetAppInfo()

	if info["name"] != "k8s-fullstack-blueprint-backend" {
		t.Errorf("expected name 'k8s-fullstack-blueprint-backend', got '%s'", info["name"])
	}

	if info["version"] != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", info["version"])
	}

	if info["go_version"] == "" {
		t.Error("expected go_version to be set")
	}
}
