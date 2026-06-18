package middleware

import (
	"strings"

	"jewellery-billing/internal/service"
	"jewellery-billing/internal/apiresponse"

	"github.com/gofiber/fiber/v2"
)

// AuthRequired is a Fiber middleware that validates the JWT Bearer token
// in the Authorization header and injects user claims into c.Locals().
//
// Downstream handlers can access:
//
//	c.Locals("userID")  → uuid.UUID
//	c.Locals("email")   → string
//	c.Locals("role")    → domain.UserRole
func AuthRequired(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return apiresponse.Unauthorized(c, "Authorization header is required")
		}

		// Expect format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return apiresponse.Unauthorized(c, "Invalid authorization format — use: Bearer <token>")
		}

		claims, err := authService.ValidateToken(parts[1])
		if err != nil {
			return apiresponse.Unauthorized(c, "Invalid or expired token")
		}

		// Store claims in context for downstream use.
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

