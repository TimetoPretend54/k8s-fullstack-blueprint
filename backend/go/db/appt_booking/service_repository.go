package appt_booking

import (
	"database/sql"
	"time"
)

// ServiceRepository handles database operations for services
type ServiceRepository struct {
	db *sql.DB
}

// NewServiceRepository creates a new service repository
func NewServiceRepository(db *sql.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

// Create inserts a new service
func (sr *ServiceRepository) Create(name, description string, duration, priceCents int) (*Service, error) {
	now := time.Now()
	service := &Service{}
	err := sr.db.QueryRow(
		"INSERT INTO services (name, description, duration_minutes, price_cents, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, description, duration_minutes, price_cents, created_at, updated_at",
		name, description, duration, priceCents, now, now,
	).Scan(&service.ID, &service.Name, &service.Description, &service.DurationMin, &service.PriceCents, &service.CreatedAt, &service.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return service, nil
}

// Update modifies an existing service
func (sr *ServiceRepository) Update(id int, name, description string, duration, priceCents int) (*Service, error) {
	now := time.Now()
	service := &Service{}
	err := sr.db.QueryRow(
		"UPDATE services SET name = $1, description = $2, duration_minutes = $3, price_cents = $4, updated_at = $5 WHERE id = $6 RETURNING id, name, description, duration_minutes, price_cents, created_at, updated_at",
		name, description, duration, priceCents, now, id,
	).Scan(&service.ID, &service.Name, &service.Description, &service.DurationMin, &service.PriceCents, &service.CreatedAt, &service.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return service, nil
}

// GetAll retrieves all services
func (sr *ServiceRepository) GetAll() ([]Service, error) {
	rows, err := sr.db.Query(
		"SELECT id, name, description, duration_minutes, price_cents, created_at, updated_at FROM services ORDER BY name",
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

// GetByID retrieves a single service by ID
func (sr *ServiceRepository) GetByID(id int) (*Service, error) {
	s := &Service{}
	err := sr.db.QueryRow(
		"SELECT id, name, description, duration_minutes, price_cents, created_at, updated_at FROM services WHERE id = $1",
		id,
	).Scan(&s.ID, &s.Name, &s.Description, &s.DurationMin, &s.PriceCents, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

// Delete removes a service
func (sr *ServiceRepository) Delete(id int) error {
	_, err := sr.db.Exec("DELETE FROM services WHERE id = $1", id)
	return err
}
