package handler

import (
	"fmt"
	"jewellery-billing/internal/apiresponse"
	"jewellery-billing/internal/domain"
	"jewellery-billing/internal/service"

	"github.com/gofiber/fiber/v2"
)

// SettingHandler exposes endpoints for shop configuration.
type SettingHandler struct {
	settingService *service.SettingService
}

func NewSettingHandler(settingService *service.SettingService) *SettingHandler {
	return &SettingHandler{settingService: settingService}
}

// Get fetches the current shop settings.
// GET /api/settings
func (h *SettingHandler) Get(c *fiber.Ctx) error {
	settings, err := h.settingService.Get(c.Context())
	if err != nil {
		return apiresponse.InternalError(c, "Failed to load settings")
	}
	return apiresponse.Success(c, fiber.StatusOK, settings)
}

// Update modifies the text-based shop settings.
// PUT /api/settings
func (h *SettingHandler) Update(c *fiber.Ctx) error {
	var req domain.UpdateShopSettingsRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Printf("Setting Update BodyParser Error: %v, Body: %s\n", err, string(c.Body()))
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	settings, err := h.settingService.Update(c.Context(), req)
	if err != nil {
		return apiresponse.BadRequest(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, settings)
}

// UploadLogo handles multipart form data to upload a new shop logo.
// POST /api/settings/logo
func (h *SettingHandler) UploadLogo(c *fiber.Ctx) error {
	file, err := c.FormFile("logo")
	if err != nil {
		return apiresponse.BadRequest(c, "No logo file provided")
	}

	path, err := h.settingService.UploadLogo(c.Context(), file)
	if err != nil {
		return apiresponse.InternalError(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{
		"message":   "Logo uploaded successfully",
		"logo_path": path,
	})
}
