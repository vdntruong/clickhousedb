package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	EnvKeyDBHost = "DB_HOST"
	EnvKeyDBPort = "DB_PORT"
	EnvKeyDBUser = "DB_USER"
	EnvKeyDBPass = "DB_PASS"
	EnvKeyDBName = "DB_NAME"

	EnvKeyDBMigrationPath = "MIGRATIONS_PATH"
)

const (
	defaultDBHost = "localhost"
	defaultDBPort = "9000"
	defaultDBUser = "admin"
	defaultDBPass = ""
	defaultDBName = "default"

	defaultDBMigrationPath = "migrations"
)

var (
	dbHost string
	dbPort string
	dbUser string
	dbPass string
	dbName string

	dbMigrationPath string

	dbAddr string
	dbDSN  string
)

// MustLoadEnv loads environment variables from .env.example and .env files.
// It sets defaults for any missing values and configures database connection parameters.
// If both .env.example and .env exist, .env values take precedence.
func MustLoadEnv() {
	_ = godotenv.Load(".env.example")
	_ = godotenv.Load(".env")

	dbHost = orEnv(EnvKeyDBHost, defaultDBHost)
	dbPort = orEnv(EnvKeyDBPort, defaultDBPort)
	dbUser = orEnv(EnvKeyDBUser, defaultDBUser)
	dbPass = orEnv(EnvKeyDBPass, defaultDBPass)
	dbName = orEnv(EnvKeyDBName, defaultDBName)
	dbMigrationPath = orEnv(EnvKeyDBMigrationPath, defaultDBMigrationPath)

	dbAddr = fmt.Sprintf("%s:%s", dbHost, dbPort)
	dbDSN = fmt.Sprintf("clickhouse://%s:%s@%s/%s", dbUser, dbPass, dbAddr, dbName)
}

// orEnv retrieves an environment variable value or uses a default if empty.
// It logs a warning when falling back to the default value.
// Parameters:
//   - envKey: Name of the environment variable to retrieve
//   - defaultValue: Value to use if environment variable is not set
//
// Returns:
//   - string: Value from environment or the default
func orEnv(envKey string, defaultValue string) string {
	var value = os.Getenv(envKey)
	if value == "" {
		value = defaultValue
		log.Printf("Warning: %s not set, using default: %s\n", envKey, defaultValue)
	}
	return value
}
