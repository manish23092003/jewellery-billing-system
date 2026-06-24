package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"jewellery-billing/internal/apiresponse"
	"jewellery-billing/internal/domain"
	"jewellery-billing/internal/repository"
	"jewellery-billing/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BillHandler struct {
	billService    *service.BillService
	settingService *service.SettingService
	pdfService     *service.PDFService
	auditRepo      repository.VerificationLogRepository
}

func NewBillHandler(billService *service.BillService, settingService *service.SettingService, pdfService *service.PDFService, auditRepo repository.VerificationLogRepository) *BillHandler {
	return &BillHandler{
		billService:    billService,
		settingService: settingService,
		pdfService:     pdfService,
		auditRepo:      auditRepo,
	}
}

// Create godoc
// POST /api/bills
func (h *BillHandler) Create(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)
	userID, _ := c.Locals("userID").(uuid.UUID)

	var req domain.CreateBillRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	invoicePrefix := h.settingService.GetInvoicePrefix(c.Context(), orgID)

	bill, err := h.billService.Create(c.Context(), orgID, req, userID, invoicePrefix)
	if err != nil {
		return apiresponse.BadRequest(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusCreated, bill)
}

// GetAll godoc
// GET /api/bills
func (h *BillHandler) GetAll(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))

	filter := domain.BillFilter{
		Search:        c.Query("search"),
		PaymentMethod: c.Query("payment_method"),
		DateFrom:      c.Query("date_from"),
		DateTo:        c.Query("date_to"),
		Page:          page,
		PerPage:       perPage,
	}

	bills, total, err := h.billService.GetAll(c.Context(), orgID, filter)
	if err != nil {
		return apiresponse.InternalError(c, err.Error())
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

// GetByID godoc
// GET /api/bills/:id
func (h *BillHandler) GetByID(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid bill ID")
	}

	bill, err := h.billService.GetByID(c.Context(), orgID, id)
	if err != nil {
		return apiresponse.NotFound(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, bill)
}

// Delete godoc
// DELETE /api/bills/:id
func (h *BillHandler) Delete(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid bill ID")
	}

	if err := h.billService.Delete(c.Context(), orgID, id); err != nil {
		return apiresponse.NotFound(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{"message": "Bill deleted"})
}

// VerifyInvoice godoc
// GET /api/public/verify/:token
func (h *BillHandler) VerifyInvoice(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return apiresponse.BadRequest(c, "Verification token is required")
	}

	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	bill, response, err := h.billService.VerifyInvoice(c.Context(), token)
	
	// Prepare audit log
	auditLog := &domain.VerificationLog{
		Token:         token,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		IsValid:       err == nil && response.VerificationStatus != "TAMPERED",
	}

	if err != nil {
		auditLog.FailureReason = err.Error()
		if h.auditRepo != nil {
			_ = h.auditRepo.LogVerification(c.Context(), auditLog)
		}
		return apiresponse.NotFound(c, "Invoice not found or invalid token")
	}

	settings, err := h.settingService.Get(c.Context(), bill.OrganizationID)
	if err == nil {
		response.ShopName = settings.ShopName
	}

	if h.auditRepo != nil {
		_ = h.auditRepo.LogVerification(c.Context(), auditLog)
	}

	return apiresponse.Success(c, fiber.StatusOK, response)
}

// DownloadPDF godoc
// GET /api/bills/:id/pdf
func (h *BillHandler) DownloadPDF(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid bill ID")
	}

	bill, err := h.billService.GetByID(c.Context(), orgID, id)
	if err != nil {
		return apiresponse.NotFound(c, err.Error())
	}

	settings, err := h.settingService.Get(c.Context(), orgID)
	if err != nil {
		return apiresponse.InternalError(c, "Failed to fetch settings")
	}

	pdfBytes, err := h.pdfService.GenerateInvoicePDF(bill, settings)
	if err != nil {
		return apiresponse.InternalError(c, "Failed to generate PDF")
	}

	// Save PDF to disk
	pdfDir := filepath.Join("uploads", "invoices")
	os.MkdirAll(pdfDir, 0755)
	pdfPath := filepath.Join(pdfDir, fmt.Sprintf("%s.pdf", bill.InvoiceNumber))
	os.WriteFile(pdfPath, pdfBytes, 0644)

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.pdf\"", bill.InvoiceNumber))

	return c.Send(pdfBytes)
}

func (h *BillHandler) AddPayment(c *fiber.Ctx) error {
	orgID, ok := c.Locals("organizationID").(uuid.UUID)
	if !ok || orgID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	billIDStr := c.Params("id")
	billID, err := uuid.Parse(billIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid bill ID format"})
	}

	var req domain.AddPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err = h.billService.AddPayment(c.Context(), orgID, billID, req.Amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Payment added successfully",
	})
}
