package database

import (
	"fmt"

	"github.com/Chateaubriand-g/bili/user_service/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatebase(cfg *config.Config) (*gorm.DB, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/biliuser?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open sql failed: %w", err)
	}

	return db, nil
}
