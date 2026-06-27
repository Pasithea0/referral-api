package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"referral-api/config"
	"referral-api/database"
	"referral-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if _, err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
	})

	handlers.SetupRoutes(app)

	go func() {
		log.Printf("Referral API starting on port %s", cfg.Port)
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
