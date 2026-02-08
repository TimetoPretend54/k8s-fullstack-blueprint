package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"k8s-fullstack-blueprint-backend/api"
)

func main() {
	// Create a new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize dependency container
	container, err := NewDependencyContainer()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Load routes
	api.SetupRoutes(e, container.HealthHandler, container.DemoDataHandler)

	// Get port from environment or default
	port := getEnv("PORT", "8080")

	// Start server
	log.Printf("Starting server on port %s", port)
	log.Fatal(e.Start(":" + port))
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
