package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"jewellery-billing/internal/domain"
)

// tokenRepository is the PostgreSQL implementation of domain.TokenRepository.
type tokenRepository struct {
	db *pgxpool.Pool
}

// NewTokenRepository returns a production-ready TokenRepository backed by pgx.
func NewTokenRepository(db *pgxpool.Pool) domain.TokenRepository {
	return &tokenRepository{db: db}
}

// ── Email Verification Tokens ──────────────────────────────────────────

func (r *tokenRepository) CreateEmailVerification(ctx context.Context, token *domain.EmailVerificationToken) error {
	query := `
		INSERT INTO email_verification_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`
	return r.db.QueryRow(ctx, query,
		token.UserID, token.Token, token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)
}

func (r *tokenRepository) GetEmailVerificationByToken(ctx context.Context, tokenStr string) (*domain.EmailVerificationToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM email_verification_tokens
		WHERE token = $1`
	token := &domain.EmailVerificationToken{}
	err := r.db.QueryRow(ctx, query, tokenStr).Scan(
		&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("verification token not found")
		}
		return nil, fmt.Errorf("failed to get verification token: %w", err)
	}
	return token, nil
}

func (r *tokenRepository) DeleteEmailVerificationByUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		"DELETE FROM email_verification_tokens WHERE user_id = $1", userID)
	return err
}

// ── Password Reset Tokens ──────────────────────────────────────────────

func (r *tokenRepository) CreatePasswordReset(ctx context.Context, token *domain.PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`
	return r.db.QueryRow(ctx, query,
		token.UserID, token.Token, token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)
}

func (r *tokenRepository) GetPasswordResetByToken(ctx context.Context, tokenStr string) (*domain.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, used, created_at
		FROM password_reset_tokens
		WHERE token = $1`
	token := &domain.PasswordResetToken{}
	err := r.db.QueryRow(ctx, query, tokenStr).Scan(
		&token.ID, &token.UserID, &token.Token, &token.ExpiresAt,
		&token.Used, &token.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("reset token not found")
		}
		return nil, fmt.Errorf("failed to get reset token: %w", err)
	}
	return token, nil
}

func (r *tokenRepository) MarkPasswordResetUsed(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		"UPDATE password_reset_tokens SET used = true WHERE id = $1", id)
	return err
}

func (r *tokenRepository) DeletePasswordResetByUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		"DELETE FROM password_reset_tokens WHERE user_id = $1", userID)
	return err
}
