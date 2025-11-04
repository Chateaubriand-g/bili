package model

import "time"

type User struct {
	ID         uint      `json:"id" gorm:"primarykey,increment"`
	UserName   string    `json:"username" gorm:"size:50,uniqueIndex,not null"`
	PassWord   string    `json:"password" gorm:"size:255,not null"`
	Email      string    `json:"email" gorm:"size:50"`
	NickName   string    `json:"nickname" gorm:"size:50,default:guet"`
	Avatar     string    `json:"avatat" gorm:"size:255,default:''"`
	CreateTime time.Time `json:"createTime" gorm:"autocreatetime"`
}
