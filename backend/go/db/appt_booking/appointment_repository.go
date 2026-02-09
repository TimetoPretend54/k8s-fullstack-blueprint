package appt_booking

import (
	"database/sql"
	"time"
)

// AppointmentRepository handles database operations for appointments
type AppointmentRepository struct {
	db *sql.DB
}

// NewAppointmentRepository creates a new appointment repository
func NewAppointmentRepository(db *sql.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

// Create inserts a new appointment
func (ar *AppointmentRepository) Create(customerName, customerEmail, customerPhone string, staffID, serviceID, durationMinutes int, appointmentDatetime time.Time, status, notes string) (*Appointment, error) {
	now := time.Now()
	appointment := &Appointment{}
	err := ar.db.QueryRow(
		`INSERT INTO appointments (customer_name, customer_email, customer_phone, staff_id, service_id, appointment_datetime, duration_minutes, status, notes, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
		 RETURNING id, customer_name, customer_email, customer_phone, staff_id, service_id, appointment_datetime, duration_minutes, status, notes, created_at, updated_at`,
		customerName, customerEmail, customerPhone, staffID, serviceID, appointmentDatetime, durationMinutes, status, notes, now, now,
	).Scan(&appointment.ID, &appointment.CustomerName, &appointment.CustomerEmail, &appointment.CustomerPhone, &appointment.StaffID, &appointment.ServiceID, &appointment.AppointmentDatetime, &appointment.DurationMinutes, &appointment.Status, &appointment.Notes, &appointment.CreatedAt, &appointment.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return appointment, nil
}

// Update modifies an existing appointment
func (ar *AppointmentRepository) Update(id int, customerName, customerEmail, customerPhone string, staffID, serviceID, durationMinutes int, appointmentDatetime time.Time, status, notes string) (*Appointment, error) {
	now := time.Now()
	appointment := &Appointment{}
	err := ar.db.QueryRow(
		`UPDATE appointments 
		 SET customer_name = $1, customer_email = $2, customer_phone = $3, staff_id = $4, service_id = $5, appointment_datetime = $6, duration_minutes = $7, status = $8, notes = $9, updated_at = $10 
		 WHERE id = $11 
		 RETURNING id, customer_name, customer_email, customer_phone, staff_id, service_id, appointment_datetime, duration_minutes, status, notes, created_at, updated_at`,
		customerName, customerEmail, customerPhone, staffID, serviceID, appointmentDatetime, durationMinutes, status, notes, now, id,
	).Scan(&appointment.ID, &appointment.CustomerName, &appointment.CustomerEmail, &appointment.CustomerPhone, &appointment.StaffID, &appointment.ServiceID, &appointment.AppointmentDatetime, &appointment.DurationMinutes, &appointment.Status, &appointment.Notes, &appointment.CreatedAt, &appointment.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return appointment, nil
}

// GetAll retrieves all appointments
func (ar *AppointmentRepository) GetAll() ([]Appointment, error) {
	rows, err := ar.db.Query(
		`SELECT id, customer_name, customer_email, customer_phone, staff_id, service_id, appointment_datetime, duration_minutes, status, notes, created_at, updated_at 
		 FROM appointments 
		 ORDER BY appointment_datetime DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var a Appointment
		if err := rows.Scan(&a.ID, &a.CustomerName, &a.CustomerEmail, &a.CustomerPhone, &a.StaffID, &a.ServiceID, &a.AppointmentDatetime, &a.DurationMinutes, &a.Status, &a.Notes, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}

// GetByID retrieves a single appointment by ID
func (ar *AppointmentRepository) GetByID(id int) (*Appointment, error) {
	a := &Appointment{}
	err := ar.db.QueryRow(
		`SELECT id, customer_name, customer_email, customer_phone, staff_id, service_id, appointment_datetime, duration_minutes, status, notes, created_at, updated_at 
		 FROM appointments 
		 WHERE id = $1`,
		id,
	).Scan(&a.ID, &a.CustomerName, &a.CustomerEmail, &a.CustomerPhone, &a.StaffID, &a.ServiceID, &a.AppointmentDatetime, &a.DurationMinutes, &a.Status, &a.Notes, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}

// GetByStaff retrieves all appointments for a specific staff member
func (ar *AppointmentRepository) GetByStaff(staffID int) ([]Appointment, error) {
	rows, err := ar.db.Query(
		`SELECT id, customer_name, customer_email, customer_phone, staff_id, service_id, appointment_datetime, duration_minutes, status, notes, created_at, updated_at 
		 FROM appointments 
		 WHERE staff_id = $1 
		 ORDER BY appointment_datetime DESC`,
		staffID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var a Appointment
		if err := rows.Scan(&a.ID, &a.CustomerName, &a.CustomerEmail, &a.CustomerPhone, &a.StaffID, &a.ServiceID, &a.AppointmentDatetime, &a.DurationMinutes, &a.Status, &a.Notes, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}

// GetByCustomerEmail retrieves all appointments for a customer by email
func (ar *AppointmentRepository) GetByCustomerEmail(email string) ([]Appointment, error) {
	rows, err := ar.db.Query(
		`SELECT id, customer_name, customer_email, customer_phone, staff_id, service_id, appointment_datetime, duration_minutes, status, notes, created_at, updated_at 
		 FROM appointments 
		 WHERE customer_email = $1 
		 ORDER BY appointment_datetime DESC`,
		email,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var a Appointment
		if err := rows.Scan(&a.ID, &a.CustomerName, &a.CustomerEmail, &a.CustomerPhone, &a.StaffID, &a.ServiceID, &a.AppointmentDatetime, &a.DurationMinutes, &a.Status, &a.Notes, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}

// GetUpcoming retrieves upcoming appointments (from now onwards)
func (ar *AppointmentRepository) GetUpcoming(limit int) ([]Appointment, error) {
	rows, err := ar.db.Query(
		`SELECT id, customer_name, customer_email, customer_phone, staff_id, service_id, appointment_datetime, duration_minutes, status, notes, created_at, updated_at 
		 FROM appointments 
		 WHERE appointment_datetime >= NOW() AND status != 'cancelled'
		 ORDER BY appointment_datetime ASC
		 LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var a Appointment
		if err := rows.Scan(&a.ID, &a.CustomerName, &a.CustomerEmail, &a.CustomerPhone, &a.StaffID, &a.ServiceID, &a.AppointmentDatetime, &a.DurationMinutes, &a.Status, &a.Notes, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}

// CheckConflict returns true if there is a conflicting appointment for the given staff at the given datetime
func (ar *AppointmentRepository) CheckConflict(staffID int, appointmentTime time.Time, durationMinutes int, excludeID ...int) (bool, error) {
	endTime := appointmentTime.Add(time.Duration(durationMinutes) * time.Minute)
	
	// Query for any existing appointment that overlaps with the requested time slot
	// Overlap condition: existing.start < new.end AND existing.end > new.start
	query := `
		SELECT COUNT(*) 
		FROM appointments 
		WHERE staff_id = $1 
		  AND status != 'cancelled'
		  AND appointment_datetime < $2 
		  AND (appointment_datetime + (duration_minutes * INTERVAL '1 minute')) > $3
	`
	args := []interface{}{staffID, endTime, appointmentTime}
	
	// Exclude current appointment ID if provided (for updates)
	if len(excludeID) > 0 {
		query += " AND id != $4"
		args = append(args, excludeID[0])
	}
	
	var count int
	err := ar.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// Delete removes an appointment
func (ar *AppointmentRepository) Delete(id int) error {
	_, err := ar.db.Exec("DELETE FROM appointments WHERE id = $1", id)
	return err
}

// Cancel updates an appointment status to 'cancelled'
func (ar *AppointmentRepository) Cancel(id int) error {
	now := time.Now()
	_, err := ar.db.Exec(
		"UPDATE appointments SET status = 'cancelled', updated_at = $1 WHERE id = $2",
		now, id,
	)
	return err
}

// Complete updates an appointment status to 'completed'
func (ar *AppointmentRepository) Complete(id int) error {
	now := time.Now()
	_, err := ar.db.Exec(
		"UPDATE appointments SET status = 'completed', updated_at = $1 WHERE id = $2",
		now, id,
	)
	return err
}

// AppointmentWithService represents an appointment joined with service price
type AppointmentWithService struct {
	ID                 int       `json:"id"`
	CustomerName       string    `json:"customer_name"`
	CustomerEmail      string    `json:"customer_email"`
	CustomerPhone      string    `json:"customer_phone"`
	StaffID            int       `json:"staff_id"`
	ServiceID          int       `json:"service_id"`
	AppointmentDatetime time.Time `json:"appointment_datetime"`
	DurationMinutes    int       `json:"duration_minutes"`
	Status             string    `json:"status"`
	Notes              string    `json:"notes"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	PriceCents         int       `json:"price_cents"`
}

// GetAllWithServiceDetails retrieves all appointments with service price
func (ar *AppointmentRepository) GetAllWithServiceDetails() ([]AppointmentWithService, error) {
	rows, err := ar.db.Query(`
		SELECT
			a.id, a.customer_name, a.customer_email, a.customer_phone,
			a.staff_id, a.service_id, a.appointment_datetime, a.duration_minutes,
			a.status, a.notes, a.created_at, a.updated_at,
			s.price_cents
		FROM appointments a
		JOIN services s ON a.service_id = s.id
		ORDER BY a.appointment_datetime DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []AppointmentWithService
	for rows.Next() {
		var a AppointmentWithService
		if err := rows.Scan(
			&a.ID, &a.CustomerName, &a.CustomerEmail, &a.CustomerPhone,
			&a.StaffID, &a.ServiceID, &a.AppointmentDatetime, &a.DurationMinutes,
			&a.Status, &a.Notes, &a.CreatedAt, &a.UpdatedAt,
			&a.PriceCents,
		); err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}

// GetByStaffWithServiceDetails retrieves appointments for a staff member with service price
func (ar *AppointmentRepository) GetByStaffWithServiceDetails(staffID int) ([]AppointmentWithService, error) {
	rows, err := ar.db.Query(`
		SELECT
			a.id, a.customer_name, a.customer_email, a.customer_phone,
			a.staff_id, a.service_id, a.appointment_datetime, a.duration_minutes,
			a.status, a.notes, a.created_at, a.updated_at,
			s.price_cents
		FROM appointments a
		JOIN services s ON a.service_id = s.id
		WHERE a.staff_id = $1
		ORDER BY a.appointment_datetime DESC
	`, staffID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []AppointmentWithService
	for rows.Next() {
		var a AppointmentWithService
		if err := rows.Scan(
			&a.ID, &a.CustomerName, &a.CustomerEmail, &a.CustomerPhone,
			&a.StaffID, &a.ServiceID, &a.AppointmentDatetime, &a.DurationMinutes,
			&a.Status, &a.Notes, &a.CreatedAt, &a.UpdatedAt,
			&a.PriceCents,
		); err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}

// GetByCustomerEmailWithServiceDetails retrieves appointments for a customer with service price
func (ar *AppointmentRepository) GetByCustomerEmailWithServiceDetails(email string) ([]AppointmentWithService, error) {
	rows, err := ar.db.Query(`
		SELECT
			a.id, a.customer_name, a.customer_email, a.customer_phone,
			a.staff_id, a.service_id, a.appointment_datetime, a.duration_minutes,
			a.status, a.notes, a.created_at, a.updated_at,
			s.price_cents
		FROM appointments a
		JOIN services s ON a.service_id = s.id
		WHERE a.customer_email = $1
		ORDER BY a.appointment_datetime DESC
	`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []AppointmentWithService
	for rows.Next() {
		var a AppointmentWithService
		if err := rows.Scan(
			&a.ID, &a.CustomerName, &a.CustomerEmail, &a.CustomerPhone,
			&a.StaffID, &a.ServiceID, &a.AppointmentDatetime, &a.DurationMinutes,
			&a.Status, &a.Notes, &a.CreatedAt, &a.UpdatedAt,
			&a.PriceCents,
		); err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}

// GetUpcomingWithServiceDetails retrieves upcoming appointments with service price
func (ar *AppointmentRepository) GetUpcomingWithServiceDetails(limit int) ([]AppointmentWithService, error) {
	if limit <= 0 {
		limit = 50 // default
	}
	rows, err := ar.db.Query(`
		SELECT
			a.id, a.customer_name, a.customer_email, a.customer_phone,
			a.staff_id, a.service_id, a.appointment_datetime, a.duration_minutes,
			a.status, a.notes, a.created_at, a.updated_at,
			s.price_cents
		FROM appointments a
		JOIN services s ON a.service_id = s.id
		WHERE a.appointment_datetime >= NOW()
		ORDER BY a.appointment_datetime ASC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []AppointmentWithService
	for rows.Next() {
		var a AppointmentWithService
		if err := rows.Scan(
			&a.ID, &a.CustomerName, &a.CustomerEmail, &a.CustomerPhone,
			&a.StaffID, &a.ServiceID, &a.AppointmentDatetime, &a.DurationMinutes,
			&a.Status, &a.Notes, &a.CreatedAt, &a.UpdatedAt,
			&a.PriceCents,
		); err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}
