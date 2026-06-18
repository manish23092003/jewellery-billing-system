package handler

import (
	"fmt"
	"strconv"

	"jewellery-billing/internal/domain"
	"jewellery-billing/internal/service"
	"jewellery-billing/internal/apiresponse"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// MetalRateHandler exposes HTTP endpoints for metal rate management.
type MetalRateHandler struct {
	rateService *service.MetalRateService
}

func NewMetalRateHandler(rateService *service.MetalRateService) *MetalRateHandler {
	return &MetalRateHandler{rateService: rateService}
}

// GetCurrentRates returns the latest rate for each metal/purity pair.
// GET /api/metal-rates
func (h *MetalRateHandler) GetCurrentRates(c *fiber.Ctx) error {
	rates, err := h.rateService.GetCurrentRates(c.Context())
	if err != nil {
		fmt.Printf("GetCurrentRates error: %v\n", err)
		return apiresponse.InternalError(c, "Failed to fetch current rates")
	}
	return apiresponse.Success(c, fiber.StatusOK, rates)
}

// GetHistory returns paginated rate history.
// GET /api/metal-rates/history?metal_type=gold&page=1&per_page=20
func (h *MetalRateHandler) GetHistory(c *fiber.Ctx) error {
	metalType := c.Query("metal_type", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))

	rates, total, err := h.rateService.GetHistory(c.Context(), metalType, page, perPage)
	if err != nil {
		return apiresponse.InternalError(c, "Failed to fetch rate history")
	}

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	return apiresponse.SuccessWithMeta(c, fiber.StatusOK, rates, &apiresponse.Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// Create adds a new metal rate.
// POST /api/metal-rates (admin only)
func (h *MetalRateHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateMetalRateRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	rate, err := h.rateService.Create(c.Context(), req)
	if err != nil {
		return apiresponse.BadRequest(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusCreated, rate)
}

// Update modifies an existing rate.
// PUT /api/metal-rates/:id (admin only)
func (h *MetalRateHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid rate ID")
	}

	var req domain.UpdateMetalRateRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	rate, err := h.rateService.Update(c.Context(), id, req)
	if err != nil {
		return apiresponse.NotFound(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, rate)
}

// Delete removes a metal rate entry.
// DELETE /api/metal-rates/:id (admin only)
func (h *MetalRateHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid rate ID")
	}

	if err := h.rateService.Delete(c.Context(), id); err != nil {
		return apiresponse.NotFound(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{"message": "Rate deleted"})
}

