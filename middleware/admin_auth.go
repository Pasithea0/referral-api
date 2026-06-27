package middleware

import (
	"referral-api/config"

	"github.com/gofiber/fiber/v2"
)

func AdminAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		password := c.Get("X-Admin-Password")
		if password == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "X-Admin-Password header required",
			})
		}

		cfg := config.Get()
		if password != cfg.AdminPassword {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Invalid admin password",
			})
		}

		return c.Next()
	}
}
