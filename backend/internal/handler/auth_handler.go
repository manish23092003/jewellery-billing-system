package handler

import (
	"jewellery-billing/internal/domain"
	"jewellery-billing/internal/service"
	"jewellery-billing/internal/apiresponse"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthHandler exposes HTTP endpoints for authentication.
type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// POST /api/auth/login
// Accepts {email, password} and returns JWT token pair + user info.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	if req.Email == "" || req.Password == "" {
		return apiresponse.BadRequest(c, "Email and password are required")
	}

	result, err := h.authService.Login(c.Context(), req)
	if err != nil {
		return apiresponse.Unauthorized(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, result)
}

// RefreshToken godoc
// POST /api/auth/refresh
// Accepts {refresh_token} and returns a new token pair.
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req domain.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	if req.RefreshToken == "" {
		return apiresponse.BadRequest(c, "Refresh token is required")
	}

	result, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return apiresponse.Unauthorized(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, result)
}

// Logout godoc
// POST /api/auth/logout (protected)
// Client-side logout — the server acknowledges the request.
// A production system would blacklist the token in Redis.
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{
		"message": "Logged out successfully",
	})
}

// Me godoc
// GET /api/auth/me (protected)
// Returns the currently authenticated user's profile.
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(uuid.UUID)
	email, _ := c.Locals("email").(string)
	role, _ := c.Locals("role").(domain.UserRole)
	name, _ := c.Locals("name").(string)

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{
		"id":    userID,
		"name":  name,
		"email": email,
		"role":  role,
	})
}

