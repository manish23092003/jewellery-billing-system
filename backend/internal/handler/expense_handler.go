package handler

import (
	"strconv"

	"jewellery-billing/internal/apiresponse"
	"jewellery-billing/internal/domain"
	"jewellery-billing/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ExpenseHandler exposes HTTP endpoints for expense operations.
type ExpenseHandler struct {
	expenseService *service.ExpenseService
}

func NewExpenseHandler(expenseService *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{expenseService: expenseService}
}

// Create godoc
// POST /api/expenses
func (h *ExpenseHandler) Create(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)
	userID, _ := c.Locals("userID").(uuid.UUID)

	var req domain.CreateExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	if req.Category == "" || req.Amount <= 0 || req.Description == "" {
		return apiresponse.BadRequest(c, "Category, positive amount, and description are required")
	}

	expense, err := h.expenseService.Create(c.Context(), orgID, req, userID)
	if err != nil {
		return apiresponse.InternalError(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusCreated, expense)
}

// GetAll godoc
// GET /api/expenses
func (h *ExpenseHandler) GetAll(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))

	filter := domain.ExpenseFilter{
		Category: c.Query("category"),
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Page:     page,
		PerPage:  perPage,
	}

	expenses, total, err := h.expenseService.GetAll(c.Context(), orgID, filter)
	if err != nil {
		return apiresponse.InternalError(c, err.Error())
	}

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	return apiresponse.SuccessWithMeta(c, fiber.StatusOK, expenses, &apiresponse.Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// Delete godoc
// DELETE /api/expenses/:id
func (h *ExpenseHandler) Delete(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid expense ID")
	}

	if err := h.expenseService.Delete(c.Context(), orgID, id); err != nil {
		return apiresponse.NotFound(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{"message": "Expense deleted"})
}
