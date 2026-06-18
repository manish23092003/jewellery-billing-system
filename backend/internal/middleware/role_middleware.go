package middleware

import (
	"jewellery-billing/internal/domain"
	"jewellery-billing/internal/apiresponse"

	"github.com/gofiber/fiber/v2"
)

// AdminOnly is a Fiber middleware that restricts access to admin users.
// Must be placed AFTER AuthRequired in the middleware chain.
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(domain.UserRole)
		if !ok {
			return apiresponse.Forbidden(c, "Unable to determine user role")
		}

		if role != domain.RoleAdmin {
			return apiresponse.Forbidden(c, "This action requires admin privileges")
		}

		return c.Next()
	}
}

// RoleRequired is a generic middleware that accepts a list of allowed roles.
// Useful when a route should be accessible to multiple (but not all) roles.
func RoleRequired(roles ...domain.UserRole) fiber.Handler {
	allowed := make(map[domain.UserRole]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}

	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(domain.UserRole)
		if !ok {
			return apiresponse.Forbidden(c, "Unable to determine user role")
		}

		if !allowed[role] {
			return apiresponse.Forbidden(c, "Insufficient permissions")
		}

		return c.Next()
	}
}

