package database

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func InitDB() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			err = fmt.Errorf("DATABASE_URL environment variable is required")
			return
		}

		db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
		if err != nil {
			err = fmt.Errorf("failed to connect to database: %w", err)
			return
		}

		if err = db.AutoMigrate(&Campaign{}, &ReferralCode{}, &ReferralRedemption{}); err != nil {
			return
		}

		log.Println("Database connection established and schema migrated")
	})

	return db, err
}

func GetDB() (*gorm.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized. Call InitDB() first")
	}
	return db, nil
}
