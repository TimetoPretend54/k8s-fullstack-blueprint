package appt_booking

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"k8s-fullstack-blueprint-backend/service/appt_booking"
)

// ScheduleHandler handles schedule endpoints
type ScheduleHandler struct {
	service *appt_booking.ApptBookingService
}

// NewScheduleHandler creates a new schedule handler
func NewScheduleHandler(service *appt_booking.ApptBookingService) *ScheduleHandler {
	return &ScheduleHandler{
		service: service,
	}
}

// ScheduleResponse represents the response for a schedule
type ScheduleResponse struct {
	ID        int    `json:"id"`
	StaffID   int    `json:"staff_id"`
	DayOfWeek int    `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetAll handles GET /api/appt_booking/schedules
func (sh *ScheduleHandler) GetAll(c echo.Context) error {
	schedules, err := sh.service.GetAllSchedules()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch schedules",
		})
	}

	response := make([]ScheduleResponse, len(schedules))
	for i, s := range schedules {
		response[i] = ScheduleResponse{
			ID:        s.ID,
			StaffID:   s.StaffID,
			DayOfWeek: s.DayOfWeek,
		StartTime: s.StartTime.Format("15:04"),
		EndTime:   s.EndTime.Format("15:04"),
			CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return c.JSON(http.StatusOK, response)
}

// GetByID handles GET /api/appt_booking/schedules/:id
func (sh *ScheduleHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid schedule ID",
		})
	}

	schedule, err := sh.service.GetScheduleByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Schedule not found",
		})
	}

	response := ScheduleResponse{
		ID:        schedule.ID,
		StaffID:   schedule.StaffID,
		DayOfWeek: schedule.DayOfWeek,
		StartTime: schedule.StartTime.Format("15:04"),
		EndTime:   schedule.EndTime.Format("15:04"),
		CreatedAt: schedule.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: schedule.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusOK, response)
}

// GetByStaff handles GET /api/appt_booking/schedules/staff/:staffId
func (sh *ScheduleHandler) GetByStaff(c echo.Context) error {
	staffIDStr := c.Param("staffId")
	staffID, err := strconv.Atoi(staffIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid staff ID",
		})
	}

	schedules, err := sh.service.GetSchedulesByStaff(staffID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch schedules for staff",
		})
	}

	response := make([]ScheduleResponse, len(schedules))
	for i, s := range schedules {
		response[i] = ScheduleResponse{
			ID:        s.ID,
			StaffID:   s.StaffID,
			DayOfWeek: s.DayOfWeek,
		StartTime: s.StartTime.Format("15:04"),
		EndTime:   s.EndTime.Format("15:04"),
			CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return c.JSON(http.StatusOK, response)
}

// ScheduleRequest represents the request for creating/updating a schedule
type ScheduleRequest struct {
	StaffID   int    `json:"staff_id"`
	DayOfWeek int    `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// Create handles POST /api/appt_booking/schedules
func (sh *ScheduleHandler) Create(c echo.Context) error {
	var req ScheduleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	// Validate required fields
	if req.StaffID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Valid staff ID is required",
		})
	}
	if req.StartTime == "" || req.EndTime == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Start time and end time are required",
		})
	}

	schedule, err := sh.service.CreateSchedule(req.StaffID, req.DayOfWeek, req.StartTime, req.EndTime)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	response := ScheduleResponse{
		ID:        schedule.ID,
		StaffID:   schedule.StaffID,
		DayOfWeek: schedule.DayOfWeek,
		StartTime: schedule.StartTime.Format("15:04"),
		EndTime:   schedule.EndTime.Format("15:04"),
		CreatedAt: schedule.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: schedule.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusCreated, response)
}

// Update handles PUT /api/appt_booking/schedules/:id
func (sh *ScheduleHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid schedule ID",
		})
	}

	var req ScheduleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	// Validate required fields
	if req.StaffID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Valid staff ID is required",
		})
	}
	if req.StartTime == "" || req.EndTime == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Start time and end time are required",
		})
	}

	schedule, err := sh.service.UpdateSchedule(id, req.DayOfWeek, req.StartTime, req.EndTime)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	response := ScheduleResponse{
		ID:        schedule.ID,
		StaffID:   schedule.StaffID,
		DayOfWeek: schedule.DayOfWeek,
		StartTime: schedule.StartTime.Format("15:04"),
		EndTime:   schedule.EndTime.Format("15:04"),
		CreatedAt: schedule.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: schedule.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusOK, response)
}

// Delete handles DELETE /api/appt_booking/schedules/:id
func (sh *ScheduleHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid schedule ID",
		})
	}

	err = sh.service.DeleteSchedule(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Schedule deleted successfully",
	})
}
