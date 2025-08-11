package models

import "time"

type User struct {
	UserID     int    `gorm:"primaryKey;autoIncrement"`
	UserName   string `gorm:"not null"`
	ProfilePic string
	Cash       float32   `gorm:"default:0"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

type Community struct {
	CommunityID        int    `gorm:"primaryKey;autoIncrement"`
	Name               string `gorm:"not null"`
	Description        string
	Picture            string
	NumFollowers       int       `gorm:"default:0"`
	DailyPostsLimit    int       `gorm:"default:20"`
	ModifiedDailyLimit int       `gorm:"default:20"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
}

type Member struct {
	UserID      int       `gorm:"primaryKey"`
	CommunityID int       `gorm:"primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Community Community `gorm:"foreignKey:CommunityID;constraint:OnDelete:CASCADE"`
}

type Moderator struct {
	UserID      int       `gorm:"primaryKey"`
	CommunityID int       `gorm:"primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Community Community `gorm:"foreignKey:CommunityID;constraint:OnDelete:CASCADE"`
}

type Wager struct {
	WagerID        int    `gorm:"primaryKey;autoIncrement"`
	CommunityID    int    // nullable foreign key
	OwnerID        int    // nullable foreign key
	Title          string `gorm:"not null"`
	Description    string `gorm:"not null"`
	Left           string `gorm:"not null"`
	Right          string `gorm:"not null"`
	Decision       string
	Explanation    string
	NumVoR         int       `gorm:"default:5"`
	NetLikes       int       `gorm:"default:0"`
	ExpirationDate time.Time `gorm:"not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`

	Community *Community `gorm:"foreignKey:CommunityID;constraint:OnDelete:SET NULL"`
	Owner     *User      `gorm:"foreignKey:OwnerID;constraint:OnDelete:SET NULL"`
}

type Gamble struct {
	UserID    int       `gorm:"primaryKey"`
	WagerID   int       `gorm:"primaryKey"`
	Amount    float32   `gorm:"not null"`
	Position  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User  User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Wager Wager `gorm:"foreignKey:WagerID;constraint:OnDelete:CASCADE"`
}

type Comment struct {
	CommentID       int `gorm:"primaryKey;autoIncrement"`
	UserID          int `gorm:"not null"`
	WagerID         int
	ParentCommentID *int
	Description     string    `gorm:"not null"`
	NetLikes        int       `gorm:"default:0"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`

	User   User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Wager  Wager    `gorm:"foreignKey:WagerID;constraint:OnDelete:CASCADE"`
	Parent *Comment `gorm:"foreignKey:ParentCommentID;constraint:OnDelete:CASCADE"`
}

type WagerLike struct {
	UserID    int       `gorm:"primaryKey"`
	WagerID   int       `gorm:"primaryKey"`
	Value     int       `gorm:"not null;check:value >= -1 AND value <= 1"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User  User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Wager Wager `gorm:"foreignKey:WagerID;constraint:OnDelete:CASCADE"`
}

type CommentLike struct {
	UserID    int       `gorm:"primaryKey"`
	CommentID int       `gorm:"primaryKey"`
	Value     int       `gorm:"not null;check:value >= -1 AND value <= 1"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Comment Comment `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE"`
}
