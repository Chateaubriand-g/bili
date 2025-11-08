package model

import "time"

type User struct {
	ID         uint      `json:"id" gorm:"primarykey,increment"`
	UserName   string    `json:"username" gorm:"size:50,uniqueIndex,not null"`
	PassWord   string    `json:"password" gorm:"size:255,not null"`
	Email      string    `json:"email" gorm:"size:50"`
	NickName   string    `json:"nickname" gorm:"size:50,default:guet"`
	Avatar     string    `json:"avatat" gorm:"size:255,default:''"`
	CreateTime time.Time `json:"create_time" gorm:"autocreatetime"`
}

type UserDTO struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	NickName string `json:"nickname"`
	Avatar   string `json:"avatat"`
}

type UserInfoUpdate struct {
	NickName string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}
