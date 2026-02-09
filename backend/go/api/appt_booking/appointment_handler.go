package appt_booking

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	appt_booking_db "k8s-fullstack-blueprint-backend/db/appt_booking"
	appt_booking_service "k8s-fullstack-blueprint-backend/service/appt_booking"
)

// AppointmentHandler handles appointment endpoints
type AppointmentHandler struct {
	service *appt_booking_service.ApptBookingService
}

// NewAppointmentHandler creates a new appointment handler
func NewAppointmentHandler(service *appt_booking_service.ApptBookingService) *AppointmentHandler {
	return &AppointmentHandler{
		service: service,
	}
}

// AppointmentResponse represents the response for an appointment (without price)
type AppointmentResponse struct {
	ID                 int    `json:"id"`
	CustomerName       string `json:"customer_name"`
	CustomerEmail      string `json:"customer_email"`
	CustomerPhone      string `json:"customer_phone"`
	StaffID            int    `json:"staff_id"`
	ServiceID          int    `json:"service_id"`
	AppointmentDatetime string `json:"appointment_datetime"`
	DurationMinutes    int    `json:"duration_minutes"`
	Status             string `json:"status"`
	Notes              string    `json:"notes"`
	CreatedAt          string    `json:"created_at"`
	UpdatedAt          string    `json:"updated_at"`
}

// AppointmentWithDetailsResponse represents an appointment with service price
type AppointmentWithDetailsResponse struct {
	ID                 int    `json:"id"`
	CustomerName       string `json:"customer_name"`
	CustomerEmail      string `json:"customer_email"`
	CustomerPhone      string `json:"customer_phone"`
	StaffID            int    `json:"staff_id"`
	ServiceID          int    `json:"service_id"`
	AppointmentDatetime string `json:"appointment_datetime"`
	DurationMinutes    int    `json:"duration_minutes"`
	Status             string `json:"status"`
	Notes              string    `json:"notes"`
	CreatedAt          string    `json:"created_at"`
	UpdatedAt          string    `json:"updated_at"`
	PriceCents         int    `json:"price_cents"`
}

// GetAll handles GET /api/appt_booking/appointments
func (ah *AppointmentHandler) GetAll(c echo.Context) error {
	// Optional query params for filtering
	staffIDStr := c.QueryParam("staff_id")
	email := c.QueryParam("email")

	var appointments []appt_booking_db.AppointmentWithService
	var err error

	if staffIDStr != "" {
		staffID, _ := strconv.Atoi(staffIDStr)
		appointments, err = ah.service.GetAppointmentsByStaffWithDetails(staffID)
	} else if email != "" {
		appointments, err = ah.service.GetAppointmentsByCustomerWithDetails(email)
	} else {
		// Default: get all appointments with service details for admin dashboard
		appointments, err = ah.service.GetAppointmentsWithDetails()
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch appointments",
		})
	}

	response := make([]AppointmentWithDetailsResponse, len(appointments))
	for i, a := range appointments {
		response[i] = AppointmentWithDetailsResponse{
			ID:                   a.ID,
			CustomerName:         a.CustomerName,
			CustomerEmail:        a.CustomerEmail,
			CustomerPhone:        a.CustomerPhone,
			StaffID:              a.StaffID,
			ServiceID:            a.ServiceID,
			AppointmentDatetime:  a.AppointmentDatetime.Format("2006-01-02T15:04:05Z07:00"),
			DurationMinutes:      a.DurationMinutes,
			Status:               a.Status,
			Notes:                a.Notes,
			CreatedAt:            a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:            a.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			PriceCents:           a.PriceCents,
		}
	}

	return c.JSON(http.StatusOK, response)
}

// GetByID handles GET /api/appt_booking/appointments/:id
func (ah *AppointmentHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid appointment ID",
		})
	}

	appointment, err := ah.service.GetAppointment(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Appointment not found",
		})
	}

	response := AppointmentResponse{
		ID:                   appointment.ID,
		CustomerName:         appointment.CustomerName,
		CustomerEmail:        appointment.CustomerEmail,
		CustomerPhone:        appointment.CustomerPhone,
		StaffID:              appointment.StaffID,
		ServiceID:            appointment.ServiceID,
		AppointmentDatetime:  appointment.AppointmentDatetime.Format("2006-01-02T15:04:05Z07:00"),
		DurationMinutes:      appointment.DurationMinutes,
		Status:               appointment.Status,
		Notes:                appointment.Notes,
		CreatedAt:            appointment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:            appointment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusOK, response)
}

// BookRequest represents the request for booking an appointment
type BookRequest struct {
	CustomerName       string    `json:"customer_name"`
	CustomerEmail      string    `json:"customer_email"`
	CustomerPhone      string    `json:"customer_phone"`
	StaffID            int       `json:"staff_id"`
	ServiceID          int       `json:"service_id"`
	AppointmentDatetime string   `json:"appointment_datetime"` // Expected format: "2006-01-02T15:04:05"
	Notes              string    `json:"notes"`
}

// Book handles POST /api/appt_booking/appointments
func (ah *AppointmentHandler) Book(c echo.Context) error {
	var req BookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	// Validate required fields
	if req.CustomerName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Customer name is required",
		})
	}
	if req.CustomerEmail == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Customer email is required",
		})
	}
	if req.StaffID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Valid staff ID is required",
		})
	}
	if req.ServiceID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Valid service ID is required",
		})
	}

	// Parse appointment datetime (expected format: "2006-01-02T15:04:05Z" or "2006-01-02T15:04:05")
	// Use time.ParseInLocation with UTC to handle timezone-aware timestamps
	apptTime, err := time.Parse(time.RFC3339, req.AppointmentDatetime)
	if err != nil {
		// Try without timezone (assume UTC)
		apptTime, err = time.Parse("2006-01-02T15:04:05", req.AppointmentDatetime)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid appointment datetime format. Use YYYY-MM-DDTHH:MM:SS or ISO 8601",
			})
		}
		// Treat naive time as UTC
		apptTime = time.Date(apptTime.Year(), apptTime.Month(), apptTime.Day(), 
			apptTime.Hour(), apptTime.Minute(), apptTime.Second(), 0, time.UTC)
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid appointment datetime format. Use YYYY-MM-DDTHH:MM:SS",
		})
	}

	appointment, err := ah.service.BookAppointment(
		req.CustomerName,
		req.CustomerEmail,
		req.CustomerPhone,
		req.StaffID,
		req.ServiceID,
		apptTime,
		req.Notes,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	response := AppointmentResponse{
		ID:                   appointment.ID,
		CustomerName:         appointment.CustomerName,
		CustomerEmail:        appointment.CustomerEmail,
		CustomerPhone:        appointment.CustomerPhone,
		StaffID:              appointment.StaffID,
		ServiceID:            appointment.ServiceID,
		AppointmentDatetime:  appointment.AppointmentDatetime.Format("2006-01-02T15:04:05Z07:00"),
		DurationMinutes:      appointment.DurationMinutes,
		Status:               appointment.Status,
		Notes:                appointment.Notes,
		CreatedAt:            appointment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:            appointment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusCreated, response)
}

// Cancel handles PUT /api/appt_booking/appointments/:id/cancel
func (ah *AppointmentHandler) Cancel(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid appointment ID",
		})
	}

	err = ah.service.CancelAppointment(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Appointment cancelled successfully",
	})
}

// Complete handles PUT /api/appt_booking/appointments/:id/complete
func (ah *AppointmentHandler) Complete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid appointment ID",
		})
	}

	err = ah.service.CompleteAppointment(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Appointment marked as completed",
	})
}
