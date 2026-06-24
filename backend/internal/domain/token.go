package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ── Email Verification Token ──────────────────────────────────────────

// EmailVerificationToken is used to verify a user's email address after registration.
type EmailVerificationToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// ── Password Reset Token ──────────────────────────────────────────────

// PasswordResetToken is used for the "forgot password" flow.
type PasswordResetToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

// ── Request DTOs ──────────────────────────────────────────────────────

// ForgotPasswordRequest is the payload for requesting a password reset.
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest is the payload for resetting a password with a token.
type ResetPasswordRequest struct {
	Token           string `json:"token"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// VerifyEmailRequest is the payload for verifying an email address.
type VerifyEmailRequest struct {
	Token string `json:"token"`
}

// ── Repository Interface ──────────────────────────────────────────────

// TokenRepository defines the data-access contract for token persistence.
type TokenRepository interface {
	// Email verification tokens
	CreateEmailVerification(ctx context.Context, token *EmailVerificationToken) error
	GetEmailVerificationByToken(ctx context.Context, token string) (*EmailVerificationToken, error)
	DeleteEmailVerificationByUser(ctx context.Context, userID uuid.UUID) error

	// Password reset tokens
	CreatePasswordReset(ctx context.Context, token *PasswordResetToken) error
	GetPasswordResetByToken(ctx context.Context, token string) (*PasswordResetToken, error)
	MarkPasswordResetUsed(ctx context.Context, id uuid.UUID) error
	DeletePasswordResetByUser(ctx context.Context, userID uuid.UUID) error
}
