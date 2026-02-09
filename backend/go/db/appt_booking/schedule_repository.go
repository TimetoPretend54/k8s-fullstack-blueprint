package appt_booking

import (
	"database/sql"
	"time"
)

// ScheduleRepository handles database operations for schedules
type ScheduleRepository struct {
	db *sql.DB
}

// NewScheduleRepository creates a new schedule repository
func NewScheduleRepository(db *sql.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// Create inserts a new schedule
func (sr *ScheduleRepository) Create(staffID, dayOfWeek int, startTime, endTime string) (*Schedule, error) {
	now := time.Now()
	schedule := &Schedule{}
	err := sr.db.QueryRow(
		"INSERT INTO schedules (staff_id, day_of_week, start_time, end_time, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, staff_id, day_of_week, start_time, end_time, created_at, updated_at",
		staffID, dayOfWeek, startTime, endTime, now, now,
	).Scan(&schedule.ID, &schedule.StaffID, &schedule.DayOfWeek, &schedule.StartTime, &schedule.EndTime, &schedule.CreatedAt, &schedule.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return schedule, nil
}

// Update modifies an existing schedule
func (sr *ScheduleRepository) Update(id int, dayOfWeek int, startTime, endTime string) (*Schedule, error) {
	now := time.Now()
	schedule := &Schedule{}
	err := sr.db.QueryRow(
		"UPDATE schedules SET day_of_week = $1, start_time = $2, end_time = $3, updated_at = $4 WHERE id = $5 RETURNING id, staff_id, day_of_week, start_time, end_time, created_at, updated_at",
		dayOfWeek, startTime, endTime, now, id,
	).Scan(&schedule.ID, &schedule.StaffID, &schedule.DayOfWeek, &schedule.StartTime, &schedule.EndTime, &schedule.CreatedAt, &schedule.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return schedule, nil
}

// GetAll retrieves all schedules
func (sr *ScheduleRepository) GetAll() ([]Schedule, error) {
	rows, err := sr.db.Query(
		"SELECT id, staff_id, day_of_week, start_time, end_time, created_at, updated_at FROM schedules ORDER BY staff_id, day_of_week",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []Schedule
	for rows.Next() {
		var s Schedule
		if err := rows.Scan(&s.ID, &s.StaffID, &s.DayOfWeek, &s.StartTime, &s.EndTime, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return schedules, nil
}

// GetByID retrieves a single schedule by ID
func (sr *ScheduleRepository) GetByID(id int) (*Schedule, error) {
	s := &Schedule{}
	err := sr.db.QueryRow(
		"SELECT id, staff_id, day_of_week, start_time, end_time, created_at, updated_at FROM schedules WHERE id = $1",
		id,
	).Scan(&s.ID, &s.StaffID, &s.DayOfWeek, &s.StartTime, &s.EndTime, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

// GetByStaff retrieves all schedules for a specific staff member
func (sr *ScheduleRepository) GetByStaff(staffID int) ([]Schedule, error) {
	rows, err := sr.db.Query(
		"SELECT id, staff_id, day_of_week, start_time, end_time, created_at, updated_at FROM schedules WHERE staff_id = $1 ORDER BY day_of_week, start_time",
		staffID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []Schedule
	for rows.Next() {
		var s Schedule
		if err := rows.Scan(&s.ID, &s.StaffID, &s.DayOfWeek, &s.StartTime, &s.EndTime, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return schedules, nil
}

// Delete removes a schedule
func (sr *ScheduleRepository) Delete(id int) error {
	_, err := sr.db.Exec("DELETE FROM schedules WHERE id = $1", id)
	return err
}

// DeleteByStaff removes all schedules for a staff member
func (sr *ScheduleRepository) DeleteByStaff(staffID int) error {
	_, err := sr.db.Exec("DELETE FROM schedules WHERE staff_id = $1", staffID)
	return err
}
