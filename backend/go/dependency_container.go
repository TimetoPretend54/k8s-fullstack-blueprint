package main

import (
	"fmt"
	"log"

	"k8s-fullstack-blueprint-backend/api"
	"k8s-fullstack-blueprint-backend/db"
	"k8s-fullstack-blueprint-backend/service"
)

// DependencyContainer holds all application dependencies
type DependencyContainer struct {
	HealthHandler   *api.HealthHandler
	DemoDataHandler *api.DemoDataHandler
	// Add other handlers as they are created
}

// NewDependencyContainer constructs and wires all dependencies
func NewDependencyContainer() (*DependencyContainer, error) {
	log.Println("Initializing application dependencies...")

	// Initialize database connection
	dbConn, err := db.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize database schema (temporary scaffold - TODO: use proper migrations)
	if err := db.InitSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	log.Println("All dependencies initialized successfully")

	// Initialize repository layer
	demoDataRepo := db.NewDemoDataRepository(dbConn)

	// Initialize service layer
	healthService := service.NewHealthService()
	demoDataService := service.NewDemoDataService(demoDataRepo)

	// Initialize API layer with dependencies
	healthHandler := api.NewHealthHandler(healthService)
	demoDataHandler := api.NewDemoDataHandler(demoDataService)

	return &DependencyContainer{
		HealthHandler:   healthHandler,
		DemoDataHandler: demoDataHandler,
	}, nil
}
