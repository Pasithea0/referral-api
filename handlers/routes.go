package handlers

import (
	"embed"
	"html"
	"referral-api/config"
	"referral-api/database"
	"referral-api/middleware"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

//go:embed static/*
var staticFiles embed.FS

func SetupRoutes(app *fiber.App) {
	cfg := config.Get()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.AllowedOrigins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization,X-Admin-Password",
		AllowCredentials: true,
	}))

	api := app.Group("/api")

	api.Get("/health", HealthHandler)
	api.Post("/webhook/signup", WebhookHandler)
	api.Post("/referral-codes", CreateReferralCodeHandler)
	api.Get("/referral-codes/:code", GetReferralCodeHandler)
	api.Get("/campaigns/:slug", GetCampaignHandler)

	admin := api.Group("/admin", middleware.AdminAuthMiddleware())
	admin.Get("/campaigns", ListCampaignsHandler)
	admin.Post("/campaigns", CreateCampaignHandler)
	admin.Get("/stats", AdminStatsHandler)

	app.Get("/", serveFile("static/dashboard.html", "text/html; charset=utf-8"))
	app.Get("/admin", serveFile("static/admin.html", "text/html; charset=utf-8"))
	app.Get("/:slug", serveCampaignPage)
}

func serveFile(path, contentType string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data, err := staticFiles.ReadFile(path)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Not found")
		}
		c.Set("Content-Type", contentType)
		return c.Send(data)
	}
}

func serveCampaignPage(c *fiber.Ctx) error {
	slug := c.Params("slug")

	if slug == "" || strings.Contains(slug, "/") || strings.HasPrefix(slug, ".") {
		return c.Status(fiber.StatusNotFound).SendString("Not found")
	}

	db, err := database.GetDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal error")
	}

	var campaign database.Campaign
	if err := db.Where("slug = ?", slug).First(&campaign).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Campaign not found")
	}

	data, err := staticFiles.ReadFile("static/campaign.html")
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Not found")
	}

	htmlContent := string(data)
	htmlContent = strings.ReplaceAll(htmlContent, "{{CAMPAIGN_SLUG}}", html.EscapeString(slug))
	htmlContent = strings.ReplaceAll(htmlContent, "{{CAMPAIGN_NAME}}", html.EscapeString(campaign.Name))
	htmlContent = strings.ReplaceAll(htmlContent, "{{CAMPAIGN_DESCRIPTION}}", html.EscapeString(campaign.Description))

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(htmlContent)
}
