package handler

import (
	"strconv"

	"jewellery-billing/internal/apiresponse"
	"jewellery-billing/internal/domain"
	"jewellery-billing/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CustomerHandler struct {
	customerService *service.CustomerService
}

func NewCustomerHandler(customerService *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: customerService}
}

// Create godoc
// POST /api/customers
func (h *CustomerHandler) Create(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)

	var req domain.CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	customer, err := h.customerService.Create(c.Context(), orgID, req)
	if err != nil {
		return apiresponse.BadRequest(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusCreated, customer)
}

// GetAll godoc
// GET /api/customers
func (h *CustomerHandler) GetAll(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)

	search := c.Query("search", "")
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	customers, total, err := h.customerService.GetAll(c.Context(), orgID, search, limit, offset)
	if err != nil {
		return apiresponse.InternalError(c, err.Error())
	}

	return apiresponse.SuccessWithMeta(c, fiber.StatusOK, customers, &apiresponse.Meta{
		Total: total,
	})
}

// GetByID godoc
// GET /api/customers/:id
func (h *CustomerHandler) GetByID(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid customer ID")
	}

	customer, err := h.customerService.GetByID(c.Context(), orgID, id)
	if err != nil {
		return apiresponse.NotFound(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, customer)
}

// Update godoc
// PUT /api/customers/:id
func (h *CustomerHandler) Update(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid customer ID")
	}

	var req domain.UpdateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	customer, err := h.customerService.Update(c.Context(), orgID, id, req)
	if err != nil {
		return apiresponse.BadRequest(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, customer)
}

// Delete godoc
// DELETE /api/customers/:id
func (h *CustomerHandler) Delete(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid customer ID")
	}

	if err := h.customerService.Delete(c.Context(), orgID, id); err != nil {
		return apiresponse.NotFound(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{"message": "Customer deleted successfully"})
}
