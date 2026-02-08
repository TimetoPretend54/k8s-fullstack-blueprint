package service

import (
	"fmt"

	"k8s-fullstack-blueprint-backend/db"
)

// DemoDataService handles business logic for demo data
type DemoDataService struct {
	repo *db.DemoDataRepository
}

// NewDemoDataService creates a new demo data service
func NewDemoDataService(repo *db.DemoDataRepository) *DemoDataService {
	return &DemoDataService{repo: repo}
}

// UpsertDemoData creates or updates a demo record
func (ds *DemoDataService) UpsertDemoData(id int, content string) (*db.DemoData, error) {
	// Basic validation
	if content == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}

	return ds.repo.Upsert(id, content)
}

// GetAllDemoData returns all demo records
func (ds *DemoDataService) GetAllDemoData() ([]db.DemoData, error) {
	return ds.repo.GetAll()
}

// GetDemoDataByID returns a specific record
func (ds *DemoDataService) GetDemoDataByID(id int) (*db.DemoData, error) {
	return ds.repo.GetByID(id)
}

// DeleteDemoData removes a record
func (ds *DemoDataService) DeleteDemoData(id int) error {
	return ds.repo.Delete(id)
}
