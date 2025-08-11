package models

import (
	"gorm.io/gorm"
)

type CommunityUpdates struct {
	Description string
	Picture     string
}

func CreateCommunity(
	DB *gorm.DB,
	UserID int,
	Name string,
	Description string,
	Picture string) (string, int) {

	community := Community{
		Name:        Name,
		Description: Description,
		Picture:     Picture,
	}
	tx := DB.Begin()
	result := tx.Model(&Community{}).Create(&community)
	if result.Error != nil {
		tx.Rollback()
		return "Internal Error", 500
	}

	member := Member{
		UserID:      UserID,
		CommunityID: community.CommunityID,
	}
	result = tx.Model(&Member{}).Create(member)
	if result.Error != nil {
		tx.Rollback()
		return "Internal Error", 500
	}

	return "Success", 200
}

func UpdateCommunity(
	DB *gorm.DB,
	CommunityID int,
	Description string,
	Picture string) (string, int) {

	updates := CommunityUpdates{
		Description: Description,
		Picture:     Picture,
	}
	result := DB.Model(&Community{}).Where("CommunityID = ?", CommunityID).UpdateColumns(updates)
	if result.Error != nil {
		return "Internal Error", 500
	}
	return "Success", 200
}

func JoinCommunity(
	DB *gorm.DB,
	UserID int,
	CommunityID int) (string, int) {

	tx := DB.Begin()
	member := Member{
		UserID:      UserID,
		CommunityID: CommunityID,
	}
	result := tx.Create(&member)

	if result.Error != nil {
		tx.Rollback()
		return "Internal Error", 500
	}

	var num_followers int
	tx.Model(&Community{}).Select("NumFollowers").Where("CommunityID = ?", CommunityID).Scan(&num_followers)

	result = tx.Model(&Community{}).Where("CommunityID = ?", CommunityID).Update("NumFollowers", num_followers+1)
	if result.Error != nil {
		tx.Rollback()
		return "Internal Error", 500
	}
	tx.Commit()
	return "Success", 200
}

func LeaveCommunity(
	DB *gorm.DB,
	UserID int,
	CommunityID int) (string, int) {

	tx := DB.Begin()
	member := Member{
		UserID:      UserID,
		CommunityID: CommunityID,
	}
	result := tx.Delete(&member)

	if result.Error != nil {
		tx.Rollback()
		return "Internal Error", 500
	}

	var num_followers int
	tx.Model(&Community{}).Select("NumFollowers").Where("CommunityID = ?", CommunityID).Scan(&num_followers)

	result = tx.Model(&Community{}).Where("CommunityID = ?", CommunityID).Update("NumFollowers", num_followers-1)
	if result.Error != nil {
		tx.Rollback()
		return "Internal Error", 500
	}
	tx.Commit()
	return "Success", 200
}
