package middleware

import (
	"slices"

	"github.com/danielreinhard1129/fiber-clean-arch/pkg/exception"
	"github.com/gofiber/fiber/v2"
)

func RequireRoles(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return exception.UnauthorizedError{Message: "Role not found in token"}
		}

		if slices.Contains(roles, userRole) {
			return c.Next()
		}

		return exception.UnauthorizedError{Message: "You do not have permission to access this resource"}
	}
}
