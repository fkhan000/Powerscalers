package models

import (
	"gorm.io/gorm"
)

func IncreaseCommunityLimit(
	DB *gorm.DB,
	UserID int,
	CommunityID int,
	Increment int) (string, int) {

	var modified_limit int
	DB.Model(&Community{}).Select("ModifiedDailyLimit").Where("CommunityID = ?", CommunityID).Scan(&modified_limit)
	result := DB.Model(&Community{}).Where("CommunityID = ?", CommunityID).Update("ModifiedDailyLimit", modified_limit+Increment)
	if result.Error != nil {
		return "Internal Error", 500
	}
	return "Success", 200
}
func LikeWager(
	DB *gorm.DB,
	UserID int,
	WagerID int,
	Value int) (string, int) {

	wagerLike := WagerLike{
		UserID:  UserID,
		WagerID: WagerID,
		Value:   Value,
	}
	tx := DB.Begin()
	var net_likes int
	tx.Model(&Wager{}).Select("NetLikes").Where("WagerID = ?", WagerID).Scan(&net_likes)
	err := tx.Model(&Wager{}).Where("WagerID = ?", WagerID).Update("NetLikes", net_likes+Value).Error

	if err != nil {
		tx.Rollback()
		return "Internal Error", 500
	}
	result := tx.Create(&wagerLike)
	if result.Error != nil {
		tx.Rollback()
		return "Internal Error", 500
	}
	tx.Commit()
	return "Success", 200
}

func LikeComment(
	DB *gorm.DB,
	UserID int,
	CommentID int,
	Value int) (string, int) {

	CommentLike := CommentLike{
		UserID:    UserID,
		CommentID: CommentID,
		Value:     Value,
	}
	tx := DB.Begin()
	var net_likes int
	tx.Model(&Comment{}).Select("NetLikes").Where("CommentID = ?", CommentID).Scan(&net_likes)
	err := tx.Model(&Comment{}).Where("CommentID = ?", CommentID).Update("NetLikes", net_likes+Value).Error

	if err != nil {
		tx.Rollback()
		return "Internal Error", 500
	}
	result := tx.Create(&CommentLike)
	if result.Error != nil {
		tx.Rollback()
		return "Internal Error", 500
	}
	tx.Commit()
	return "Success", 200
}
