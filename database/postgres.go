package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func getConnectionString() string {
	// Load environment variables from .env file in development
	if env := os.Getenv("ENVIRONMENT"); env != "production" {
		log.Println("[ENV]:", env)

		err := godotenv.Load()
		if err != nil {
			log.Fatal("[FATAL]: failed to connect database", err)
		}
	} else {
		log.Println("[ENV]: production")
	}

	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
	)

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname,
	)
}

func runMigrations(dbURL, migrationsPath string) error {
	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	log.Printf("Migrations applied successfully. Current version: %d, dirty: %v", version, dirty)
	return nil
}

func ConnectDatabase() (*sql.DB, error) {
	connectionString := getConnectionString()

	// Open a new database connection
	database, err := sql.Open("postgres", getConnectionString())

	// Return an error if the connection failed
	if err != nil {
		msg := "invalid connection string provided"
		log.Println("[FAIL]:", msg)

		return nil, fmt.Errorf("%v", msg)
	}

	// Ping the database to verify the connection
	err = database.Ping()
	if err != nil {
		msg := "failed to establish a connection to the database"
		log.Println("[FAIL]:", msg)

		return nil, fmt.Errorf("%v", msg)
	}

	if err := runMigrations(connectionString, "migrations"); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return database, nil
}
