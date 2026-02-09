package appt_booking

import (
	"errors"
	"time"

	"k8s-fullstack-blueprint-backend/db/appt_booking"
)

// ApptBookingService handles business logic for appointment booking
type ApptBookingService struct {
	serviceRepo      *appt_booking.ServiceRepository
	staffRepo        *appt_booking.StaffRepository
	staffServiceRepo *appt_booking.StaffServiceRepository
	scheduleRepo     *appt_booking.ScheduleRepository
	appointmentRepo  *appt_booking.AppointmentRepository
}

// NewApptBookingService creates a new appointment booking service
func NewApptBookingService(
	serviceRepo *appt_booking.ServiceRepository,
	staffRepo *appt_booking.StaffRepository,
	staffServiceRepo *appt_booking.StaffServiceRepository,
	scheduleRepo *appt_booking.ScheduleRepository,
	appointmentRepo *appt_booking.AppointmentRepository,
) *ApptBookingService {
	return &ApptBookingService{
		serviceRepo:      serviceRepo,
		staffRepo:        staffRepo,
		staffServiceRepo: staffServiceRepo,
		scheduleRepo:     scheduleRepo,
		appointmentRepo:  appointmentRepo,
	}
}

// ========== Service Operations ==========

// CreateService creates a new service
func (s *ApptBookingService) CreateService(name, description string, durationMinutes, priceCents int) (*appt_booking.Service, error) {
	// Validation
	if name == "" {
		return nil, errors.New("service name is required")
	}
	if durationMinutes <= 0 {
		return nil, errors.New("duration must be positive")
	}
	if priceCents < 0 {
		return nil, errors.New("price cannot be negative")
	}

	return s.serviceRepo.Create(name, description, durationMinutes, priceCents)
}

// UpdateService modifies an existing service
func (s *ApptBookingService) UpdateService(id int, name, description string, durationMinutes, priceCents int) (*appt_booking.Service, error) {
	// Validation
	if name == "" {
		return nil, errors.New("service name is required")
	}
	if durationMinutes <= 0 {
		return nil, errors.New("duration must be positive")
	}
	if priceCents < 0 {
		return nil, errors.New("price cannot be negative")
	}

	return s.serviceRepo.Update(id, name, description, durationMinutes, priceCents)
}

// GetAllServices retrieves all services
func (s *ApptBookingService) GetAllServices() ([]appt_booking.Service, error) {
	return s.serviceRepo.GetAll()
}

// GetServiceByID retrieves a service by ID
func (s *ApptBookingService) GetServiceByID(id int) (*appt_booking.Service, error) {
	return s.serviceRepo.GetByID(id)
}

// DeleteService removes a service
func (s *ApptBookingService) DeleteService(id int) error {
	// Check if service is in use by staff
	staffList, err := s.staffServiceRepo.GetStaffForService(id)
	if err != nil {
		return err
	}
	if len(staffList) > 0 {
		return errors.New("cannot delete service that is assigned to staff members")
	}

	// Check if service has existing appointments
	// Note: We could also cascade delete or restrict based on business rules
	// For now, we'll prevent deletion if appointments exist
	appointments, err := s.appointmentRepo.GetAll()
	if err != nil {
		return err
	}
	for _, a := range appointments {
		if a.ServiceID == id {
			return errors.New("cannot delete service that has existing appointments")
		}
	}

	return s.serviceRepo.Delete(id)
}

// ========== Staff Operations ==========

// CreateStaff creates a new staff member
func (s *ApptBookingService) CreateStaff(name, email, phone, role string) (*appt_booking.Staff, error) {
	// Validation
	if name == "" {
		return nil, errors.New("staff name is required")
	}
	if email == "" {
		return nil, errors.New("staff email is required")
	}
	if role == "" {
		return nil, errors.New("staff role is required")
	}
	// Basic email format check
	if !contains(email, "@") {
		return nil, errors.New("invalid email format")
	}

	// Check if email already exists
	existing, err := s.staffRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("staff with this email already exists")
	}

	return s.staffRepo.Create(name, email, phone, role)
}

