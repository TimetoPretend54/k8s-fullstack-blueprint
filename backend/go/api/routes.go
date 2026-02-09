package api

import (
	"github.com/labstack/echo/v4"

	"k8s-fullstack-blueprint-backend/api/appt_booking"
)

// SetupRoutes configures all API routes
func SetupRoutes(
	e *echo.Echo,
	healthHandler *HealthHandler,
	demoDataHandler *DemoDataHandler,
	serviceHandler *appt_booking.ServiceHandler,
	staffHandler *appt_booking.StaffHandler,
	scheduleHandler *appt_booking.ScheduleHandler,
	appointmentHandler *appt_booking.AppointmentHandler,
) {
	// Health check endpoints
	e.GET("/", healthHandler.Root)
	e.GET("/health", healthHandler.Check)
	e.GET("/info", healthHandler.Info)

	// Demo data endpoints
	e.GET("/api/demo-data", demoDataHandler.GetAll)
	e.GET("/api/demo-data/:id", demoDataHandler.GetByID)
	e.POST("/api/demo-data", demoDataHandler.Upsert)

	// Appointment Booking endpoints
	// Services
	e.GET("/api/appt_booking/services", serviceHandler.GetAll)
	e.GET("/api/appt_booking/services/:id", serviceHandler.GetByID)
	e.POST("/api/appt_booking/services", serviceHandler.Create)
	e.PUT("/api/appt_booking/services/:id", serviceHandler.Update)
	e.DELETE("/api/appt_booking/services/:id", serviceHandler.Delete)

	// Staff
	e.GET("/api/appt_booking/staff", staffHandler.GetAll)
	e.GET("/api/appt_booking/staff/:id", staffHandler.GetByID)
	e.POST("/api/appt_booking/staff", staffHandler.Create)
	e.PUT("/api/appt_booking/staff/:id", staffHandler.Update)
	e.DELETE("/api/appt_booking/staff/:id", staffHandler.Delete)
	e.GET("/api/appt_booking/staff/by-service/:serviceId", staffHandler.GetByService)

	// Schedules
	e.GET("/api/appt_booking/schedules", scheduleHandler.GetAll)
	e.GET("/api/appt_booking/schedules/:id", scheduleHandler.GetByID)
	e.GET("/api/appt_booking/schedules/staff/:staffId", scheduleHandler.GetByStaff)
	e.POST("/api/appt_booking/schedules", scheduleHandler.Create)
	e.PUT("/api/appt_booking/schedules/:id", scheduleHandler.Update)
	e.DELETE("/api/appt_booking/schedules/:id", scheduleHandler.Delete)

	// Appointments
	e.GET("/api/appt_booking/appointments", appointmentHandler.GetAll)
	e.GET("/api/appt_booking/appointments/:id", appointmentHandler.GetByID)
	e.POST("/api/appt_booking/appointments", appointmentHandler.Book)
	e.PUT("/api/appt_booking/appointments/:id/cancel", appointmentHandler.Cancel)
	e.PUT("/api/appt_booking/appointments/:id/complete", appointmentHandler.Complete)
}
