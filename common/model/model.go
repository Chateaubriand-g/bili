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
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time;type:datetime"`
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
