package models

import (
	"fmt"
	"powerscalers/backend/src/services"
	"time"

	"gorm.io/gorm"
)

type Winner struct {
	UserID int
	Cash   float32
	Amount float32
}

func RewardWager(
	DB *gorm.DB,
	WagerID int) (string, int) {
	DB.AutoMigrate(Wager{})

	var wagerDetails map[string]interface{}
	tx := DB.Model(&Wager{}).Select("Decision", "Left", "Right", "LeftAmount", "RightAmount").Where("WagerID = ?", WagerID).Scan(&wagerDetails)

	decision := wagerDetails["Decision"].(string)
	if err := tx.Model(&Wager{}).
		Select("decision").
		Where("wager_id = ?", WagerID).
		Row().
		Scan(&decision); err != nil {
		tx.Rollback()
		return "Internal Error", 500
	}
	if decision == "" {
		tx.Rollback()
		return "No decision set for wager", 400
	}

	var losingTotal float32
	var winningTotal float32

	if decision == wagerDetails["Left"].(string) {
		winningTotal = wagerDetails["LeftAmount"].(float32)
		losingTotal = wagerDetails["RightAmount"].(float32)
	} else {
		winningTotal = wagerDetails["RightAmount"].(float32)
		losingTotal = wagerDetails["LeftAmount"].(float32)
	}

	var winners []Winner
	DB.Model(&Gamble{}).Select("User.UserID, User.Cash, Cash.Amount").Joins("JOIN User on User.UserID = Gamble.UserID").Where("WagerID = ? AND Position = ?", WagerID, decision).Scan(&winners)

	tx = DB.Begin()
	for i := range winners {
		reward := (winners[i].Amount / winningTotal) * losingTotal

		if err := tx.Model(&User{}).Where("UserID = ?", winners[i].UserID).Update("Cash", gorm.Expr("cash + ?", reward)).Error; err != nil {
			tx.Rollback()
			return "Internal Error", 500
		}
	}
	tx.Commit()

	return "Success", 200
}

func MakeWager(
	DB *gorm.DB,
	OwnerID int,
	CommunityID int,
	Title string,
	Description string,
	Left string,
	Right string,
	NumVoR int,
	ExpirationDate time.Time) (string, int) {

	wager := Wager{
		CommunityID:    CommunityID,
		OwnerID:        OwnerID,
		Title:          Title,
		Description:    Description,
		Left:           Left,
		Right:          Right,
		NumVoR:         NumVoR,
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
	UserID int,
	WagerID int,
	Amount float32,
	Position string) (string, int) {

	var wagerInfo map[string]interface{}
	DB.Model(&Wager{}).Select("ExpirationDate", "Left", "Right", "LeftAmount", "RightAmount").Where("WagerID = ?", WagerID).Scan(&wagerInfo)

	expDate := wagerInfo["ExpirationDate"].(time.Time)
	left := wagerInfo["Left"].(string)
	right := wagerInfo["Right"].(string)

	if time.Now().Unix() > expDate.Unix() {
		return "Wager Expiration Date Has Passed", 501
	}
	if Position != left && Position != right {
		return "Wager Position Is Invalid", 502
	}

	var cash float32
	DB.Model(&User{}).Select("Cash").Where("UserID = ?", UserID).Scan(&cash)
	if cash < Amount {
		return "Insufficient Amount of Money", 503
	}

	tx := DB.Begin()
	err := tx.Model(&User{}).Where("UserID = ?", UserID).Update("Cash", cash-Amount).Error
	if err != nil {
		tx.Rollback()
		return "Internal Database Error", 500
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
		return "Internal Database Error", 500
	}
	if Position == left {
		leftAmount := wagerInfo["LeftAmount"].(float32)
		DB.Model(&Wager{}).Where("WagerID = ?", WagerID).Update("LeftAmount", leftAmount+Amount)
	} else {
		rightAmount := wagerInfo["RightAmount"].(float32)
		DB.Model(&Wager{}).Where("WagerID = ?", WagerID).Update("RightAmount", rightAmount+Amount)
	}
	tx.Commit()
	return "Success", 200
}

func TerminateWager(
	DB *gorm.DB,
	WagerID int) (string, int) {
	var wagerInfo map[string]interface{}
	DB.Model(&Wager{}).Select("Title", "Description", "Left", "Right", "NumVoR").Where("WagerID = ?", WagerID).Scan(&wagerInfo)
	comments, err := FetchComments(
		DB,
		WagerID,
		"Likes",
		"Descending",
		0,
		wagerInfo["NumVoR"].(int),
		1,
	)
	if err != nil {
		return fmt.Sprintf("Error retrieving voices of reason: %s", err), 504
	}
	voicesOfReason := []services.VoiceOfReason{}
	for _, comment := range comments {
		VoR := services.VoiceOfReason{
			Reason:   comment.Description,
			UserName: comment.UserName,
			Position: comment.Position,
			NumLikes: comment.NetLikes,
		}
		voicesOfReason = append(voicesOfReason, VoR)
	}
	wagerDecisionDetails := services.WagerDecisionDetails{
		Title:          wagerInfo["Title"].(string),
		Description:    wagerInfo["Description"].(string),
		Left:           wagerInfo["Left"].(string),
		Right:          wagerInfo["Right"].(string),
		VoicesOfReason: voicesOfReason,
	}
	decision, explanation, err := services.Decide(wagerDecisionDetails)
	if err != nil {
		return fmt.Sprintf("Error determining verdict for wager: %s", err), 505
	}

	DB.Model(&Wager{}).
		Where("WagerID = ?", WagerID).
		Updates(map[string]interface{}{
			"decision":    decision,
			"explanation": explanation,
		})
	return "Success", 200
}

func LoadWagers(
	DB *gorm.DB,
	CommunityID int,
	UserID int,
	Offset int,
	Limit int,
	SortCategory string,
	Ascending bool) ([]Wager, error) {

	sortParam2Field := map[string]string{"Likes": "w.NetLikes", "Time": "CreatedAt", "Amount": "Amount"}

	var sortCateg string
	var ok bool
	if sortCateg, ok = sortParam2Field[SortCategory]; !ok {
		return nil, fmt.Errorf("Invalid Sort Category Provided")
	}
	var sort string
	if Ascending {
		sort = "ASC"
	} else {
		sort = "DESC"
	}
	if UserID == -1 && CommunityID == -1 {
		return nil, fmt.Errorf("Cannot have both user and community IDs as -1")
	}

	var whereClause string
	if UserID == -1 {
		whereClause = fmt.Sprintf("w.community_id = %d", CommunityID)
	} else if CommunityID == -1 {
		whereClause = fmt.Sprintf("w.owner_id = %d", UserID)
	} else {
		return nil, fmt.Errorf("Cannot have both user and community IDs be greater than 0")
	}
	orderClause := sortCateg + " " + sort
	var results []Wager
	DB.Model(&Wager{}).
		Select(
			`w.*,
			u.user_name,
			u.profile_pic
		`).
		Joins("JOIN users AS u ON u.user_id == w.owner_id").
		Where(whereClause).
		Order(orderClause).
		Offset(Offset).
		Limit(Limit).
		Scan(&results)
	return results, nil
}
