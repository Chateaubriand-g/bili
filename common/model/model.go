package model

import "time"

type User struct {
	ID          uint      `json:"id"`
	UserName    string    `json:"username" gorm:"username;type:varchar(64)"`
	PassWord    string    `json:"password" gorm:"password;type:varchar(64)"`
	NickName    string    `json:"nickname" gorm:"nickname;type:varchar(32)"`
	Gender      string    `json:"gender" gorm:"type:varchar(32)"`
	Description string    `json:"description" gorm:"type:text"`
	Email       string    `json:"email" gorm:"type:varchar(128)"`
	Avatar      string    `json:"avatat" gorm:"type:varchar(255)"`
	CreateTime  time.Time `json:"create_time" gorm:"type:datatime"`
}

type UserDTO struct {
	ID          uint   `json:"id"`
	UserName    string `json:"username"`
	NickName    string `json:"nickname"`
	Gender      string `json:"gender"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
}

type UserInfoUpdate struct {
	NickName    string `json:"nickname"`
	Description string `json:"description"`
	Gender      string `json:"gender"`
}

type UserAvatatUpdate struct {
	Avatar string `json:"avatar"`
}