// UpdateStaff modifies an existing staff member
func (s *ApptBookingService) UpdateStaff(id int, name, email, phone, role string) (*appt_booking.Staff, error) {
	// Validation
	if name == "" {
		return nil, errors.New("staff name is required")
	}
	if email == "" {
		return nil, errors.New("staff email is required")
	}
	if role == "" {
		return nil, errors.New("staff role is required")
	}

	// Check if email is used by another staff member
	existing, err := s.staffRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.ID != id {
		return nil, errors.New("email is already used by another staff member")
	}

	return s.staffRepo.Update(id, name, email, phone, role)
}

// GetAllStaff retrieves all staff members
func (s *ApptBookingService) GetAllStaff() ([]appt_booking.Staff, error) {
	return s.staffRepo.GetAll()
}

// GetStaffByID retrieves a staff member by ID
func (s *ApptBookingService) GetStaffByID(id int) (*appt_booking.Staff, error) {
	return s.staffRepo.GetByID(id)
}

// DeleteStaff removes a staff member
func (s *ApptBookingService) DeleteStaff(id int) error {
	// Check if staff has existing appointments
	appointments, err := s.appointmentRepo.GetByStaff(id)
	if err != nil {
		return err
	}
	if len(appointments) > 0 {
		return errors.New("cannot delete staff member with existing appointments")
	}

	// Cascade delete schedules and staff_service assignments
	// This is handled by ON DELETE CASCADE in the database schema
	return s.staffRepo.Delete(id)
}

// ========== Staff-Service Assignment Operations ==========

// AssignServiceToStaff links a service to a staff member
func (s *ApptBookingService) AssignServiceToStaff(staffID, serviceID int) error {
	// Validate staff exists
	staff, err := s.staffRepo.GetByID(staffID)
	if err != nil {
		return err
	}
	if staff == nil {
		return errors.New("staff not found")
	}

	// Validate service exists
	service, err := s.serviceRepo.GetByID(serviceID)
	if err != nil {
		return err
	}
	if service == nil {
		return errors.New("service not found")
	}

	// Check if already assigned
	existingServices, err := s.staffServiceRepo.GetServicesForStaff(staffID)
	if err != nil {
		return err
	}
	for _, svc := range existingServices {
		if svc.ID == serviceID {
			return errors.New("service is already assigned to this staff member")
		}
	}

	return s.staffServiceRepo.Assign(staffID, serviceID)
}

// UnassignServiceFromStaff removes a service from a staff member
func (s *ApptBookingService) UnassignServiceFromStaff(staffID, serviceID int) error {
	return s.staffServiceRepo.Unassign(staffID, serviceID)
}

// GetServicesForStaff retrieves all services offered by a staff member
func (s *ApptBookingService) GetServicesForStaff(staffID int) ([]appt_booking.Service, error) {
	return s.staffServiceRepo.GetServicesForStaff(staffID)
}

// GetStaffForService retrieves all staff members who offer a service
func (s *ApptBookingService) GetStaffForService(serviceID int) ([]appt_booking.Staff, error) {
	return s.staffServiceRepo.GetStaffForService(serviceID)
}

// ========== Schedule Operations ==========

