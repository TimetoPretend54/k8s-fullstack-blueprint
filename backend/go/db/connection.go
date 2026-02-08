package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// DB is the global database connection pool
var DB *sql.DB

// Connect establishes a connection to the PostgreSQL database
func Connect() (*sql.DB, error) {
	// Get database connection parameters from environment
	host := getEnv("DB_HOST", "")
	port := getEnv("DB_PORT", "")
	user := getEnv("DB_USER", "")
	password := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "")
	sslMode := getEnv("DB_SSL_MODE", "")

	// Validate required parameters
	if host == "" || port == "" || user == "" || password == "" || dbName == "" {
		return nil, errors.New("all database connection parameters (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME) must be set")
	}

	// Construct database URL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbName, sslMode)

	log.Printf("Connecting to database: %s", maskPassword(dbURL))

	// Open connection
	db, err := sql.Open("postgres", dbURL)
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

	log.Println("Database connection established successfully")
	DB = db
	return db, nil
}

// maskPassword hides the password in the database URL for logging
func maskPassword(dbURL string) string {
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		// If parsing fails, return a generic message
		return "postgres://<user>:<password>@<host>:<port>/<database>?sslmode=disable"
	}

	// Replace password with asterisks
	if parsedURL.User != nil {
		// Get the username
		username := parsedURL.User.Username()
		// Reconstruct URL with masked password
		masked := &url.URL{
			Scheme:   parsedURL.Scheme,
			User:     url.User(username), // no password
			Host:     parsedURL.Host,
			Path:     parsedURL.Path,
			RawQuery: parsedURL.RawQuery,
			Fragment: parsedURL.Fragment,
		}
		return masked.String()
	}

	return dbURL
}

// InitSchema creates necessary tables if they don't exist
// TODO/WARNING: This is a temporary scaffold solution. For production,
// use proper database migrations (e.g., goose, golang-migrate) to manage schema changes.
func InitSchema() error {
	log.Println("Initializing database schema...")

	// Create demo_data table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS demo_data (
		id SERIAL PRIMARY KEY,
		content TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create demo_data table: %w", err)
	}

	log.Println("Database schema initialized successfully (demo_data table ready)")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
