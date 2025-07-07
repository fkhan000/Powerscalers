package models

import (
	"gorm.io/gorm"
)

func IncreaseCommunityLimit(
	DB *gorm.DB,
	UserID int32,
	CommunityID int32,
	Increment int32) (string, int32) {

	var modified_limit int32
	DB.Model(&Community{}).Select("ModifiedDailyLimit").Where("CommunityID = ?", CommunityID).Scan(&modified_limit)
	result := DB.Model(&Community{}).Where("CommunityID = ?", CommunityID).Update("ModifiedDailyLimit", modified_limit+Increment)
	if result.Error != nil {
		return "Internal Error", 500
	}
	return "Success", 200
}
func LikeWager(
	DB *gorm.DB,
	UserID int32,
	WagerID int32,
	Value int32) (string, int32) {

	wagerLike := WagerLike{
		UserID:  UserID,
		WagerID: WagerID,
		Value:   Value,
	}
	tx := DB.Begin()
	var net_likes int32
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
	UserID int32,
	CommentID int32,
	Value int32) (string, int32) {

	CommentLike := CommentLike{
		UserID:    UserID,
		CommentID: CommentID,
		Value:     Value,
	}
	tx := DB.Begin()
	var net_likes int32
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
