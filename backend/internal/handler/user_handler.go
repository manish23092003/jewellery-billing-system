package handler

import (
	"jewellery-billing/internal/domain"
	"jewellery-billing/internal/service"
	"jewellery-billing/internal/apiresponse"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UserHandler exposes HTTP endpoints for user management (admin-only).
type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Create godoc
// POST /api/users (admin)
// Creates a new user account with the given name, email, password, and role.
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return apiresponse.BadRequest(c, "Name, email, and password are required")
	}

	if !req.Role.IsValid() {
		return apiresponse.BadRequest(c, "Invalid role — must be 'admin' or 'staff'")
	}

	result, err := h.userService.Create(c.Context(), req)
	if err != nil {
		return apiresponse.Error(c, fiber.StatusConflict, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusCreated, result)
}

// GetAll godoc
// GET /api/users (admin)
// Returns all users in the system.
func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	users, err := h.userService.GetAll(c.Context())
	if err != nil {
		return apiresponse.InternalError(c, "Failed to fetch users")
	}

	return apiresponse.Success(c, fiber.StatusOK, users)
}

// GetByID godoc
// GET /api/users/:id (admin)
// Returns a single user by UUID.
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid user ID format")
	}

	user, err := h.userService.GetByID(c.Context(), id)
	if err != nil {
		return apiresponse.NotFound(c, "User not found")
	}

	return apiresponse.Success(c, fiber.StatusOK, user)
}

// Update godoc
// PUT /api/users/:id (admin)
// Partially updates a user — only non-empty fields are applied.
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid user ID format")
	}

	var req domain.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.BadRequest(c, "Invalid request body")
	}

	result, err := h.userService.Update(c.Context(), id, req)
	if err != nil {
		return apiresponse.Error(c, fiber.StatusConflict, err.Error())
	}

	return apiresponse.Success(c, fiber.StatusOK, result)
}

// Delete godoc
// DELETE /api/users/:id (admin)
// Permanently removes a user account.
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apiresponse.BadRequest(c, "Invalid user ID format")
	}

	if err := h.userService.Delete(c.Context(), id); err != nil {
		return apiresponse.NotFound(c, "User not found")
	}

	return apiresponse.Success(c, fiber.StatusOK, fiber.Map{
		"message": "User deleted successfully",
	})
}

