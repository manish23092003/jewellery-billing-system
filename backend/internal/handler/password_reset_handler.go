package handler

import (
	"jewellery-billing/internal/apiresponse"
	"jewellery-billing/internal/domain"
	"jewellery-billing/internal/service"

	"github.com/gofiber/fiber/v2"
)

// PasswordResetHandler handles password reset endpoints.
type PasswordResetHandler struct {
	passwordResetService *service.PasswordResetService
}

func NewPasswordResetHandler(passwordResetService *service.PasswordResetService) *PasswordResetHandler {
	return &PasswordResetHandler{passwordResetService: passwordResetService}
}

// ForgotPassword godoc
// POST /api/auth/forgot-password
func (h *PasswordResetHandler) ForgotPassword(c *fiber.Ctx) error {
	var req domain.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	if err := h.passwordResetService.ForgotPassword(c.Context(), req); err != nil {
		return apiresponse.BadRequest(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{
		"message": "If an account with that email exists, a password reset link has been sent.",
	})
}

// ResetPassword godoc
// POST /api/auth/reset-password
func (h *PasswordResetHandler) ResetPassword(c *fiber.Ctx) error {
	var req domain.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	if err := h.passwordResetService.ResetPassword(c.Context(), req); err != nil {
		return apiresponse.BadRequest(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{
		"message": "Password has been reset successfully.",
	})
}
