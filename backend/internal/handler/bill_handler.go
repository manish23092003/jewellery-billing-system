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

// BillHandler exposes HTTP endpoints for billing operations.
type BillHandler struct {
	billService    *service.BillService
	settingService *service.SettingService
	pdfService     *service.PDFService
}

func NewBillHandler(billService *service.BillService, settingService *service.SettingService, pdfService *service.PDFService) *BillHandler {
	return &BillHandler{
		billService:    billService,
		settingService: settingService,
		pdfService:     pdfService,
	}
}

// Create creates a new bill with multiple items.
// POST /api/bills
func (h *BillHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateBillRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	// Get the authenticated user's ID from context
	createdBy, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return apiresponse.Unauthorized(c, "Unable to identify user")
	}

	settings, _ := h.settingService.Get(c.Context())
	prefix := "INV"
	if settings != nil && settings.InvoicePrefix != "" {
		prefix = settings.InvoicePrefix
	}

	bill, err := h.billService.Create(c.Context(), req, createdBy, prefix)
	if err != nil {
		return apiresponse.BadRequest(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusCreated, bill)
}

// GetByID returns a single bill with all its items.
// GET /api/bills/:id
func (h *BillHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid bill ID")
	}

	bill, err := h.billService.GetByID(c.Context(), id)
	if err != nil {
		return apiresponse.NotFound(c, "Bill not found")
	}

	return apiresponse.Success(c, fiber.StatusOK, bill)
}

// GetAll returns a paginated list of bills with optional filters.
// GET /api/bills?search=INV&payment_method=cash&date_from=2026-01-01&date_to=2026-12-31&page=1&per_page=20
func (h *BillHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))

	filter := domain.BillFilter{
		Search:        c.Query("search", ""),
		PaymentMethod: c.Query("payment_method", ""),
		DateFrom:      c.Query("date_from", ""),
		DateTo:        c.Query("date_to", ""),
		Page:          page,
		PerPage:       perPage,
	}

	bills, total, err := h.billService.GetAll(c.Context(), filter)
	if err != nil {
		return apiresponse.InternalError(c, "Failed to fetch bills")
	}

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	return apiresponse.SuccessWithMeta(c, fiber.StatusOK, bills, &apiresponse.Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// Delete removes a bill and its items (admin only).
// DELETE /api/bills/:id
func (h *BillHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid bill ID")
	}

	if err := h.billService.Delete(c.Context(), id); err != nil {
		return apiresponse.NotFound(c, "Bill not found")
	}

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{"message": "Bill deleted"})
}

// DownloadPDF returns a PDF invoice.
// GET /api/bills/:id/pdf
func (h *BillHandler) DownloadPDF(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid bill ID")
	}
	bill, err := h.billService.GetByID(c.Context(), id)
	if err != nil {
		return apiresponse.NotFound(c, "Bill not found")
	}
	settings, _ := h.settingService.Get(c.Context())
	if settings == nil {
		settings = &domain.ShopSettings{ShopName: "Aura Jewels"}
	}

	pdfBytes, err := h.pdfService.GenerateInvoicePDF(bill, settings)
	if err != nil {
		return apiresponse.InternalError(c, "Failed to generate PDF")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=invoice_%s.pdf", bill.InvoiceNumber))
	return c.Send(pdfBytes)
}

