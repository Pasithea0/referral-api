package handlers

import (
	"crypto/rand"
	"math/big"
	"referral-api/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const codeCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func generateReferralCode() (string, error) {
	code := make([]byte, 10)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(codeCharset))))
		if err != nil {
			return "", err
		}
		code[i] = codeCharset[n.Int64()]
	}
	return string(code), nil
}

type CreateReferralCodeRequest struct {
	CampaignSlug string `json:"campaign_slug" validate:"required"`
	Email        string `json:"email" validate:"required"`
	Discord      string `json:"discord" validate:"required"`
}

func CreateReferralCodeHandler(c *fiber.Ctx) error {
	var req CreateReferralCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.CampaignSlug == "" || req.Email == "" || req.Discord == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "campaign_slug, email, and discord are required",
		})
	}

	db, err := database.GetDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database unavailable",
		})
	}

	var campaign database.Campaign
	if err := db.Where("slug = ?", req.CampaignSlug).First(&campaign).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Campaign not found",
		})
	}

	var existing database.ReferralCode
	result := db.Where("campaign_id = ? AND (owner_email = ? OR owner_discord = ?)",
		campaign.ID, req.Email, req.Discord).First(&existing)

	if result.Error == nil {
		var count int64
		db.Model(&database.ReferralRedemption{}).
			Where("referral_code_id = ?", existing.ID).
			Count(&count)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"code":        existing.Code,
			"share_link":  campaign.BaseURL + "?refer=" + existing.Code,
			"redemptions": count,
			"existing":    true,
		})
	}

	code, err := generateReferralCode()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate referral code",
		})
	}

	refCode := database.ReferralCode{
		ID:           uuid.New().String(),
		Code:         code,
		CampaignID:   campaign.ID,
		OwnerEmail:   req.Email,
		OwnerDiscord: req.Discord,
	}

	if err := db.Create(&refCode).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create referral code",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":        code,
		"share_link":  campaign.BaseURL + "?refer=" + code,
		"redemptions": 0,
		"existing":    false,
	})
}

func GetReferralCodeHandler(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Code is required",
		})
	}

	db, err := database.GetDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database unavailable",
		})
	}

	var refCode database.ReferralCode
	if err := db.Preload("Campaign").Where("code = ?", code).First(&refCode).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Referral code not found",
		})
	}

	var count int64
	db.Model(&database.ReferralRedemption{}).
		Where("referral_code_id = ?", refCode.ID).
		Count(&count)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":        refCode.Code,
		"campaign":    refCode.Campaign.Name,
		"campaign_slug": refCode.Campaign.Slug,
		"owner_email": refCode.OwnerEmail,
		"owner_discord": refCode.OwnerDiscord,
		"redemptions": count,
		"share_link":  refCode.Campaign.BaseURL + "?refer=" + refCode.Code,
	})
}