// CreateSchedule creates a new schedule entry for a staff member
func (s *ApptBookingService) CreateSchedule(staffID, dayOfWeek int, startTime, endTime string) (*appt_booking.Schedule, error) {
	// Validation
	if dayOfWeek < 0 || dayOfWeek > 6 {
		return nil, errors.New("day of week must be between 0 (Sunday) and 6 (Saturday)")
	}
	if startTime == "" || endTime == "" {
		return nil, errors.New("start time and end time are required")
	}
	// Basic time format validation (HH:MM)
	if !isValidTimeFormat(startTime) || !isValidTimeFormat(endTime) {
		return nil, errors.New("time must be in HH:MM format")
	}

	// Validate staff exists
	staff, err := s.staffRepo.GetByID(staffID)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, errors.New("staff not found")
	}

	// Parse input times to time.Time for validation
	startTimeParsed, err := time.Parse("15:04", startTime)
	if err != nil {
		return nil, errors.New("invalid start time format, must be HH:MM")
	}
	endTimeParsed, err := time.Parse("15:04", endTime)
	if err != nil {
		return nil, errors.New("invalid end time format, must be HH:MM")
	}

	// Check for overlapping schedules for the same staff and day
	existingSchedules, err := s.scheduleRepo.GetByStaff(staffID)
	if err != nil {
		return nil, err
	}
	for _, sch := range existingSchedules {
		if sch.DayOfWeek == dayOfWeek {
			// Simple overlap check: if time ranges intersect
			if timesOverlap(startTimeParsed, endTimeParsed, sch.StartTime, sch.EndTime) {
				return nil, errors.New("schedule overlaps with an existing schedule for this staff member on the same day")
			}
		}
	}

	return s.scheduleRepo.Create(staffID, dayOfWeek, startTime, endTime)
}

// UpdateSchedule modifies an existing schedule
func (s *ApptBookingService) UpdateSchedule(id int, dayOfWeek int, startTime, endTime string) (*appt_booking.Schedule, error) {
	// Validation
	if dayOfWeek < 0 || dayOfWeek > 6 {
		return nil, errors.New("day of week must be between 0 (Sunday) and 6 (Saturday)")
	}
	if startTime == "" || endTime == "" {
		return nil, errors.New("start time and end time are required")
	}
	if !isValidTimeFormat(startTime) || !isValidTimeFormat(endTime) {
		return nil, errors.New("time must be in HH:MM format")
	}

	// Get existing schedule to check staff ID
	existing, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("schedule not found")
	}

	// Parse input times to time.Time for validation
	startTimeParsed, err := time.Parse("15:04", startTime)
	if err != nil {
		return nil, errors.New("invalid start time format, must be HH:MM")
	}
	endTimeParsed, err := time.Parse("15:04", endTime)
	if err != nil {
		return nil, errors.New("invalid end time format, must be HH:MM")
	}

	// Check for overlapping schedules (excluding current)
	schedules, err := s.scheduleRepo.GetByStaff(existing.StaffID)
	if err != nil {
		return nil, err
	}
	for _, sch := range schedules {
		if sch.ID != id && sch.DayOfWeek == dayOfWeek {
			if timesOverlap(startTimeParsed, endTimeParsed, sch.StartTime, sch.EndTime) {
				return nil, errors.New("schedule overlaps with an existing schedule for this staff member on the same day")
			}
		}
	}

	return s.scheduleRepo.Update(id, dayOfWeek, startTime, endTime)
}

// GetAllSchedules retrieves all schedules
func (s *ApptBookingService) GetAllSchedules() ([]appt_booking.Schedule, error) {
	return s.scheduleRepo.GetAll()
}

// GetScheduleByID retrieves a schedule by ID
func (s *ApptBookingService) GetScheduleByID(id int) (*appt_booking.Schedule, error) {
	return s.scheduleRepo.GetByID(id)
}

// GetSchedulesByStaff retrieves all schedules for a staff member
func (s *ApptBookingService) GetSchedulesByStaff(staffID int) ([]appt_booking.Schedule, error) {
	return s.scheduleRepo.GetByStaff(staffID)
}

// DeleteSchedule removes a schedule
func (s *ApptBookingService) DeleteSchedule(id int) error {
	return s.scheduleRepo.Delete(id)
}

// ========== Appointment Operations ==========

