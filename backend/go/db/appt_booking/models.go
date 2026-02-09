package appt_booking

import (
	"time"
)

// Service represents a booking service offered by staff
type Service struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	DurationMin  int       `json:"duration_minutes" db:"duration_minutes"`
	PriceCents   int       `json:"price_cents" db:"price_cents"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Staff represents a provider or admin
type Staff struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	Role      string    `json:"role" db:"role"` // "provider" or "admin"
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// StaffService is a junction table linking staff to services (many-to-many)
type StaffService struct {
	StaffID  int `json:"staff_id" db:"staff_id"`
	ServiceID int `json:"service_id" db:"service_id"`
}

// Schedule represents a recurring availability slot for staff
type Schedule struct {
	ID         int       `json:"id" db:"id"`
	StaffID    int       `json:"staff_id" db:"staff_id"`
	DayOfWeek  int       `json:"day_of_week" db:"day_of_week"`   // 0=Sunday, 6=Saturday
	StartTime  time.Time `json:"start_time" db:"start_time"`     // TIME type, stores time of day
	EndTime    time.Time `json:"end_time" db:"end_time"`         // TIME type, stores time of day
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// Appointment represents a booked appointment
type Appointment struct {
	ID                 int       `json:"id" db:"id"`
	CustomerName       string    `json:"customer_name" db:"customer_name"`
	CustomerEmail      string    `json:"customer_email" db:"customer_email"`
	CustomerPhone      string    `json:"customer_phone" db:"customer_phone"`
	StaffID            int       `json:"staff_id" db:"staff_id"`
	ServiceID          int       `json:"service_id" db:"service_id"`
	AppointmentDatetime time.Time `json:"appointment_datetime" db:"appointment_datetime"`
	DurationMinutes    int       `json:"duration_minutes" db:"duration_minutes"`
	Status             string    `json:"status" db:"status"` // "confirmed", "cancelled", "completed"
	Notes              string    `json:"notes" db:"notes"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}
