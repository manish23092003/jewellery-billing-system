package handler

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"jewellery-billing/internal/apiresponse"
	"jewellery-billing/internal/domain"
	"jewellery-billing/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// EmailVerificationHandler handles email verification endpoints.
type EmailVerificationHandler struct {
	authService *service.AuthService
	tokenRepo   domain.TokenRepository
	emailSender service.EmailSender
	userRepo    domain.UserRepository
}

func NewEmailVerificationHandler(
	authService *service.AuthService,
	tokenRepo domain.TokenRepository,
	emailSender service.EmailSender,
	userRepo domain.UserRepository,
) *EmailVerificationHandler {
	return &EmailVerificationHandler{
		authService: authService,
		tokenRepo:   tokenRepo,
		emailSender: emailSender,
		userRepo:    userRepo,
	}
}

// VerifyEmail godoc
// POST /api/auth/verify-email
func (h *EmailVerificationHandler) VerifyEmail(c *fiber.Ctx) error {
	var req domain.VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	if req.Token == "" {
		return apiresponse.BadRequest(c, "Token is required")
	}

	if err := h.authService.VerifyEmail(c.Context(), h.tokenRepo, req.Token); err != nil {
		return apiresponse.BadRequest(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{
		"message": "Email verified successfully.",
	})
}

// ResendVerification godoc
// POST /api/auth/resend-verification (protected)
func (h *EmailVerificationHandler) ResendVerification(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(uuid.UUID)

	user, err := h.userRepo.GetByID(c.Context(), userID)
	if err != nil {
		return apiresponse.NotFound(c, "User not found")
	}

	if user.EmailVerified {
		return apiresponse.BadRequest(c, "Email is already verified")
	}

	// Delete existing verification tokens
	_ = h.tokenRepo.DeleteEmailVerificationByUser(c.Context(), userID)

	// Generate new token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return apiresponse.InternalError(c, "Failed to generate verification token")
	}
	tokenStr := hex.EncodeToString(bytes)

	evToken := &domain.EmailVerificationToken{
		UserID:    userID,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := h.tokenRepo.CreateEmailVerification(c.Context(), evToken); err != nil {
		return apiresponse.InternalError(c, "Failed to create verification token")
	}

	_ = h.emailSender.SendVerificationEmail(user.Email, user.Name, tokenStr)

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{
		"message": "Verification email sent.",
	})
}
