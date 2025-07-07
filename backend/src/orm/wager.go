package models

import (
	"powerscalers/backend/src/services"
	"time"

	"gorm.io/gorm"
)

type Winner struct {
	UserID int32
	Cash   float32
	Amount float32
}

func TerminateWager(
	DB *gorm.DB,
	WagerID int32) (string, int32) {
	DB.AutoMigrate(Wager{})

	var decision string
	DB.Model(&Wager{}).Select("Decision").Where("WagerID = ?", WagerID).Scan(&decision)

	var losing_total float32
	var winning_total float32
	DB.Model(&Gamble{}).Select("SUM(Amount)").Where("WagerID = ? AND Position != ?", WagerID, decision).Scan(&losing_total)
	DB.Model(&Gamble{}).Select("SUM(Amount)").Where("WagerID = ? AND Position = ?", WagerID, decision).Scan(&winning_total)
	var winners []Winner
	DB.Model(&Gamble{}).Select("User.UserID, User.Cash, Cash.Amount").Joins("JOIN User on User.UserID = Gamble.UserID").Where("WagerID = ? AND Position = ?", WagerID, decision).Scan(&winners)

	tx := DB.Begin()
	for i := range winners {
		reward := (winners[i].Amount / winning_total) * losing_total
		newCash := winners[i].Cash + reward

		if err := tx.Model(&User{}).Where("UserID = ?", winners[i].UserID).Update("Cash", newCash).Error; err != nil {
			tx.Rollback()
			return "Internal Error", 500
		}
	}
	tx.Commit()

	return "Success", 200
}

func MakeWager(
	DB *gorm.DB,
	OwnerID int32,
	CommunityID int32,
	Title string,
	Description string,
	ExpirationDate time.Time) (string, int32) {

	decision, explanation := services.Decide(Description)
	wager := Wager{
		CommunityID:    CommunityID,
		OwnerID:        OwnerID,
		Title:          Title,
		Decision:       decision,
		Explanation:    explanation,
		Description:    Description,
		ExpirationDate: ExpirationDate,
	}
	result := DB.Create(&wager)
	if result.Error != nil {
		return "Internal Error", 500
	}
	return "Success", 200
}

func MakeGamble(
	DB *gorm.DB,
	UserID int32,
	WagerID int32,
	Amount float32,
	Position string) (string, int32) {

	var exp_date time.Time
	DB.Model(&Wager{}).Select("ExpirationDate").Where("WagerID = ?", WagerID).Scan(&exp_date)

	if time.Now().Unix() > exp_date.Unix() {
		return "Wager Expiration Date Has Passed", 501
	}

	var cash float32
	DB.Model(&User{}).Select("Cash").Where("UserID = ?", UserID).Scan(&cash)
	if cash < Amount {
		return "Insufficient Amount of Money", 501
	}

	tx := DB.Begin()
	err := tx.Model(&User{}).Where("UserID = ?", UserID).Update("Cash", cash-Amount).Error
	if err != nil {
		tx.Rollback()
		return "Internal Error", 500
	}
	gamble := Gamble{
		UserID:   UserID,
		WagerID:  WagerID,
		Amount:   Amount,
		Position: Position,
	}
	result := tx.Create(&gamble)

	if result.Error != nil {
		tx.Rollback()
		return "Internal Error", 500
	}
	tx.Commit()
	return "Success", 200
}
