package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"jewellery-billing/internal/domain"
)

// PasswordResetService handles forgot/reset password flows.
type PasswordResetService struct {
	userRepo    domain.UserRepository
	tokenRepo   domain.TokenRepository
	emailSender EmailSender
}

// NewPasswordResetService creates a PasswordResetService with required dependencies.
func NewPasswordResetService(
	userRepo domain.UserRepository,
	tokenRepo domain.TokenRepository,
	emailSender EmailSender,
) *PasswordResetService {
	return &PasswordResetService{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		emailSender: emailSender,
	}
}

// ForgotPassword generates a password reset token and sends the reset email.
// Returns nil even if user doesn't exist (prevents email enumeration).
func (s *PasswordResetService) ForgotPassword(ctx context.Context, req domain.ForgotPasswordRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Don't reveal whether the email exists — just return success
		return nil
	}

	// Delete any existing reset tokens for this user
	_ = s.tokenRepo.DeletePasswordResetByUser(ctx, user.ID)

	// Generate new token
	tokenStr, err := generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	resetToken := &domain.PasswordResetToken{
		UserID:    user.ID,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := s.tokenRepo.CreatePasswordReset(ctx, resetToken); err != nil {
		return fmt.Errorf("failed to store reset token: %w", err)
	}

	// Send the reset email
	if err := s.emailSender.SendPasswordResetEmail(user.Email, user.Name, tokenStr); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

// ResetPassword validates the token and sets the new password.
func (s *PasswordResetService) ResetPassword(ctx context.Context, req domain.ResetPasswordRequest) error {
	if req.Token == "" {
		return fmt.Errorf("token is required")
	}
	if len(req.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if req.Password != req.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	// Find the token
	resetToken, err := s.tokenRepo.GetPasswordResetByToken(ctx, req.Token)
	if err != nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	// Check if token is expired
	if time.Now().After(resetToken.ExpiresAt) {
		return fmt.Errorf("reset token has expired")
	}

	// Check if token was already used
	if resetToken.Used {
		return fmt.Errorf("reset token has already been used")
	}

	// Fetch the user
	user, err := s.userRepo.GetByID(ctx, resetToken.UserID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Hash the new password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update the password
	user.PasswordHash = string(hash)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	_ = s.tokenRepo.MarkPasswordResetUsed(ctx, resetToken.ID)

	return nil
}
