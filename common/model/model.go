package model

import "time"

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserName    string    `json:"username" gorm:"column:username;type:varchar(64)"`
	PassWord    string    `json:"password" gorm:"column:password;type:varchar(64)"`
	NickName    string    `json:"nickname" gorm:"column:nickname;type:varchar(32)"`
	Gender      string    `json:"gender" gorm:"column:gender;type:varchar(32)"`
	Description string    `json:"description" gorm:"column:description;type:text"`
	Email       string    `json:"email" gorm:"column:email;type:varchar(128)"`
	Avatar      string    `json:"avatar" gorm:"column:avatar;type:varchar(255)"`
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time;type:datetime;->"`
}

type Notification struct {
	ID         uint64    `json:"id" gorm:"primarykey"`
	UserID     uint64    `json:"user_id"`
	Type       int8      `json:"type"`
	FromUserID uint64    `json:"from_user_id"`
	BizID      uint64    `json:"biz_id"`
	Payload    string    `json:"payload"`
	IsRead     uint8     `json:"is_read"`
	DeleteddAt time.Time `json:"deleted_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type Comment struct {
	ID        uint64 `gorm:"primarykey"`
	VideoID   uint64 `gorm:"index"`
	UserID    uint64
	Content   string `gorm:"type:text"`
	ParentID  uint64 `gorm:"index;default:0:`
	LikeCount int
}

type CommentsLike struct {
	CommentID uint64 `gorm:"primarykey"`
	UserID    uint64 `gorm:"primarykey"`
}

type Video struct {
	ID       uint64 `gorm:"primarykey"`
	Length   uint64
	FileUrl  string `gorm:"type:varchar(255)"`
	CoverUrl string `gorm:"type:varchar(255)"`
	Name     string `gorm:"varchar(128)"`
	Intro    string `gorm:"varchar(512)"`
	OwnerID  uint64
}
