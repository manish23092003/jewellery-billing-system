package handler

import (
	"jewellery-billing/internal/apiresponse"
	"jewellery-billing/internal/domain"
	"jewellery-billing/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UserHandler exposes HTTP endpoints for user management (admin only).
type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Create godoc
// POST /api/users
func (h *UserHandler) Create(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)

	var req domain.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return apiresponse.BadRequest(c, "Name, email, and password are required")
	}
	if !req.Role.IsValid() {
		req.Role = domain.RoleStaff
	}

	result, err := h.userService.Create(c.Context(), orgID, req)
	if err != nil {
		return apiresponse.BadRequest(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusCreated, result)
}

// GetAll godoc
// GET /api/users
func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	orgID, _ := c.Locals("organizationID").(uuid.UUID)

	users, err := h.userService.GetAll(c.Context(), orgID)
	if err != nil {
		return apiresponse.InternalError(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, users)
}

// GetByID godoc
// GET /api/users/:id
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid user ID")
	}

	user, err := h.userService.GetByID(c.Context(), id)
	if err != nil {
		return apiresponse.NotFound(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, user)
}

// Update godoc
// PUT /api/users/:id
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid user ID")
	}

	var req domain.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	result, err := h.userService.Update(c.Context(), id, req)
	if err != nil {
		return apiresponse.BadRequest(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, result)
}

// Delete godoc
// DELETE /api/users/:id
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid user ID")
	}

	if err := h.userService.Delete(c.Context(), id); err != nil {
		return apiresponse.NotFound(c, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{"message": "User deleted"})
}
