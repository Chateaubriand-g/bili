package dao

import (
	"bili/auth_service/model"

	"gorm.io/gorm"
)

type UserDAO interface {
	FindByUsername(username string) (*model.User, error)
	Create(user *model.User) error
}

type userDAO struct{ DB *gorm.DB }

func NewUserDAO(db *gorm.DB) UserDAO { return &userDAO{DB: db} }

func (d *userDAO) FindByUsername(username string) (*model.User, error) {
	var user model.User

	err := d.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *userDAO) Create(user *model.User) error {
	return d.DB.Create(user).Error
}
