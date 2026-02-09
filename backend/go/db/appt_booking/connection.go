package appt_booking

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Connect establishes a connection to the appointment_booking database.
// It reuses the main database connection parameters (host, port, user, password, sslmode)
// but connects to a separate database specified by APPT_BOOKING_DB_NAME (or default "appt_booking").
// It automatically creates the database if it doesn't exist.
func Connect() (*sql.DB, error) {
	// Get database connection parameters from environment.
	// For host, port, user, password, sslmode: use same vars as main DB.
	// For database name: use APPT_BOOKING_DB_NAME, fallback to "appt_booking".
	host := getEnv("DB_HOST", "")
	port := getEnv("DB_PORT", "")
	user := getEnv("DB_USER", "")
	password := getEnv("DB_PASSWORD", "")
	dbName := getEnv("APPT_BOOKING_DB_NAME", "appt_booking")
	sslMode := getEnv("DB_SSL_MODE", "disable")

	// Validate required parameters
	if host == "" || port == "" || user == "" || password == "" {
		return nil, errors.New("all database connection parameters (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD) must be set")
	}

	// First, connect to the default 'postgres' database to create the target database if needed
	postgresURL := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=%s",
		user, password, host, port, sslMode)
	log.Printf("Connecting to system database (postgres) to ensure target database exists...")
	
	adminDB, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to system database: %w", err)
	}
	defer adminDB.Close()

	// Configure connection pool for admin connection
	adminDB.SetMaxOpenConns(5)
	adminDB.SetMaxIdleConns(5)
	adminDB.SetConnMaxLifetime(5 * time.Minute)

	// Test admin connection
	if err := adminDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping system database: %w", err)
	}

	// Check if target database exists, create if not
	var exists bool
	err = adminDB.QueryRow("SELECT 1 FROM pg_database WHERE datname = $1", dbName).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check if database exists: %w", err)
	}
	
	if !exists {
		log.Printf("Database '%s' does not exist, creating it...", dbName)
		_, err = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Database '%s' created successfully", dbName)
	} else {
		log.Printf("Database '%s' already exists", dbName)
	}

	// Now connect to the target database
	targetDBURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbName, sslMode)
	log.Printf("Connecting to appointment booking database: %s", maskPassword(targetDBURL))

	db, err := sql.Open("postgres", targetDBURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Appointment booking database connection established successfully")
	return db, nil
}

