package db

import (
	"database/sql"
	"time"
)

// DemoDataRepository handles database operations for demo data
type DemoDataRepository struct {
	db *sql.DB
}

// NewDemoDataRepository creates a new demo data repository
func NewDemoDataRepository(db *sql.DB) *DemoDataRepository {
	return &DemoDataRepository{db: db}
}

// Upsert inserts or updates a demo record
// If id is 0, it creates a new record; otherwise it updates the existing one
func (dr *DemoDataRepository) Upsert(id int, content string) (*DemoData, error) {
	now := time.Now()

	var data DemoData
	var err error

	if id == 0 {
		// Insert new record
		err = dr.db.QueryRow(
			"INSERT INTO demo_data (content, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id, content, created_at, updated_at",
			content, now, now,
		).Scan(&data.ID, &data.Content, &data.CreatedAt, &data.UpdatedAt)
	} else {
		// Update existing record
		err = dr.db.QueryRow(
			"UPDATE demo_data SET content = $1, updated_at = $2 WHERE id = $3 RETURNING id, content, created_at, updated_at",
			content, now, id,
		).Scan(&data.ID, &data.Content, &data.CreatedAt, &data.UpdatedAt)
	}

	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetAll retrieves all demo records
func (dr *DemoDataRepository) GetAll() ([]DemoData, error) {
	rows, err := dr.db.Query(
		"SELECT id, content, created_at, updated_at FROM demo_data ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []DemoData
	for rows.Next() {
		var data DemoData
		if err := rows.Scan(&data.ID, &data.Content, &data.CreatedAt, &data.UpdatedAt); err != nil {
			return nil, err
		}
		records = append(records, data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// GetByID retrieves a single record by ID
func (dr *DemoDataRepository) GetByID(id int) (*DemoData, error) {
	data := &DemoData{}
	err := dr.db.QueryRow(
		"SELECT id, content, created_at, updated_at FROM demo_data WHERE id = $1",
		id,
	).Scan(&data.ID, &data.Content, &data.CreatedAt, &data.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}

// Delete removes a record
func (dr *DemoDataRepository) Delete(id int) error {
	_, err := dr.db.Exec("DELETE FROM demo_data WHERE id = $1", id)
	return err
}
