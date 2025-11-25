package dao

import (
	//"github.com/Chateaubriand-g/bili/auth_service/model"
	"context"
	"strconv"
	"time"

	"github.com/Chateaubriand-g/bili/common/model"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

type UserDAO interface {
	FindByUsername(username string) (*model.User, error)
	Create(user *model.User) error
	SaveRefreshToken(token string, user_id uint64) error
	DeleteRefreshToken(token string, user_id uint64) error
	IsTokenVaild(token string, user_id uint64) (bool, error)
}

type userDAO struct {
	DB  *gorm.DB
	RDS *redis.Client
}

func NewUserDAO(db *gorm.DB, rds *redis.Client) UserDAO { return &userDAO{DB: db, RDS: rds} }

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

func (d *userDAO) SaveRefreshToken(token string, user_id uint64) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	key := "refresh_token:" + token
	value := map[string]interface{}{
		"user_id": user_id,
	}
	if err := d.RDS.HSet(ctx, key, value).Err(); err != nil {
		return err
	}

	return d.RDS.Expire(ctx, key, 24*time.Hour).Err()
}

func (d *userDAO) DeleteRefreshToken(token string, user_id uint64) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	key := "refresh_token" + token
	return d.RDS.Del(ctx, key).Err()
}

func (d *userDAO) IsTokenVaild(token string, user_id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	key := "refresh_token" + token
	data, err := d.RDS.HGet(ctx, key, "user_id").Result()
	if err != nil {
		return false, err
	}
	uidStr := strconv.FormatUint(user_id, 10)
	if data == uidStr {
		return true, nil
	}
	return false, nil
}
