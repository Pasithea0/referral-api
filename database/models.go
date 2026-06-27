package database

import "time"

type Campaign struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Slug        string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Name        string    `gorm:"type:varchar(255);not null"`
	BaseURL     string    `gorm:"type:varchar(500);not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

type ReferralCode struct {
	ID           string   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code         string   `gorm:"type:varchar(10);uniqueIndex;not null"`
	CampaignID   string   `gorm:"type:uuid;not null;uniqueIndex:idx_campaign_email;uniqueIndex:idx_campaign_discord"`
	Campaign     Campaign `gorm:"foreignKey:CampaignID"`
	OwnerEmail   string   `gorm:"type:varchar(255);not null;uniqueIndex:idx_campaign_email"`
	OwnerDiscord string   `gorm:"type:varchar(255);not null;uniqueIndex:idx_campaign_discord"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ReferralRedemption struct {
	ID             string       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	NewUserID      string       `gorm:"type:varchar(255);uniqueIndex;not null"`
	ReferralCodeID string       `gorm:"type:uuid;not null;index"`
	ReferralCode   ReferralCode `gorm:"foreignKey:ReferralCodeID"`
	CreatedAt      time.Time    `gorm:"not null"`
}

func (Campaign) TableName() string {
	return "campaigns"
}

func (ReferralCode) TableName() string {
	return "referral_codes"
}

func (ReferralRedemption) TableName() string {
	return "referral_redemptions"
}
