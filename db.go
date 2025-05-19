package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse" // ClickHouse driver for migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"         // File source driver for migrate
)

// MustLoadClickHouseConn creates a ClickHouse connection or terminates the program if connection fails.
// It uses environment variables for connection parameters and performs validation checks.
// Returns:
//   - driver.Conn: An established ClickHouse connection
//
// Exits the program if the connection cannot be established.
func MustLoadClickHouseConn() driver.Conn {
	var opts = dbOptions()

	conn, err := LoadClickHouseConn(opts)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func dbOptions() *clickhouse.Options {
	return &clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", dbHost, dbPort)},
		Auth: clickhouse.Auth{
			Database: dbName,
			Username: dbUser,
			Password: dbPass,
		},
	}
}

// LoadClickHouseConn establishes a connection to the ClickHouse database using the provided options.
// It verifies connectivity with a ping test and returns the connection or an error.
// Returns:
//   - driver.Conn: Active ClickHouse connection on success
//   - error: Error details if connection fails
func LoadClickHouseConn(opts *clickhouse.Options) (driver.Conn, error) {
	conn, err := clickhouse.Open(opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}

// LoadClickHouseDB establishes a standard database/sql connection to ClickHouse.
// This supports SQL operations using the standard library interface.
// Returns:
//   - *sql.DB: Database connection object
//   - error: Connection error if any
func LoadClickHouseDB(opts *clickhouse.Options) (*sql.DB, error) {
	db := clickhouse.OpenDB(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

// Migrate applies all pending database migrations from the configured migrations directory.
// It tracks migration versions to ensure proper schema evolution and reports migration status.
// Returns:
//   - error: nil on success, detailed error on failure
//
// Note: No error is returned if no new migrations are available to apply.
func Migrate() error {
	if _, err := os.Stat(dbMigrationPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory not found: %s", dbMigrationPath)
	}

	m, err := migrate.New("file://"+dbMigrationPath, dbDSN)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance: %w", err)
	}

	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			log.Printf("Warning: Error closing migration source: %v", sourceErr)
		}
		if dbErr != nil {
			log.Printf("Warning: Error closing migration database connection: %v", dbErr)
		}
	}()

	currentVersion, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("failed to get current migration version: %w", err)
	}
	if dirty {
		log.Println("Warning: Database is in a dirty state. Consider manual intervention.")
		// You might want to force a version and then try to migrate up/down.
		// Example: if err := m.Force(PREVIOUS_CLEAN_VERSION); err != nil { ... }
		return fmt.Errorf("database is dirty. Current version: %d", currentVersion)
	}

	if errors.Is(err, migrate.ErrNilVersion) {
		log.Println("No migrations applied yet.")
	} else {
		log.Printf("Current migration version: %d, Dirty: %v\n", currentVersion, dirty)
	}

	log.Println("Applying migrations...")
	err = m.Up() // Apply all available up migrations
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No new migrations to apply.")
			return nil // Not an error if there's no change
		}
		return fmt.Errorf("an error occurred while applying migrations: %w", err)
	}

	log.Println("Migrations applied successfully!")

	finalVersion, dirty, err := m.Version()
	if err != nil {
		log.Printf("Warning: Could not get version after migration: %v", err)
	} else {
		log.Printf("New migration version: %d, Dirty: %v\n", finalVersion, dirty)
	}

	// Check for errors from source and database connections
	sourceErr, dbErr := m.Close()
	if sourceErr != nil {
		log.Printf("Warning: Error closing migration source: %v", sourceErr)
	}
	if dbErr != nil {
		log.Printf("Warning: Error closing migration database connection: %v", dbErr)
	}

	return nil
}
