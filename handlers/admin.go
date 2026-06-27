package handlers

import (
	"referral-api/database"

	"github.com/gofiber/fiber/v2"
)

type recentRedemption struct {
	NewUserID   string `json:"new_user_id"`
	Code        string `json:"code"`
	OwnerEmail  string `json:"owner_email"`
	Campaign    string `json:"campaign"`
	CreatedAt   string `json:"created_at"`
}

type topCode struct {
	Code        string `json:"code"`
	OwnerEmail  string `json:"owner_email"`
	OwnerDiscord string `json:"owner_discord"`
	Campaign    string `json:"campaign"`
	Redemptions int64  `json:"redemptions"`
}

type adminStatsResponse struct {
	TotalCampaigns   int64             `json:"total_campaigns"`
	TotalCodes       int64             `json:"total_codes"`
	TotalRedemptions int64             `json:"total_redemptions"`
	RecentRedemptions []recentRedemption `json:"recent_redemptions"`
	TopCodes         []topCode         `json:"top_codes"`
}

func AdminStatsHandler(c *fiber.Ctx) error {
	db, err := database.GetDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database unavailable",
		})
	}

	campaignID := c.Query("campaign_id")

	var totalCampaigns int64
	db.Model(&database.Campaign{}).Count(&totalCampaigns)

	var totalCodes int64
	codesQuery := db.Model(&database.ReferralCode{})
	if campaignID != "" {
		codesQuery = codesQuery.Where("campaign_id = ?", campaignID)
	}
	codesQuery.Count(&totalCodes)

	var totalRedemptions int64
	if campaignID != "" {
		db.Raw(`
			SELECT COUNT(*) FROM referral_redemptions rr
			JOIN referral_codes rc ON rr.referral_code_id = rc.id
			WHERE rc.campaign_id = ?
		`, campaignID).Scan(&totalRedemptions)
	} else {
		db.Model(&database.ReferralRedemption{}).Count(&totalRedemptions)
	}

	var recent []struct {
		NewUserID    string
		Code         string
		OwnerEmail   string
		CampaignName string
		CreatedAt    string
	}

	var recentQuery string
	var recentArgs []interface{}
	if campaignID != "" {
		recentQuery = `
			SELECT rr.new_user_id, rc.code, rc.owner_email, c.name as campaign_name,
			       TO_CHAR(rr.created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at
			FROM referral_redemptions rr
			JOIN referral_codes rc ON rr.referral_code_id = rc.id
			JOIN campaigns c ON rc.campaign_id = c.id
			WHERE rc.campaign_id = ?
			ORDER BY rr.created_at DESC
			LIMIT 20
		`
		recentArgs = []interface{}{campaignID}
	} else {
		recentQuery = `
			SELECT rr.new_user_id, rc.code, rc.owner_email, c.name as campaign_name,
			       TO_CHAR(rr.created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at
			FROM referral_redemptions rr
			JOIN referral_codes rc ON rr.referral_code_id = rc.id
			JOIN campaigns c ON rc.campaign_id = c.id
			ORDER BY rr.created_at DESC
			LIMIT 20
		`
	}
	db.Raw(recentQuery, recentArgs...).Scan(&recent)

	recentRedemptions := make([]recentRedemption, 0, len(recent))
	for _, r := range recent {
		recentRedemptions = append(recentRedemptions, recentRedemption{
			NewUserID:  r.NewUserID,
			Code:       r.Code,
			OwnerEmail: r.OwnerEmail,
			Campaign:   r.CampaignName,
			CreatedAt:  r.CreatedAt,
		})
	}

	var top []struct {
		Code         string
		OwnerEmail   string
		OwnerDiscord string
		CampaignName string
		Redemptions  int64
	}

	var topQuery string
	var topArgs []interface{}
	if campaignID != "" {
		topQuery = `
			SELECT rc.code, rc.owner_email, rc.owner_discord, c.name as campaign_name,
			       COUNT(rr.id) as redemptions
			FROM referral_codes rc
			JOIN campaigns c ON rc.campaign_id = c.id
			LEFT JOIN referral_redemptions rr ON rr.referral_code_id = rc.id
			WHERE rc.campaign_id = ?
			GROUP BY rc.id, rc.code, rc.owner_email, rc.owner_discord, c.name
			ORDER BY redemptions DESC
			LIMIT 10
		`
		topArgs = []interface{}{campaignID}
	} else {
		topQuery = `
			SELECT rc.code, rc.owner_email, rc.owner_discord, c.name as campaign_name,
			       COUNT(rr.id) as redemptions
			FROM referral_codes rc
			JOIN campaigns c ON rc.campaign_id = c.id
			LEFT JOIN referral_redemptions rr ON rr.referral_code_id = rc.id
			GROUP BY rc.id, rc.code, rc.owner_email, rc.owner_discord, c.name
			ORDER BY redemptions DESC
			LIMIT 10
		`
	}
	db.Raw(topQuery, topArgs...).Scan(&top)

	topCodes := make([]topCode, 0, len(top))
	for _, t := range top {
		topCodes = append(topCodes, topCode{
			Code:         t.Code,
			OwnerEmail:   t.OwnerEmail,
			OwnerDiscord: t.OwnerDiscord,
			Campaign:     t.CampaignName,
			Redemptions:  t.Redemptions,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total_campaigns":   totalCampaigns,
		"total_codes":       totalCodes,
		"total_redemptions": totalRedemptions,
		"recent_redemptions": recentRedemptions,
		"top_codes":         topCodes,
	})
}
