package handlers

import (
	"referral-api/database"

	"github.com/gofiber/fiber/v2"
)

func HealthHandler(c *fiber.Ctx) error {
	_, err := database.GetDB()
	dbOk := err == nil

	status := fiber.StatusOK
	if !dbOk {
		status = fiber.StatusServiceUnavailable
	}

	return c.Status(status).JSON(fiber.Map{
		"status":   "ok",
		"database": dbOk,
	})
}
