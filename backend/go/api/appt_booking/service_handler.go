package appt_booking

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"k8s-fullstack-blueprint-backend/service/appt_booking"
)

// ServiceHandler handles service endpoints
type ServiceHandler struct {
	service *appt_booking.ApptBookingService
}

// NewServiceHandler creates a new service handler
func NewServiceHandler(service *appt_booking.ApptBookingService) *ServiceHandler {
	return &ServiceHandler{
		service: service,
	}
}

// ServiceResponse represents the response for a service
type ServiceResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DurationMin int    `json:"duration_min"`
	PriceCents  int    `json:"price_cents"`
}

// GetAll handles GET /api/appt_booking/services
func (sh *ServiceHandler) GetAll(c echo.Context) error {
	services, err := sh.service.GetAllServices()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch services",
		})
	}

	response := make([]ServiceResponse, len(services))
	for i, s := range services {
		response[i] = ServiceResponse{
			ID:          s.ID,
			Name:        s.Name,
			Description: s.Description,
			DurationMin: s.DurationMin,
			PriceCents:  s.PriceCents,
		}
	}

	return c.JSON(http.StatusOK, response)
}

// GetByID handles GET /api/appt_booking/services/:id
func (sh *ServiceHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid service ID",
		})
	}

	service, err := sh.service.GetServiceByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Service not found",
		})
	}

	response := ServiceResponse{
		ID:          service.ID,
		Name:        service.Name,
		Description: service.Description,
		DurationMin: service.DurationMin,
		PriceCents:  service.PriceCents,
	}

	return c.JSON(http.StatusOK, response)
}

// ServiceRequest represents the request for creating/updating a service
type ServiceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DurationMin int    `json:"duration_min"`
	PriceCents  int    `json:"price_cents"`
}

// Create handles POST /api/appt_booking/services
func (sh *ServiceHandler) Create(c echo.Context) error {
	var req ServiceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	// Validate required fields
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Service name is required",
		})
	}
	if req.DurationMin <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Duration must be positive",
		})
	}

	service, err := sh.service.CreateService(req.Name, req.Description, req.DurationMin, req.PriceCents)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create service",
		})
	}

	response := ServiceResponse{
		ID:          service.ID,
		Name:        service.Name,
		Description: service.Description,
		DurationMin: service.DurationMin,
		PriceCents:  service.PriceCents,
	}

	return c.JSON(http.StatusCreated, response)
}

// Update handles PUT /api/appt_booking/services/:id
func (sh *ServiceHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid service ID",
		})
	}

	var req ServiceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	// Validate required fields
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Service name is required",
		})
	}
	if req.DurationMin <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Duration must be positive",
		})
	}

	service, err := sh.service.UpdateService(id, req.Name, req.Description, req.DurationMin, req.PriceCents)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update service",
		})
	}

	response := ServiceResponse{
		ID:          service.ID,
		Name:        service.Name,
		Description: service.Description,
		DurationMin: service.DurationMin,
		PriceCents:  service.PriceCents,
	}

	return c.JSON(http.StatusOK, response)
}

// Delete handles DELETE /api/appt_booking/services/:id
func (sh *ServiceHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid service ID",
		})
	}

	err = sh.service.DeleteService(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete service",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Service deleted successfully",
	})
}
