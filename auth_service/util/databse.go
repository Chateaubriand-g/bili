package util

import (
	"bili/auth_service/config"
	"bili/auth_service/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CreateDB(config *config.Config) (*gorm.DB, error) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/biliuser?charset=utf8mb4&parseTime=True&loc=Local",
		config.DSN.User,
		config.DSN.Password,
		config.DSN.Host,
		config.DSN.Port,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm open failed: %w", err)
	}

	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate databse: %w", err)
	}
	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
	)
}
