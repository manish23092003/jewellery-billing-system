package handler

import (
	"jewellery-billing/internal/apiresponse"
	"jewellery-billing/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AnalyticsHandler exposes HTTP endpoints for dashboard analytics.
type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

// GetDashboard godoc
// GET /api/analytics/dashboard
func (h *AnalyticsHandler) GetDashboard(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)

	data, err := h.analyticsService.GetDashboard(c.Context(), orgID)
	if err != nil {
		return apiresponse.InternalError(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, data)
}