// BookAppointment creates a new appointment with conflict checking
func (s *ApptBookingService) BookAppointment(
	customerName, customerEmail, customerPhone string,
	staffID, serviceID int,
	appointmentDatetime time.Time,
	notes string,
) (*appt_booking.Appointment, error) {
	// Validation
	if customerName == "" {
		return nil, errors.New("customer name is required")
	}
	if customerEmail == "" {
		return nil, errors.New("customer email is required")
	}
	if !contains(customerEmail, "@") {
		return nil, errors.New("invalid email format")
	}
	if staffID <= 0 {
		return nil, errors.New("valid staff ID is required")
	}
	if serviceID <= 0 {
		return nil, errors.New("valid service ID is required")
	}

	// Validate staff exists
	staff, err := s.staffRepo.GetByID(staffID)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, errors.New("staff not found")
	}

	// Validate service exists
	service, err := s.serviceRepo.GetByID(serviceID)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, errors.New("service not found")
	}

	// Check if staff offers this service
	staffServices, err := s.staffServiceRepo.GetServicesForStaff(staffID)
	if err != nil {
		return nil, err
	}
	serviceOffered := false
	for _, svc := range staffServices {
		if svc.ID == serviceID {
			serviceOffered = true
			break
		}
	}
	if !serviceOffered {
		return nil, errors.New("staff member does not offer this service")
	}

	// Check staff schedule for the appointment day/time
	// Convert appointment from UTC to staff's local timezone (America/Los_Angeles)
	// Schedules are stored in local time, so we need to compare apples-to-apples
	// TODO: Add DB field for "staff"'s local timezone
	loc, _ := time.LoadLocation("America/Los_Angeles")
	appointmentLocal := appointmentDatetime.In(loc)
	appointmentDay := int(appointmentLocal.Weekday())
	// Go's time.Weekday: Sunday=0, Monday=1, etc. matches our schema
	schedules, err := s.scheduleRepo.GetByStaff(staffID)
	if err != nil {
		return nil, err
	}
	appointmentEndTimeLocal := appointmentLocal.Add(time.Duration(service.DurationMin) * time.Minute)

	isWithinSchedule := false
	for _, sch := range schedules {
		if sch.DayOfWeek == appointmentDay {
			if isTimeWithin(sch.StartTime, sch.EndTime, appointmentLocal, appointmentEndTimeLocal) {
				isWithinSchedule = true
				break
			}
		}
	}

	if !isWithinSchedule {
		return nil, errors.New("appointment time is outside staff member's working hours")
	}

	// Check for conflicts with existing appointments
	hasConflict, err := s.appointmentRepo.CheckConflict(staffID, appointmentDatetime, service.DurationMin)
	if err != nil {
		return nil, err
	}
	if hasConflict {
		return nil, errors.New("appointment time conflicts with an existing appointment")
	}

	// Create the appointment
	return s.appointmentRepo.Create(
		customerName,
		customerEmail,
		customerPhone,
		staffID,
		serviceID,
		service.DurationMin,
		appointmentDatetime,
		"confirmed",
		notes,
	)
}

// GetAppointment retrieves an appointment by ID
func (s *ApptBookingService) GetAppointment(id int) (*appt_booking.Appointment, error) {
	return s.appointmentRepo.GetByID(id)
}

// GetAppointmentsByStaff retrieves all appointments for a staff member
func (s *ApptBookingService) GetAppointmentsByStaff(staffID int) ([]appt_booking.Appointment, error) {
	return s.appointmentRepo.GetByStaff(staffID)
}

// GetAppointmentsByCustomer retrieves all appointments for a customer by email
func (s *ApptBookingService) GetAppointmentsByCustomer(email string) ([]appt_booking.Appointment, error) {
	return s.appointmentRepo.GetByCustomerEmail(email)
}

// GetUpcomingAppointments retrieves upcoming appointments
func (s *ApptBookingService) GetUpcomingAppointments(limit int) ([]appt_booking.Appointment, error) {
	if limit <= 0 {
		limit = 50 // default
	}
	return s.appointmentRepo.GetUpcoming(limit)
}

// CancelAppointment cancels an appointment
func (s *ApptBookingService) CancelAppointment(id int) error {
	// Check if appointment exists
	appt, err := s.appointmentRepo.GetByID(id)
	if err != nil {
		return err
	}
	if appt == nil {
		return errors.New("appointment not found")
	}

	if appt.Status == "cancelled" {
		return errors.New("appointment is already cancelled")
	}

	return s.appointmentRepo.Cancel(id)
}

