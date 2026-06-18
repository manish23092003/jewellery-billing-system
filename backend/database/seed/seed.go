package seed

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// SeedAdminUser creates a default admin account if no admin exists yet.
// This runs at every startup but is idempotent — it skips if an admin
// is already present.
//
// Default credentials:
//
//	Email:    admin@jewellery.com
//	Password: Admin@123
func SeedAdminUser(pool *pgxpool.Pool) error {
	ctx := context.Background()

	// Check whether any admin user already exists.
	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Info().Msg("Admin user already exists — skipping seed")
		return nil
	}

	// Hash the default password at runtime (never store pre-computed hashes in code).
	hash, err := bcrypt.GenerateFromPassword([]byte("Admin@123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx,
		`INSERT INTO users (name, email, password_hash, role)
		 VALUES ($1, $2, $3, $4)`,
		"Administrator", "admin@jewellery.com", string(hash), "admin",
	)
	if err != nil {
		return err
	}

	log.Info().Msg("✓ Default admin user created (admin@jewellery.com / Admin@123)")
	return nil
}
