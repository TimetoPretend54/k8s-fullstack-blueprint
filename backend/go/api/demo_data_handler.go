package api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"k8s-fullstack-blueprint-backend/service"
)

// DemoDataHandler handles demo data endpoints
type DemoDataHandler struct {
	demoDataService *service.DemoDataService
}

// NewDemoDataHandler creates a new demo data handler
func NewDemoDataHandler(demoDataService *service.DemoDataService) *DemoDataHandler {
	return &DemoDataHandler{
		demoDataService: demoDataService,
	}
}

// DemoDataRequest represents the request for creating/updating demo data
type DemoDataRequest struct {
	Content string `json:"content"`
}

// DemoDataResponse represents the response for demo data
type DemoDataResponse struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetAll handles GET /api/demo-data
func (dh *DemoDataHandler) GetAll(c echo.Context) error {
	records, err := dh.demoDataService.GetAllDemoData()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch demo data",
		})
	}

	// Convert to response format
	response := make([]DemoDataResponse, len(records))
	for i, record := range records {
		response[i] = DemoDataResponse{
			ID:        record.ID,
			Content:   record.Content,
			CreatedAt: record.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: record.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return c.JSON(http.StatusOK, response)
}

// Upsert handles POST /api/demo-data (upsert operation)
func (dh *DemoDataHandler) Upsert(c echo.Context) error {
	var req DemoDataRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	// Validate content
	if req.Content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Content cannot be empty",
		})
	}

	// For simplicity, we'll always create a new record (id = 0)
	// In a real app, you'd pass the ID from the request for updates
	record, err := dh.demoDataService.UpsertDemoData(0, req.Content)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to save demo data",
		})
	}

	response := DemoDataResponse{
		ID:        record.ID,
		Content:   record.Content,
		CreatedAt: record.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: record.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusCreated, response)
}

// GetByID handles GET /api/demo-data/:id
func (dh *DemoDataHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID",
		})
	}

	record, err := dh.demoDataService.GetDemoDataByID(id)
	if err != nil {
		// For simplicity, if not found, return 404
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Demo data not found",
		})
	}

	response := DemoDataResponse{
		ID:        record.ID,
		Content:   record.Content,
		CreatedAt: record.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: record.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusOK, response)
}
