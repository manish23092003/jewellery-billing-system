package handler

import (
	"strconv"

	"jewellery-billing/internal/apiresponse"
	"jewellery-billing/internal/domain"
	"jewellery-billing/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ExpenseHandler exposes HTTP endpoints for expense management.
type ExpenseHandler struct {
	expenseService *service.ExpenseService
}

func NewExpenseHandler(expenseService *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{expenseService: expenseService}
}

// Create creates a new expense.
// POST /api/expenses
func (h *ExpenseHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	createdBy, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return apiresponse.Unauthorized(c, "Unable to identify user")
	}

	expense, err := h.expenseService.Create(c.Context(), req, createdBy)
	if err != nil {
		return apiresponse.BadRequest(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusCreated, expense)
}

// GetAll returns a paginated list of expenses with optional filters.
// GET /api/expenses?category=salary&date_from=2026-06-01&page=1
func (h *ExpenseHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))

	filter := domain.ExpenseFilter{
		Category: c.Query("category", ""),
		DateFrom: c.Query("date_from", ""),
		DateTo:   c.Query("date_to", ""),
		Page:     page,
		PerPage:  perPage,
	}

	expenses, total, err := h.expenseService.GetAll(c.Context(), filter)
	if err != nil {
		return apiresponse.InternalError(c, "Failed to fetch expenses")
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

// Delete removes an expense entry.
// DELETE /api/expenses/:id
func (h *ExpenseHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid expense ID")
	}

	if err := h.expenseService.Delete(c.Context(), id); err != nil {
		return apiresponse.NotFound(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{"message": "Expense deleted"})
}