// maskPassword hides the password in the database URL for logging
func maskPassword(dbURL string) string {
	// Simple masking: replace password between : and @
	// For more robust masking, parse URL. But this is sufficient for logs.
	for i := 0; i < len(dbURL); i++ {
		if dbURL[i] == ':' && i+1 < len(dbURL) && dbURL[i+1] != '/' {
			// Find the @ after the password
			for j := i + 1; j < len(dbURL); j++ {
				if dbURL[j] == '@' {
					return dbURL[:i+1] + "****" + dbURL[j:]
				}
			}
			break
		}
	}
	return dbURL
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// InitSchema creates all necessary tables for the appointment booking feature if they don't exist.
// This is a temporary scaffold solution. For production, use proper database migrations.
// TODO: Replace with proper migration tool (e.g., golang-migrate, goose) for versioned, repeatable migrations.
// WARNING: Current approach is not suitable for production deployments without migration strategy.
func InitSchema(db *sql.DB) error {
	log.Println("Initializing appointment booking database schema...")

	// Create services table
	createServicesTable := `
	CREATE TABLE IF NOT EXISTS services (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		duration_minutes INTEGER NOT NULL CHECK (duration_minutes > 0),
		price_cents INTEGER NOT NULL CHECK (price_cents >= 0),
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`

	_, err := db.Exec(createServicesTable)
	if err != nil {
		return fmt.Errorf("failed to create services table: %w", err)
	}

	// Create staff table
	createStaffTable := `
	CREATE TABLE IF NOT EXISTS staff (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		phone VARCHAR(50),
		role VARCHAR(100) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`

	_, err = db.Exec(createStaffTable)
	if err != nil {
		return fmt.Errorf("failed to create staff table: %w", err)
	}

	// Create staff_services junction table
	createStaffServicesTable := `
	CREATE TABLE IF NOT EXISTS staff_services (
		staff_id INTEGER NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
		service_id INTEGER NOT NULL REFERENCES services(id) ON DELETE CASCADE,
		PRIMARY KEY (staff_id, service_id)
	)`

	_, err = db.Exec(createStaffServicesTable)
	if err != nil {
		return fmt.Errorf("failed to create staff_services table: %w", err)
	}

	// Create schedules table
	createSchedulesTable := `
	CREATE TABLE IF NOT EXISTS schedules (
		id SERIAL PRIMARY KEY,
		staff_id INTEGER NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
		day_of_week INTEGER NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
		start_time TIME NOT NULL,
		end_time TIME NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		UNIQUE(staff_id, day_of_week, start_time, end_time)
	)`

	_, err = db.Exec(createSchedulesTable)
	if err != nil {
		return fmt.Errorf("failed to create schedules table: %w", err)
	}

	// Create appointments table
	createAppointmentsTable := `
	CREATE TABLE IF NOT EXISTS appointments (
		id SERIAL PRIMARY KEY,
		customer_name VARCHAR(255) NOT NULL,
		customer_email VARCHAR(255) NOT NULL,
		customer_phone VARCHAR(50),
		staff_id INTEGER NOT NULL REFERENCES staff(id),
		service_id INTEGER NOT NULL REFERENCES services(id),
		appointment_datetime TIMESTAMP NOT NULL,
		duration_minutes INTEGER NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'confirmed',
		notes TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`

	_, err = db.Exec(createAppointmentsTable)
	if err != nil {
		return fmt.Errorf("failed to create appointments table: %w", err)
	}

	// Create indexes for better query performance
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_appt_booking_appointments_staff_datetime ON appointments(staff_id, appointment_datetime)
	`)
	if err != nil {
		return fmt.Errorf("failed to create index on appointments(staff_id, appointment_datetime): %w", err)
	}

	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_appt_booking_appointments_customer_email ON appointments(customer_email)
	`)
	if err != nil {
		return fmt.Errorf("failed to create index on appointments(customer_email): %w", err)
	}

	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_appt_booking_schedules_staff_day ON schedules(staff_id, day_of_week)
	`)
	if err != nil {
		return fmt.Errorf("failed to create index on schedules(staff_id, day_of_week): %w", err)
	}

	log.Println("Appointment booking database schema initialized successfully")
	return nil
}

// SeedSampleData inserts preloaded sample data for testing and demonstration.
// This is idempotent - can be safely called multiple times.
func SeedSampleData(db *sql.DB) error {
	log.Println("Seeding sample data for appointment booking...")

	// Check if data already exists to avoid duplicates
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM services").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing services: %w", err)
	}
	if count > 0 {
		log.Println("Sample data already exists, skipping seeding")
		return nil
	}

	// Insert services: Haircut (30 min, $35), Beard Trim (15 min, $20), Full Grooming (60 min, $60)
	services := []struct {
		name        string
		description string
		duration    int
		price       int
	}{
		{"Haircut", "Classic haircut with styling", 30, 3500},
		{"Beard Trim", "Beard trimming and shaping", 15, 2000},
		{"Full Grooming", "Complete haircut and beard grooming", 60, 6000},
	}

	serviceIDs := make(map[string]int)
	for _, s := range services {
		var id int
		err := db.QueryRow(
			"INSERT INTO services (name, description, duration_minutes, price_cents) VALUES ($1, $2, $3, $4) RETURNING id",
			s.name, s.description, s.duration, s.price,
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("failed to insert service %s: %w", s.name, err)
		}
		serviceIDs[s.name] = id
		log.Printf("  Inserted service: %s (ID: %d)", s.name, id)
	}

	// Insert staff: John Smith, Jane Doe, Admin User
	staff := []struct {
		name  string
		email string
		phone string
		role  string
	}{
		{"John Smith", "john@example.com", "555-0101", "provider"},
		{"Jane Doe", "jane@example.com", "555-0102", "provider"},
		{"Admin User", "admin@example.com", "555-0100", "admin"},
	}

	staffIDs := make(map[string]int)
	for _, s := range staff {
		var id int
		err := db.QueryRow(
			"INSERT INTO staff (name, email, phone, role) VALUES ($1, $2, $3, $4) RETURNING id",
			s.name, s.email, s.phone, s.role,
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("failed to insert staff %s: %w", s.name, err)
		}
		staffIDs[s.name] = id
		log.Printf("  Inserted staff: %s (ID: %d)", s.name, id)
	}

	// Insert staff-service assignments
	// John: Haircut, Beard Trim
	// Jane: Haircut, Full Grooming
	assignments := []struct {
		staffName  string
		serviceName string
	}{
		{"John Smith", "Haircut"},
		{"John Smith", "Beard Trim"},
		{"Jane Doe", "Haircut"},
		{"Jane Doe", "Full Grooming"},
	}

	for _, a := range assignments {
		_, err := db.Exec(
			"INSERT INTO staff_services (staff_id, service_id) VALUES ($1, $2)",
			staffIDs[a.staffName], serviceIDs[a.serviceName],
		)
		if err != nil {
			// Ignore duplicate key errors
			if !isDuplicateKeyError(err) {
				return fmt.Errorf("failed to assign service %s to staff %s: %w", a.serviceName, a.staffName, err)
			}
		}
		log.Printf("  Assigned %s to %s", a.serviceName, a.staffName)
	}

	// Insert schedules
	// John: Mon-Fri 9:00-17:00, Sat 10:00-14:00
	// Jane: Mon-Fri 10:00-18:00
	schedules := []struct {
		staffName string
		day       int    // 0=Sun, 1=Mon, ..., 6=Sat
		start     string // HH:MM
		end       string // HH:MM
	}{
		{"John Smith", 1, "09:00", "17:00"},
		{"John Smith", 2, "09:00", "17:00"},
		{"John Smith", 3, "09:00", "17:00"},
		{"John Smith", 4, "09:00", "17:00"},
		{"John Smith", 5, "09:00", "17:00"},
		{"John Smith", 6, "10:00", "14:00"},
		{"Jane Doe", 1, "10:00", "18:00"},
		{"Jane Doe", 2, "10:00", "18:00"},
		{"Jane Doe", 3, "10:00", "18:00"},
		{"Jane Doe", 4, "10:00", "18:00"},
		{"Jane Doe", 5, "10:00", "18:00"},
	}

	for _, s := range schedules {
		_, err := db.Exec(
			"INSERT INTO schedules (staff_id, day_of_week, start_time, end_time) VALUES ($1, $2, $3, $4)",
			staffIDs[s.staffName], s.day, s.start, s.end,
		)
		if err != nil {
			// Ignore duplicate key errors
			if !isDuplicateKeyError(err) {
				return fmt.Errorf("failed to insert schedule for %s day %d: %w", s.staffName, s.day, err)
			}
		}
		log.Printf("  Inserted schedule: %s day %d %s-%s", s.staffName, s.day, s.start, s.end)
	}

	// Insert some sample appointments (upcoming)
	// Use a helper to create appointments
	createAppointment := func(staffName string, serviceName string, year, month, day, hour, minute int, duration int, status string) error {
		staffID := staffIDs[staffName]
		serviceID := serviceIDs[serviceName]
		datetime := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local)
		_, err := db.Exec(
			"INSERT INTO appointments (customer_name, customer_email, customer_phone, staff_id, service_id, appointment_datetime, duration_minutes, status, notes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
			"Sample Customer", "customer@example.com", "555-1234", staffID, serviceID, datetime, duration, status, "",
		)
		return err
	}

	// Create a few appointments for John and Jane
	now := time.Now()
	year, month, day := now.Year(), int(now.Month()), now.Day()
	_ = createAppointment("John Smith", "Haircut", year, month, day+1, 9, 30, 30, "confirmed")
	_ = createAppointment("John Smith", "Beard Trim", year, month, day+2, 14, 0, 15, "confirmed")
	_ = createAppointment("Jane Doe", "Full Grooming", year, month, day+3, 11, 0, 60, "confirmed")
	_ = createAppointment("Jane Doe", "Haircut", year, month, day+4, 13, 30, 30, "completed")
	_ = createAppointment("John Smith", "Haircut", year, month, day+5, 10, 0, 30, "confirmed")

	log.Println("Sample data seeding completed successfully")
	return nil
}

// isDuplicateKeyError checks if an error is a PostgreSQL duplicate key violation
func isDuplicateKeyError(err error) bool {
	// PostgreSQL error code for unique_violation is 23505
	// This is a simple check; for more robust handling, use pgerror package
	if err == nil {
		return false
	}
	// Check if error message contains "duplicate key" or "unique constraint"
	// This is a heuristic; proper way is to check sql.State
	if err.Error()[0:23] == "ERROR: duplicate key value" {
		return true
	}
	return false
}
