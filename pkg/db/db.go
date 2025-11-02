package db

import (
	"bili/pkg/config"
)

var DB *gorm.DB

func InitDB(config Config) error {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/steam?charset=utf8mb4&parseTime=True&loc=Local",
						config.DBUser,
						config.DBPassword,
						config.DBHost,
						config.DBPort,
	)

	DB,err = gorm.Open(mysql.Open(dsn),&gorm.Config{})
	if err!=nil {
		return fmt.ErrorF("gorm open failed: %w",err)
	}
	return nil
}