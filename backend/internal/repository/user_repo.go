package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"jewellery-billing/internal/domain"
)

// userRepository is the PostgreSQL implementation of domain.UserRepository.
type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository returns a production-ready UserRepository backed by pgx.
func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &userRepository{db: db}
}

// ── Queries ────────────────────────────────────────────────────────────

const (
	queryInsertUser = `
		INSERT INTO users (organization_id, name, email, password_hash, role, is_active, email_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	queryGetUserByID = `
		SELECT id, organization_id, name, email, password_hash, role, is_active, email_verified, created_at, updated_at
		FROM users WHERE id = $1`

	queryGetUserByEmail = `
		SELECT id, organization_id, name, email, password_hash, role, is_active, email_verified, created_at, updated_at
		FROM users WHERE email = $1`

	queryGetAllUsersByOrg = `
		SELECT id, organization_id, name, email, password_hash, role, is_active, email_verified, created_at, updated_at
		FROM users WHERE organization_id = $1 ORDER BY created_at DESC`

	queryUpdateUser = `
		UPDATE users
		SET name = $2, email = $3, password_hash = $4, role = $5,
		    is_active = $6, email_verified = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	queryDeleteUser = `DELETE FROM users WHERE id = $1`

	queryCountUsersByOrg = `SELECT COUNT(*) FROM users WHERE organization_id = $1`
)

// ── Implementation ─────────────────────────────────────────────────────

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.QueryRow(ctx, queryInsertUser,
		user.OrganizationID, user.Name, user.Email, user.PasswordHash,
		user.Role, user.IsActive, user.EmailVerified,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow(ctx, queryGetUserByID, id).Scan(
		&user.ID, &user.OrganizationID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetByEmail is org-agnostic — used during login to find the user across all orgs.
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow(ctx, queryGetUserByEmail, email).Scan(
		&user.ID, &user.OrganizationID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetAllByOrg(ctx context.Context, orgID uuid.UUID) ([]domain.User, error) {
	rows, err := r.db.Query(ctx, queryGetAllUsersByOrg, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID, &user.OrganizationID, &user.Name, &user.Email, &user.PasswordHash,
			&user.Role, &user.IsActive, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	err := r.db.QueryRow(ctx, queryUpdateUser,
		user.ID, user.Name, user.Email, user.PasswordHash,
		user.Role, user.IsActive, user.EmailVerified,
	).Scan(&user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, queryDeleteUser, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *userRepository) CountByOrg(ctx context.Context, orgID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, queryCountUsersByOrg, orgID).Scan(&count)
	return count, err
}
