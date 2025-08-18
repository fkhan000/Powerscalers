package models

import (
	"fmt"

	"gorm.io/gorm"
)

func IncreaseCommunityLimit(
	DB *gorm.DB,
	UserID int,
	CommunityID int,
	Increment int) (string, int) {

	var modified_limit int
	DB.Model(&Community{}).Select("NumPostsAllowed").Where("CommunityID = ?", CommunityID).Scan(&modified_limit)
	result := DB.Model(&Community{}).Where("CommunityID = ?", CommunityID).Update("NumPostsAllowed", modified_limit+Increment)
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
		return fmt.Sprintf("Error updating net likes for wager %s", err), 500
	}
	result := tx.Create(&CommentLike)
	if result.Error != nil {
		tx.Rollback()
		return fmt.Sprintf("Error creating comment %s", result.Error), 500
	}
	tx.Commit()
	return "Success", 200
}

func LoadUserProfile(
	DB *gorm.DB,
	UserID int) (User, error) {

	var joinedCommunities []CommunityPreview

	DB.Table("community as c").
		Select(`
			c.community_id,
			c.Name,
			c.Description,
			c.Picture,
			c.NumFollowers
		`).
		Joins("JOIN member AS m ON m.community_id = c.community_id").
		Where("m.user_id = ?", UserID).
		Scan(&joinedCommunities)

	var user User
	DB.Table("users").
		Where("user_id = ?", UserID).
		Scan(&user)
	user.JoinedCommunities = joinedCommunities

	return user, nil
}
func LoadForYouPage(
	DB *gorm.DB,
	UserID int,
	Offset int,
	Limit int) ([]Wager, error) {

	var communities []int
	if err := DB.Model(&Member{}).
		Select("community_id").
		Where("user_id = ?", UserID).
		Scan(&communities).Error; err != nil {
		return nil, err
	}

	var wagers []Wager
	err := DB.Table("wagers AS w").
		Select(`
            w.*,
            u.user_name,
            u.profile_pic,
			c.name,
			c.picture
        `).
		Joins("JOIN users AS u ON u.user_id = w.owner_id").
		Joins("JOIN community AS c ON c.community_id = w.community_id").
		Where("w.community_id IN ? AND (EXTRACT(EPOCH FROM (NOW() - w.created_at)) / 86400) < 7", communities).
		Order("(w.net_likes / POWER(EXTRACT(EPOCH FROM (NOW() - w.created_at)) / 3600 + 2, 1.2)) + (RANDOM() * 0.1) DESC").
		Offset(Offset).
		Limit(Limit).
		Scan(&wagers).Error

	if err != nil {
		return nil, err
	}
	return wagers, nil
}
