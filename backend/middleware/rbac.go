package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := fmt.Sprintf("%v", c.Locals("role"))
		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}
		return c.Status(403).JSON(fiber.Map{"error": "Insufficient permissions"})
	}
}

func GetUserID(c *fiber.Ctx) string {
	if v, ok := c.Locals("user_id").(uuid.UUID); ok {
		return v.String()
	}
	return ""
}

func GetTenantID(c *fiber.Ctx) string {
	if v, ok := c.Locals("tenant_id").(uuid.UUID); ok {
		return v.String()
	}
	return ""
}

func GetBranchID(c *fiber.Ctx) string {
	if v, ok := c.Locals("branch_id").(uuid.UUID); ok {
		return v.String()
	}
	return ""
}
