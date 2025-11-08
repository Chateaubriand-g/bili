package database

import (
	"fmt"

	"github.com/Chateaubriand-g/bili/common/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CreateDB(dsn string) (*gorm.DB, error) {
	var err error
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
