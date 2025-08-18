package models

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type FullComment struct {
	CommentID      int       `gorm:"column:comment_id"`
	UserID         int       `gorm:"column:user_id"`
	UserName       string    `gorm:"column:user_name"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	NetLikes       int       `gorm:"column:net_likes"`
	Description    string    `gorm:"column:description"`
	Position       string    `gorm:"column:position"`
	ProfilePicture string    `gorm:"column:profile_pic"`
}

func AddComment(
	DB *gorm.DB,
	UserID int,
	WagerID int,
	ParentCommentID int,
	Description string) (string, int) {

	Comment := Comment{
		UserID:          UserID,
		WagerID:         WagerID,
		ParentCommentID: &ParentCommentID,
		Description:     Description,
	}

	tx := DB.Begin()

	err := tx.Create(&Comment).Error

	if err != nil {
		tx.Rollback()
		return fmt.Sprintf("Error creating comment: %s", err), 500
	}
	var totalComments int
	tx.Table("wager as w").Select("total_comments").Where("wager_id = ?", WagerID).Scan(totalComments)

	err = tx.Table("wager as w").Where("wager_id", WagerID).Update("total_comments", totalComments).Error
	if err != nil {
		tx.Rollback()
		return fmt.Sprintf("Error updating total comments: %s", err), 500
	}
	tx.Commit()
	return "Success", 200
}

func FetchComments(
	DB *gorm.DB,
	WagerID int,
	SortBy string,
	Order string,
	Offset int,
	Limit int,
	depth int,
) ([]FullComment, error) {

	sortCol := map[string]string{"Likes": "c.net_likes", "Time": "c.created_at"}[SortBy]
	if sortCol == "" {
		return nil, fmt.Errorf("unexpected sort option")
	}
	if strings.ToUpper(Order) != "ASC" && strings.ToUpper(Order) != "DESC" {
		return nil, fmt.Errorf("unexpected order option")
	}
	orderClause := sortCol + " " + strings.ToUpper(Order)

	var allIDs []int
	var parentIDs []int

	for i := 0; i < depth; i++ {
		var nextIDs []int
		q := DB.Table("comments AS c").
			Where("c.wager_id = ?", WagerID)

		if i == 0 {

			q = q.Where("c.parent_comment_id IS NULL").
				Order(orderClause).
				Offset(Offset).
				Limit(Limit)
		} else {
			if len(parentIDs) == 0 {
				break
			}
			q = q.Where("c.parent_comment_id IN ?", parentIDs)
		}

		if err := q.Pluck("c.comment_id", &nextIDs).Error; err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}
		if len(nextIDs) == 0 {
			break
		}
		allIDs = append(allIDs, nextIDs...)
		parentIDs = nextIDs
	}

	if len(allIDs) == 0 {
		return []FullComment{}, nil
	}

	var results []FullComment
	tx := DB.Table("comments AS c").
		Select(`
            c.comment_id,
            c.user_id,
            c.created_at,
            c.net_likes,
            c.description,
            COALESCE(g.position, '') AS position,
            u.profile_pic,
            u.user_name
        `).
		Joins("JOIN users   AS u ON u.user_id = c.user_id").
		Joins("LEFT JOIN gambles AS g ON g.user_id = c.user_id AND g.wager_id = c.wager_id").
		Where("c.comment_id IN ?", allIDs).
		Order(orderClause).
		Offset(Offset).
		Limit(Limit).
		Scan(&results)
	if tx.Error != nil {
		return nil, fmt.Errorf("database error: %w", tx.Error)
	}
	return results, nil
}
