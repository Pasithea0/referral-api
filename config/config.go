package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	AdminPassword string
	AllowedOrigins []string
}

var cfg *Config

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		return nil, fmt.Errorf("ADMIN_PASSWORD environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	var allowedOrigins []string
	if origins != "" {
		allowedOrigins = strings.Split(origins, ",")
	} else {
		allowedOrigins = []string{"http://localhost:3000", "https://theintrodb.org"}
	}

	cfg = &Config{
		Port:           port,
		DatabaseURL:    databaseURL,
		AdminPassword:  adminPassword,
		AllowedOrigins: allowedOrigins,
	}

	return cfg, nil
}

func Get() Config {
	if cfg == nil {
		return Config{}
	}
	return *cfg
}
