package handler

import (
	"jewellery-billing/internal/apiresponse"
	"jewellery-billing/internal/service"

	"github.com/gofiber/fiber/v2"
)

// AnalyticsHandler exposes HTTP endpoints for the dashboard.
type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

// GetDashboard returns the aggregated dashboard metrics and trends.
// GET /api/analytics/dashboard
func (h *AnalyticsHandler) GetDashboard(c *fiber.Ctx) error {
	data, err := h.analyticsService.GetDashboard(c.Context())
	if err != nil {
		return apiresponse.InternalError(c, "Failed to load dashboard analytics")
	}

	return apiresponse.Success(c, fiber.StatusOK, data)
}
