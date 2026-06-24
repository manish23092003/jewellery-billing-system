package handler

import (
	"jewellery-billing/internal/apiresponse"
	"jewellery-billing/internal/domain"
	"jewellery-billing/internal/service"

	"github.com/gofiber/fiber/v2"
)

// RegistrationHandler handles business registration endpoints.
type RegistrationHandler struct {
	registrationService *service.RegistrationService
}

func NewRegistrationHandler(registrationService *service.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{registrationService: registrationService}
}

// Register godoc
// POST /api/auth/register
// Creates a new organization + admin user and returns JWT tokens.
func (h *RegistrationHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	result, err := h.registrationService.Register(c.Context(), req)
	if err != nil {
		return apiresponse.BadRequest(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusCreated, result)
}