// CompleteAppointment marks an appointment as completed
func (s *ApptBookingService) CompleteAppointment(id int) error {
	appt, err := s.appointmentRepo.GetByID(id)
	if err != nil {
		return err
	}
	if appt == nil {
		return errors.New("appointment not found")
	}

	if appt.Status == "completed" {
		return errors.New("appointment is already completed")
	}
	if appt.Status == "cancelled" {
		return errors.New("cannot complete a cancelled appointment")
	}

	return s.appointmentRepo.Complete(id)
}

// GetAppointmentsWithDetails retrieves all appointments with service price for revenue calculation
func (s *ApptBookingService) GetAppointmentsWithDetails() ([]appt_booking.AppointmentWithService, error) {
	return s.appointmentRepo.GetAllWithServiceDetails()
}

// GetAppointmentsByStaffWithDetails retrieves appointments for a staff member with service price
func (s *ApptBookingService) GetAppointmentsByStaffWithDetails(staffID int) ([]appt_booking.AppointmentWithService, error) {
	return s.appointmentRepo.GetByStaffWithServiceDetails(staffID)
}

// GetAppointmentsByCustomerWithDetails retrieves appointments for a customer with service price
func (s *ApptBookingService) GetAppointmentsByCustomerWithDetails(email string) ([]appt_booking.AppointmentWithService, error) {
	return s.appointmentRepo.GetByCustomerEmailWithServiceDetails(email)
}

// GetUpcomingAppointmentsWithDetails retrieves upcoming appointments with service price
func (s *ApptBookingService) GetUpcomingAppointmentsWithDetails(limit int) ([]appt_booking.AppointmentWithService, error) {
	return s.appointmentRepo.GetUpcomingWithServiceDetails(limit)
}

// ========== Helper Functions ==========

// isValidTimeFormat checks if time string is in HH:MM format
func isValidTimeFormat(t string) bool {
	if len(t) != 5 {
		return false
	}
	if t[2] != ':' {
		return false
	}
	hour := t[0:2]
	minute := t[3:5]
	// Simple check: digits only
	for i := 0; i < 2; i++ {
		if hour[i] < '0' || hour[i] > '9' {
			return false
		}
	}
	for i := 0; i < 2; i++ {
		if minute[i] < '0' || minute[i] > '9' {
			return false
		}
	}
	// Check numeric ranges
	h := int(hour[0]-'0')*10 + int(hour[1]-'0')
	m := int(minute[0]-'0')*10 + int(minute[1]-'0')
	return h >= 0 && h <= 23 && m >= 0 && m <= 59
}

// timesOverlap checks if two time ranges overlap
// All parameters are time.Time values (date component ignored, only time-of-day used)
func timesOverlap(start1, end1, start2, end2 time.Time) bool {
	// Convert to minutes since midnight for easy comparison
	toMinutes := func(t time.Time) int {
		return t.Hour()*60 + t.Minute()
	}
	s1 := toMinutes(start1)
	e1 := toMinutes(end1)
	s2 := toMinutes(start2)
	e2 := toMinutes(end2)
	// Overlap if: s1 < e2 AND s2 < e1
	return s1 < e2 && s2 < e1
}

// isTimeWithin checks if an appointment time range fits within a schedule slot
// scheduleStart/End: schedule time-of-day (time.Time with zero date)
// apptStart/End: appointment datetime (full time.Time, but only time-of-day considered)
func isTimeWithin(scheduleStart, scheduleEnd, apptStart, apptEnd time.Time) bool {
	toMinutes := func(t time.Time) int {
		return t.Hour()*60 + t.Minute()
	}
	ss := toMinutes(scheduleStart)
	se := toMinutes(scheduleEnd)
	as := toMinutes(apptStart)
	ae := toMinutes(apptEnd)
	// Appointment must be fully within schedule
	return as >= ss && ae <= se
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
