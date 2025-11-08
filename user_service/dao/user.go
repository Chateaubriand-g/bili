package dao

import (
	"fmt"

	"github.com/Chateaubriand-g/bili/common/model"
	"gorm.io/gorm"
)

type UserDAO interface {
	FindByUserID(userID string) (*model.User, error)
	Update(data *model.User) error
}

type userDAO struct {
	DB *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &userDAO{
		DB: db,
	}
}

func (u *userDAO) FindByUserID(userID string) (*model.User, error) {
	var res model.User

	if err := u.DB.Where("id = ?", userID).First(&res).Error; err != nil {
		return nil, fmt.Errorf("sql search failed: %w", err)
	}

	return &res, nil
}

func (u *userDAO) Update(userID string, data *model.User) error {
	if err := u.DB.Model(&model.User{}).Where("id = ?", userID).Updates(data).Error; err != nil {
		return fmt.Errorf("sql updates failed: %w", err)
	}
	return nil
}
