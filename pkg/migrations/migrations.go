package migrations

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

// Migrate applies migrations to the database
func Migrate(db *sql.DB) error {
	// Create a new migrate instance using the PostgreSQL driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrations driver: %v", err)
	}

	// Initialize migrations source (in this case, local directory migrations)
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"song", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}

	// Run migrations
	log.Println("Applying database migrations...")
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply.")
		} else {
			log.Fatalf("Migration error: %v", err) // Critical error
		}
	}

	log.Println("Migrations applied successfully!")
	return nil
}
