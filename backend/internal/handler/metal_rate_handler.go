package handler

import (
	"strconv"

	"jewellery-billing/internal/apiresponse"
	"jewellery-billing/internal/domain"
	"jewellery-billing/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// MetalRateHandler exposes HTTP endpoints for metal rate operations.
type MetalRateHandler struct {
	metalRateService *service.MetalRateService
}

func NewMetalRateHandler(metalRateService *service.MetalRateService) *MetalRateHandler {
	return &MetalRateHandler{metalRateService: metalRateService}
}

// Create godoc
// POST /api/metal-rates
func (h *MetalRateHandler) Create(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)

	var req domain.CreateMetalRateRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	rate, err := h.metalRateService.Create(c.Context(), orgID, req)
	if err != nil {
		return apiresponse.BadRequest(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusCreated, rate)
}

// GetCurrentRates godoc
// GET /api/metal-rates/latest
func (h *MetalRateHandler) GetCurrentRates(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)

	rates, err := h.metalRateService.GetCurrentRates(c.Context(), orgID)
	if err != nil {
		return apiresponse.InternalError(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, rates)
}

// GetHistory godoc
// GET /api/metal-rates
func (h *MetalRateHandler) GetHistory(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)

	metalType := c.Query("metal_type")
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	rates, total, err := h.metalRateService.GetHistory(c.Context(), orgID, metalType, limit, offset)
	if err != nil {
		return apiresponse.InternalError(c, err.Error())
	}

	return apiresponse.SuccessWithMeta(c, fiber.StatusOK, rates, &apiresponse.Meta{
		Total: total,
	})
}

// Delete godoc
// DELETE /api/metal-rates/:id
func (h *MetalRateHandler) Delete(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid rate ID")
	}

	if err := h.metalRateService.Delete(c.Context(), orgID, id); err != nil {
		return apiresponse.NotFound(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{"message": "Metal rate deleted"})
}
