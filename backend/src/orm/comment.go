package models

import (
	"gorm.io/gorm"
)

func AddComment(
	DB *gorm.DB,
	UserID int32,
	WagerID int32,
	ParentCommentID int32,
	Description string) (string, int32) {

	Comment := Comment{
		UserID:          UserID,
		WagerID:         WagerID,
		ParentCommentID: &ParentCommentID,
		Description:     Description,
	}
	result := DB.Create(&Comment)
	if result.Error != nil {
		return "Internal Error", 500
	}
	return "Success", 200
}
