package handlers

import (
	"referral-api/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SignupWebhookRequest struct {
	NewUserID    string `json:"new_user_id" validate:"required"`
	ReferralCode string `json:"referral_code" validate:"required"`
}

func WebhookHandler(c *fiber.Ctx) error {
	var req SignupWebhookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.NewUserID == "" || req.ReferralCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "new_user_id and referral_code are required",
		})
	}

	db, err := database.GetDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database unavailable",
		})
	}

	var code database.ReferralCode
	if err := db.Where("code = ?", req.ReferralCode).First(&code).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Referral code not found",
		})
	}

	var existing database.ReferralRedemption
	if err := db.Where("new_user_id = ?", req.NewUserID).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User already recorded as referred",
		})
	}

	redemption := database.ReferralRedemption{
		ID:             uuid.New().String(),
		NewUserID:      req.NewUserID,
		ReferralCodeID: code.ID,
	}

	if err := db.Create(&redemption).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to record referral",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"ok": true,
	})
}
