package handlers

import (
	"referral-api/database"

	"github.com/gofiber/fiber/v2"
)

type CreateCampaignRequest struct {
	Slug        string `json:"slug" validate:"required"`
	Name        string `json:"name" validate:"required"`
	BaseURL     string `json:"base_url" validate:"required"`
	Description string `json:"description"`
}

func CreateCampaignHandler(c *fiber.Ctx) error {
	var req CreateCampaignRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Slug == "" || req.Name == "" || req.BaseURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "slug, name, and base_url are required",
		})
	}

	db, err := database.GetDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database unavailable",
		})
	}

	campaign := database.Campaign{
		Slug:        req.Slug,
		Name:        req.Name,
		BaseURL:     req.BaseURL,
		Description: req.Description,
	}

	if err := db.Create(&campaign).Error; err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Campaign slug already exists",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":          campaign.ID,
		"slug":        campaign.Slug,
		"name":        campaign.Name,
		"base_url":    campaign.BaseURL,
		"description": campaign.Description,
	})
}

func ListCampaignsHandler(c *fiber.Ctx) error {
	db, err := database.GetDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database unavailable",
		})
	}

	var campaigns []database.Campaign
	if err := db.Order("created_at DESC").Find(&campaigns).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch campaigns",
		})
	}

	type campaignWithStats struct {
		ID          string `json:"id"`
		Slug        string `json:"slug"`
		Name        string `json:"name"`
		BaseURL     string `json:"base_url"`
		Description string `json:"description"`
		CodesCount  int64  `json:"codes_count"`
		Redemptions int64  `json:"redemptions"`
		CreatedAt   string `json:"created_at"`
	}

	result := make([]campaignWithStats, 0, len(campaigns))
	for _, c := range campaigns {
		var codesCount int64
		db.Model(&database.ReferralCode{}).Where("campaign_id = ?", c.ID).Count(&codesCount)

		var redemptions int64
		db.Raw(`
			SELECT COUNT(*) FROM referral_redemptions rr
			JOIN referral_codes rc ON rr.referral_code_id = rc.id
			WHERE rc.campaign_id = ?
		`, c.ID).Scan(&redemptions)

		result = append(result, campaignWithStats{
			ID:          c.ID,
			Slug:        c.Slug,
			Name:        c.Name,
			BaseURL:     c.BaseURL,
			Description: c.Description,
			CodesCount:  codesCount,
			Redemptions: redemptions,
			CreatedAt:   c.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"campaigns": result,
	})
}

func GetCampaignHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")

	db, err := database.GetDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database unavailable",
		})
	}

	var campaign database.Campaign
	if err := db.Where("slug = ?", slug).First(&campaign).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Campaign not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":          campaign.ID,
		"slug":        campaign.Slug,
		"name":        campaign.Name,
		"base_url":    campaign.BaseURL,
		"description": campaign.Description,
	})
}
