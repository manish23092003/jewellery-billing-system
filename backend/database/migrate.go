package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // file source
	"github.com/rs/zerolog/log"
)

// RunMigrations applies all pending up-migrations from the database/migrations
// directory. It is safe to call on every startup — already-applied migrations
// are skipped automatically (ErrNoChange).
func RunMigrations(databaseURL string) error {
	m, err := migrate.New(
		"file://database/migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	version, dirty, _ := m.Version()
	log.Info().
		Uint("version", version).
		Bool("dirty", dirty).
		Msg("✓ Database migrations completed")

	return nil
}
