package handler

import (
	"jewellery-billing/internal/apiresponse"
	"jewellery-billing/internal/domain"
	"jewellery-billing/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthHandler exposes HTTP endpoints for authentication.
type AuthHandler struct {
	authService *service.AuthService
	orgRepo     domain.OrganizationRepository
}

func NewAuthHandler(authService *service.AuthService, orgRepo domain.OrganizationRepository) *AuthHandler {
	return &AuthHandler{authService: authService, orgRepo: orgRepo}
}

// Login godoc
// POST /api/auth/login
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
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{
		"message": "Logged out successfully",
	})
}

// Me godoc
// GET /api/auth/me (protected)
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(uuid.UUID)
	orgID, _ := c.Locals("organizationID").(uuid.UUID)
	email, _ := c.Locals("email").(string)
	role, _ := c.Locals("role").(domain.UserRole)

	// Fetch organization details
	org, _ := h.orgRepo.GetByID(c.Context(), orgID)
	var orgResponse *domain.OrganizationResponse
	if org != nil {
		resp := org.ToResponse()
		orgResponse = &resp
	}

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{
		"id":              userID,
		"organization_id": orgID,
		"email":           email,
		"role":            role,
		"organization":    orgResponse,
	})
}
