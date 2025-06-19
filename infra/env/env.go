package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	KeyDBHost = "DB_HOST"
	KeyDBPort = "DB_PORT"
	KeyDBUser = "DB_USER"
	KeyDBPass = "DB_PASS"
	KeyDBName = "DB_NAME"

	KeyDBMigrationPath = "MIGRATIONS_PATH"
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
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	DBMigrationPath string

	DBAddr string
	DBDSN  string
)

// MustLoadEnv loads environment variables from .env.example and .env files.
// It sets defaults for any missing values and configures database connection parameters.
// If both .env.example and .env exist, .env values take precedence.
func MustLoadEnv() {
	_ = godotenv.Load(".env.example")
	_ = godotenv.Load(".env")

	DBHost = orEnv(KeyDBHost, defaultDBHost)
	DBPort = orEnv(KeyDBPort, defaultDBPort)
	DBUser = orEnv(KeyDBUser, defaultDBUser)
	DBPass = orEnv(KeyDBPass, defaultDBPass)
	DBName = orEnv(KeyDBName, defaultDBName)
	DBMigrationPath = orEnv(KeyDBMigrationPath, defaultDBMigrationPath)

	DBAddr = fmt.Sprintf("%s:%s", DBHost, DBPort)
	DBDSN = fmt.Sprintf("clickhouse://%s:%s@%s/%s", DBUser, DBPass, DBAddr, DBName)
}

// orEnv retrieves an environment variable value or uses a default if empty.
// It logs a warning when falling back to the default value.
//
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
