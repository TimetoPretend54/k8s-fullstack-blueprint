package appt_booking

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"k8s-fullstack-blueprint-backend/service/appt_booking"
)

// StaffHandler handles staff endpoints
type StaffHandler struct {
	service *appt_booking.ApptBookingService
}

// NewStaffHandler creates a new staff handler
func NewStaffHandler(service *appt_booking.ApptBookingService) *StaffHandler {
	return &StaffHandler{
		service: service,
	}
}

// StaffResponse represents the response for a staff member
type StaffResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetAll handles GET /api/appt_booking/staff
func (sh *StaffHandler) GetAll(c echo.Context) error {
	staffList, err := sh.service.GetAllStaff()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch staff",
		})
	}

	response := make([]StaffResponse, len(staffList))
	for i, s := range staffList {
		response[i] = StaffResponse{
			ID:        s.ID,
			Name:      s.Name,
			Email:     s.Email,
			Phone:     s.Phone,
			Role:      s.Role,
			CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return c.JSON(http.StatusOK, response)
}

// GetByID handles GET /api/appt_booking/staff/:id
func (sh *StaffHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid staff ID",
		})
	}

	staff, err := sh.service.GetStaffByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Staff not found",
		})
	}

	response := StaffResponse{
		ID:        staff.ID,
		Name:      staff.Name,
		Email:     staff.Email,
		Phone:     staff.Phone,
		Role:      staff.Role,
		CreatedAt: staff.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: staff.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusOK, response)
}

// StaffRequest represents the request for creating/updating a staff member
type StaffRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
}

// Create handles POST /api/appt_booking/staff
func (sh *StaffHandler) Create(c echo.Context) error {
	var req StaffRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	// Validate required fields
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Staff name is required",
		})
	}
	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Staff email is required",
		})
	}
	if req.Role == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Staff role is required",
		})
	}

	staff, err := sh.service.CreateStaff(req.Name, req.Email, req.Phone, req.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	response := StaffResponse{
		ID:        staff.ID,
		Name:      staff.Name,
		Email:     staff.Email,
		Phone:     staff.Phone,
		Role:      staff.Role,
		CreatedAt: staff.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: staff.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusCreated, response)
}

// Update handles PUT /api/appt_booking/staff/:id
func (sh *StaffHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid staff ID",
		})
	}

	var req StaffRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	// Validate required fields
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Staff name is required",
		})
	}
	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Staff email is required",
		})
	}
	if req.Role == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Staff role is required",
		})
	}

	staff, err := sh.service.UpdateStaff(id, req.Name, req.Email, req.Phone, req.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	response := StaffResponse{
		ID:        staff.ID,
		Name:      staff.Name,
		Email:     staff.Email,
		Phone:     staff.Phone,
		Role:      staff.Role,
		CreatedAt: staff.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: staff.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusOK, response)
}

// Delete handles DELETE /api/appt_booking/staff/:id
func (sh *StaffHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid staff ID",
		})
	}

	err = sh.service.DeleteStaff(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Staff deleted successfully",
	})
}

// GetByService retrieves all staff members who offer a specific service
func (sh *StaffHandler) GetByService(c echo.Context) error {
	serviceIDStr := c.Param("serviceId")
	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid service ID",
		})
	}

	staffList, err := sh.service.GetStaffForService(serviceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch staff for service",
		})
	}

	response := make([]StaffResponse, len(staffList))
	for i, s := range staffList {
		response[i] = StaffResponse{
			ID:        s.ID,
			Name:      s.Name,
			Email:     s.Email,
			Phone:     s.Phone,
			Role:      s.Role,
			CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return c.JSON(http.StatusOK, response)
}
