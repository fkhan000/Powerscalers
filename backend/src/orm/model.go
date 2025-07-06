package models

import "time"

type User struct {
	UserID     int32  `gorm:"primaryKey;autoIncrement"`
	UserName   string `gorm:"not null"`
	ProfilePic string
	Cash       float32   `gorm:"default:0"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

type Community struct {
	CommunityID  int32  `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"not null"`
	Description  string
	Picture      string
	NumFollowers int       `gorm:"default:0"`
	RateLimit    int       `gorm:"default:0"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

type Member struct {
	UserID      int32     `gorm:"primaryKey"`
	CommunityID int32     `gorm:"primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Community Community `gorm:"foreignKey:CommunityID;constraint:OnDelete:CASCADE"`
}

type Moderator struct {
	UserID      int32     `gorm:"primaryKey"`
	CommunityID int32     `gorm:"primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Community Community `gorm:"foreignKey:CommunityID;constraint:OnDelete:CASCADE"`
}

type Wager struct {
	WagerID        int32  `gorm:"primaryKey;autoIncrement"`
	CommunityID    int32  // nullable foreign key
	OwnerID        int32  // nullable foreign key
	Title          string `gorm:"not null"`
	Description    string
	Decision       string
	Explanation    string
	NetLikes       int32 `gorm:"default:0"`
	ExpirationDate time.Time
	CreatedAt      time.Time `gorm:"autoCreateTime"`

	Community *Community `gorm:"foreignKey:CommunityID;constraint:OnDelete:SET NULL"`
	Owner     *User      `gorm:"foreignKey:OwnerID;constraint:OnDelete:SET NULL"`
}

type Gamble struct {
	UserID    int32     `gorm:"primaryKey"`
	WagerID   int32     `gorm:"primaryKey"`
	Amount    float32   `gorm:"not null"`
	Position  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User  User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Wager Wager `gorm:"foreignKey:WagerID;constraint:OnDelete:CASCADE"`
}

type Comment struct {
	CommentID   int32     `gorm:"primaryKey;autoIncrement"`
	UserID      int32     `gorm:"not null"`
	WagerID     int32     `gorm:"not null"`
	Description string    `gorm:"not null"`
	NetLikes    int32     `gorm:"default:0"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	User  User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Wager Wager `gorm:"foreignKey:WagerID;constraint:OnDelete:CASCADE"`
}

type WagerLike struct {
	UserID    int32     `gorm:"primaryKey"`
	WagerID   int32     `gorm:"primaryKey"`
	Value     int32     `gorm:"not null;check:value >= -1 AND value <= 1"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User  User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Wager Wager `gorm:"foreignKey:WagerID;constraint:OnDelete:CASCADE"`
}

type CommentLike struct {
	UserID    int32     `gorm:"primaryKey"`
	CommentID int32     `gorm:"primaryKey"`
	Value     int32     `gorm:"not null;check:value >= -1 AND value <= 1"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Comment Comment `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE"`
}
