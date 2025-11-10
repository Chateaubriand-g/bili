package util

import (
	"fmt"

	"github.com/Chateaubriand-g/bili/common/database"
	"github.com/Chateaubriand-g/bili/notification_service/config"

	"gorm.io/gorm"
)

func InitDatabase(config *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/biliuser?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB.User,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
	)
	return database.CreateDB(dsn)
}
