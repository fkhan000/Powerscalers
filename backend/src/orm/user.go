package models

import (
	"powerscalers/backend/src/services"
	"time"

	"gorm.io/gorm"
)

func (u *User) makeWager(
	DB *gorm.DB,
	CommunityID *int32,
	Title string,
	Description string,
	ExpirationDate time.Time) int32 {

	DB.AutoMigrate(Wager{})

	decision, explanation := services.Decide(Description)
	wager := Wager{
		CommunityID:    CommunityID,
		Title:          Title,
		Decision:       decision,
		Explanation:    explanation,
		Description:    Description,
		ExpirationDate: ExpirationDate,
	}
	result := DB.Create(&wager)
	if result != nil {
		return 500
	}
	return 200
}
