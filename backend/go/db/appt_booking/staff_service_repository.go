package appt_booking

import (
	"database/sql"
)

// StaffServiceRepository handles operations for the staff_services junction table
type StaffServiceRepository struct {
	db *sql.DB
}

// NewStaffServiceRepository creates a new staff service repository
func NewStaffServiceRepository(db *sql.DB) *StaffServiceRepository {
	return &StaffServiceRepository{db: db}
}

// Assign links a staff member to a service
func (sr *StaffServiceRepository) Assign(staffID, serviceID int) error {
	_, err := sr.db.Exec(
		"INSERT INTO staff_services (staff_id, service_id) VALUES ($1, $2)",
		staffID, serviceID,
	)
	return err
}

// Unassign removes a service from a staff member
func (sr *StaffServiceRepository) Unassign(staffID, serviceID int) error {
	_, err := sr.db.Exec(
		"DELETE FROM staff_services WHERE staff_id = $1 AND service_id = $2",
		staffID, serviceID,
	)
	return err
}

// GetServicesForStaff retrieves all services offered by a specific staff member
func (sr *StaffServiceRepository) GetServicesForStaff(staffID int) ([]Service, error) {
	rows, err := sr.db.Query(
		`SELECT s.id, s.name, s.description, s.duration_minutes, s.price_cents, s.created_at, s.updated_at 
		 FROM services s
		 INNER JOIN staff_services ss ON s.id = ss.service_id
		 WHERE ss.staff_id = $1
		 ORDER BY s.name`,
		staffID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var s Service
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.DurationMin, &s.PriceCents, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return services, nil
}

// GetStaffForService retrieves all staff members who offer a specific service
func (sr *StaffServiceRepository) GetStaffForService(serviceID int) ([]Staff, error) {
	rows, err := sr.db.Query(
		`SELECT st.id, st.name, st.email, st.phone, st.role, st.created_at, st.updated_at
		 FROM staff st
		 INNER JOIN staff_services ss ON st.id = ss.staff_id
		 WHERE ss.service_id = $1
		 ORDER BY st.name`,
		serviceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var staffList []Staff
	for rows.Next() {
		var s Staff
		if err := rows.Scan(&s.ID, &s.Name, &s.Email, &s.Phone, &s.Role, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		staffList = append(staffList, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return staffList, nil
}

// RemoveAllForStaff removes all service assignments for a staff member
func (sr *StaffServiceRepository) RemoveAllForStaff(staffID int) error {
	_, err := sr.db.Exec("DELETE FROM staff_services WHERE staff_id = $1", staffID)
	return err
}

// RemoveAllForService removes all staff assignments for a service
func (sr *StaffServiceRepository) RemoveAllForService(serviceID int) error {
	_, err := sr.db.Exec("DELETE FROM staff_services WHERE service_id = $1", serviceID)
	return err
}
