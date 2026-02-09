package main

import (
	"database/sql"
	"fmt"
	"log"

	"k8s-fullstack-blueprint-backend/api"
	"k8s-fullstack-blueprint-backend/api/appt_booking"
	"k8s-fullstack-blueprint-backend/db"
	"k8s-fullstack-blueprint-backend/service"
	appt_booking_db "k8s-fullstack-blueprint-backend/db/appt_booking"
	appt_booking_service "k8s-fullstack-blueprint-backend/service/appt_booking"
)

// DependencyContainer holds all application dependencies
type DependencyContainer struct {
	HealthHandler      *api.HealthHandler
	DemoDataHandler    *api.DemoDataHandler
	// Appointment Booking handlers
	ServiceHandler     *appt_booking.ServiceHandler
	StaffHandler       *appt_booking.StaffHandler
	ScheduleHandler    *appt_booking.ScheduleHandler
	AppointmentHandler *appt_booking.AppointmentHandler
	// Repositories (for direct access if needed)
	ApptBookingDB      *sql.DB
	ServiceRepo        *appt_booking_db.ServiceRepository
	StaffRepo          *appt_booking_db.StaffRepository
	StaffServiceRepo   *appt_booking_db.StaffServiceRepository
	ScheduleRepo       *appt_booking_db.ScheduleRepository
	AppointmentRepo    *appt_booking_db.AppointmentRepository
	ApptBookingService *appt_booking_service.ApptBookingService
}

// NewDependencyContainer constructs and wires all dependencies
func NewDependencyContainer() (*DependencyContainer, error) {
	log.Println("Initializing application dependencies...")

	// Initialize main database connection (for demo_data)
	dbConn, err := db.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize main database schema
	if err := db.InitSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	// Initialize appointment booking database connection (separate database)
	apptBookingDB, err := appt_booking_db.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to appointment booking database: %w", err)
	}

	// Initialize appointment booking database schema
	if err := appt_booking_db.InitSchema(apptBookingDB); err != nil {
		return nil, fmt.Errorf("failed to initialize appointment booking schema: %w", err)
	}

	// Seed sample data for appointment booking (idempotent)
	if err := appt_booking_db.SeedSampleData(apptBookingDB); err != nil {
		return nil, fmt.Errorf("failed to seed appointment booking sample data: %w", err)
	}

	log.Println("All dependencies initialized successfully")

	// Initialize main repository layer
	demoDataRepo := db.NewDemoDataRepository(dbConn)

	// Initialize appointment booking repository layer
	serviceRepo := appt_booking_db.NewServiceRepository(apptBookingDB)
	staffRepo := appt_booking_db.NewStaffRepository(apptBookingDB)
	staffServiceRepo := appt_booking_db.NewStaffServiceRepository(apptBookingDB)
	scheduleRepo := appt_booking_db.NewScheduleRepository(apptBookingDB)
	appointmentRepo := appt_booking_db.NewAppointmentRepository(apptBookingDB)

	// Initialize service layer
	healthService := service.NewHealthService()
	demoDataService := service.NewDemoDataService(demoDataRepo)
	apptBookingService := appt_booking_service.NewApptBookingService(serviceRepo, staffRepo, staffServiceRepo, scheduleRepo, appointmentRepo)

	// Initialize API layer with dependencies
	healthHandler := api.NewHealthHandler(healthService)
	demoDataHandler := api.NewDemoDataHandler(demoDataService)
	serviceHandler := appt_booking.NewServiceHandler(apptBookingService)
	staffHandler := appt_booking.NewStaffHandler(apptBookingService)
	scheduleHandler := appt_booking.NewScheduleHandler(apptBookingService)
	appointmentHandler := appt_booking.NewAppointmentHandler(apptBookingService)

	return &DependencyContainer{
		HealthHandler:      healthHandler,
		DemoDataHandler:    demoDataHandler,
		ServiceHandler:     serviceHandler,
		StaffHandler:       staffHandler,
		ScheduleHandler:    scheduleHandler,
		AppointmentHandler: appointmentHandler,
		ApptBookingDB:      apptBookingDB,
		ServiceRepo:        serviceRepo,
		StaffRepo:          staffRepo,
		StaffServiceRepo:   staffServiceRepo,
		ScheduleRepo:       scheduleRepo,
		AppointmentRepo:    appointmentRepo,
		ApptBookingService: apptBookingService,
	}, nil
}
