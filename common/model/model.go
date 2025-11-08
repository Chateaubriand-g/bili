package model

import "time"

type User struct {
	ID          uint      `json:"id"`
	UserName    string    `json:"username"`
	PassWord    string    `json:"password"`
	NickName    string    `json:"nickname"`
	Gender      string    `json:"gender"`
	Description string    `json:"description"`
	Email       string    `json:"email"`
	Avatar      string    `json:"avatat"`
	CreateTime  time.Time `json:"create_time"`
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
