package appt_booking

import (
	"database/sql"
	"time"
)

// StaffRepository handles database operations for staff
type StaffRepository struct {
	db *sql.DB
}

// NewStaffRepository creates a new staff repository
func NewStaffRepository(db *sql.DB) *StaffRepository {
	return &StaffRepository{db: db}
}

// Create inserts a new staff member
func (sr *StaffRepository) Create(name, email, phone, role string) (*Staff, error) {
	now := time.Now()
	staff := &Staff{}
	err := sr.db.QueryRow(
		"INSERT INTO staff (name, email, phone, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, email, phone, role, created_at, updated_at",
		name, email, phone, role, now, now,
	).Scan(&staff.ID, &staff.Name, &staff.Email, &staff.Phone, &staff.Role, &staff.CreatedAt, &staff.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return staff, nil
}

// Update modifies an existing staff member
func (sr *StaffRepository) Update(id int, name, email, phone, role string) (*Staff, error) {
	now := time.Now()
	staff := &Staff{}
	err := sr.db.QueryRow(
		"UPDATE staff SET name = $1, email = $2, phone = $3, role = $4, updated_at = $5 WHERE id = $6 RETURNING id, name, email, phone, role, created_at, updated_at",
		name, email, phone, role, now, id,
	).Scan(&staff.ID, &staff.Name, &staff.Email, &staff.Phone, &staff.Role, &staff.CreatedAt, &staff.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return staff, nil
}

// GetAll retrieves all staff members
func (sr *StaffRepository) GetAll() ([]Staff, error) {
	rows, err := sr.db.Query(
		"SELECT id, name, email, phone, role, created_at, updated_at FROM staff ORDER BY name",
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

// GetByID retrieves a single staff member by ID
func (sr *StaffRepository) GetByID(id int) (*Staff, error) {
	s := &Staff{}
	err := sr.db.QueryRow(
		"SELECT id, name, email, phone, role, created_at, updated_at FROM staff WHERE id = $1",
		id,
	).Scan(&s.ID, &s.Name, &s.Email, &s.Phone, &s.Role, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

// GetByEmail retrieves a staff member by email
func (sr *StaffRepository) GetByEmail(email string) (*Staff, error) {
	s := &Staff{}
	err := sr.db.QueryRow(
		"SELECT id, name, email, phone, role, created_at, updated_at FROM staff WHERE email = $1",
		email,
	).Scan(&s.ID, &s.Name, &s.Email, &s.Phone, &s.Role, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

// Delete removes a staff member
func (sr *StaffRepository) Delete(id int) error {
	_, err := sr.db.Exec("DELETE FROM staff WHERE id = $1", id)
	return err
}
