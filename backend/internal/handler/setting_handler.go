package handler

import (
	"fmt"
	"os"
	"path/filepath"

	"jewellery-billing/internal/apiresponse"
	"jewellery-billing/internal/domain"
	"jewellery-billing/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// SettingHandler exposes HTTP endpoints for shop settings.
type SettingHandler struct {
	settingService *service.SettingService
}

func NewSettingHandler(settingService *service.SettingService) *SettingHandler {
	return &SettingHandler{settingService: settingService}
}

// Get godoc
// GET /api/settings
func (h *SettingHandler) Get(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)

	settings, err := h.settingService.Get(c.Context(), orgID)
	if err != nil {
		return apiresponse.InternalError(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, settings)
}

// Update godoc
// PUT /api/settings
func (h *SettingHandler) Update(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)

	var req domain.UpdateShopSettingsRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	result, err := h.settingService.Update(c.Context(), orgID, req)
	if err != nil {
		return apiresponse.InternalError(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, result)
}

// UploadLogo godoc
// POST /api/settings/logo
func (h *SettingHandler) UploadLogo(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)

	file, err := c.FormFile("logo")
	if err != nil {
		return apiresponse.BadRequest(c, "Logo file is required")
	}

	// Save to uploads directory scoped by org
	uploadDir := filepath.Join("uploads", "logos")
	os.MkdirAll(uploadDir, 0755)
	filename := fmt.Sprintf("%s_%s", orgID.String()[:8], file.Filename)
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveFile(file, savePath); err != nil {
		return apiresponse.InternalError(c, "Failed to save logo")
	}

	logoPath := fmt.Sprintf("/uploads/logos/%s", filename)
	if err := h.settingService.UpdateLogo(c.Context(), orgID, logoPath); err != nil {
		return apiresponse.InternalError(c, "Failed to update logo path")
	}

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{
		"logo_path": logoPath,
	})
}
